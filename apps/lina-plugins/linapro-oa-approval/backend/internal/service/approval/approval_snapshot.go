// This file implements the frozen approval node snapshot: JSON serialization,
// participant checks, and order normalization for append actions. Runtime
// actions always work on the snapshot so later flow configuration changes
// never affect already submitted requests.

package approval

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
)

// snapshotNode mirrors one frozen approval node inside node_snapshot.
type snapshotNode struct {
	Order        int64  `json:"order"`
	ApproverId   int64  `json:"approverId"`
	ApproverName string `json:"approverName"`
}

// marshalSnapshotNodes serializes the ordered node snapshot into the stored
// JSON text.
func marshalSnapshotNodes(nodes []snapshotNode) (string, error) {
	normalized := normalizeSnapshotOrders(nodes)
	content, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// parseSnapshotNodes deserializes the stored node snapshot text and returns an
// empty slice for empty input.
func parseSnapshotNodes(snapshot string) ([]snapshotNode, error) {
	trimmed := strings.TrimSpace(snapshot)
	if trimmed == "" {
		return []snapshotNode{}, nil
	}
	nodes := make([]snapshotNode, 0)
	if err := json.Unmarshal([]byte(trimmed), &nodes); err != nil {
		// A malformed stored snapshot is a data-integrity problem; wrap it so
		// callers never surface a raw JSON parser error to HTTP clients.
		return nil, gerror.Wrap(err, "parse approval node snapshot failed")
	}
	return normalizeSnapshotOrders(nodes), nil
}

// normalizeSnapshotOrders rewrites node orders from list position so the
// snapshot always stays 1-based and consecutive after appends.
func normalizeSnapshotOrders(nodes []snapshotNode) []snapshotNode {
	normalized := make([]snapshotNode, len(nodes))
	copy(normalized, nodes)
	for index := range normalized {
		normalized[index].Order = int64(index + 1)
	}
	return normalized
}

// currentSnapshotNode returns the pending node of one pending request.
func currentSnapshotNode(nodes []snapshotNode, currentNodeOrder int) (snapshotNode, bool) {
	if currentNodeOrder < 1 || currentNodeOrder > len(nodes) {
		return snapshotNode{}, false
	}
	return nodes[currentNodeOrder-1], true
}

// isRequestParticipant reports whether the user is the applicant or one of the
// snapshot approvers. Non-participants are treated as not-found by callers so
// request existence never leaks.
func isRequestParticipant(request *RequestEntity, nodes []snapshotNode, userID int64) bool {
	if request == nil {
		return false
	}
	if request.ApplicantId == userID {
		return true
	}
	for _, node := range nodes {
		if node.ApproverId == userID {
			return true
		}
	}
	return false
}
