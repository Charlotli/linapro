// This file verifies the frozen approval node snapshot helpers: order
// normalization, parsing, participant boundaries, and attachment validation.

package approval

import (
	"strings"
	"testing"
)

// TestNormalizeSnapshotOrders verifies snapshot orders are rewritten from
// list position so appends and removals keep the chain 1-based consecutive.
func TestNormalizeSnapshotOrders(t *testing.T) {
	nodes := []snapshotNode{
		{Order: 7, ApproverId: 2, ApproverName: "alice"},
		{Order: 3, ApproverId: 3, ApproverName: "bob"},
	}
	normalized := normalizeSnapshotOrders(nodes)
	if normalized[0].Order != 1 || normalized[1].Order != 2 {
		t.Fatalf("expected consecutive orders 1,2 got %d,%d", normalized[0].Order, normalized[1].Order)
	}
	if normalized[0].ApproverId != 2 || normalized[1].ApproverId != 3 {
		t.Fatal("approver identity must survive normalization")
	}
}

// TestParseSnapshotNodes verifies empty snapshots stay empty and malformed
// snapshots are rejected.
func TestParseSnapshotNodes(t *testing.T) {
	nodes, err := parseSnapshotNodes("")
	if err != nil || len(nodes) != 0 {
		t.Fatalf("expected empty snapshot, got %v err=%v", nodes, err)
	}
	if _, err = parseSnapshotNodes("{broken"); err == nil {
		t.Fatal("expected malformed snapshot to fail")
	}
	parsed, err := parseSnapshotNodes(`[{"order":1,"approverId":2,"approverName":"alice"}]`)
	if err != nil || len(parsed) != 1 || parsed[0].ApproverName != "alice" {
		t.Fatalf("expected one parsed node, got %v err=%v", parsed, err)
	}
}

// TestCurrentSnapshotNode verifies the pending node lookup and bounds.
func TestCurrentSnapshotNode(t *testing.T) {
	nodes := []snapshotNode{
		{Order: 1, ApproverId: 2},
		{Order: 2, ApproverId: 3},
	}
	current, ok := currentSnapshotNode(nodes, 1)
	if !ok || current.ApproverId != 2 {
		t.Fatalf("expected first node approver 2, got %v ok=%v", current, ok)
	}
	if _, ok = currentSnapshotNode(nodes, 3); ok {
		t.Fatal("out-of-range node order must not resolve")
	}
	if _, ok = currentSnapshotNode(nodes, 0); ok {
		t.Fatal("zero node order must not resolve")
	}
}

// TestIsRequestParticipant verifies the participant boundary: the applicant
// and any snapshot approver participate; other users do not.
func TestIsRequestParticipant(t *testing.T) {
	request := &RequestEntity{ApplicantId: 1}
	nodes := []snapshotNode{{Order: 1, ApproverId: 2}, {Order: 2, ApproverId: 3}}
	if !isRequestParticipant(request, nodes, 1) {
		t.Fatal("applicant must participate")
	}
	if !isRequestParticipant(request, nodes, 3) {
		t.Fatal("snapshot approver must participate")
	}
	if isRequestParticipant(request, nodes, 9) {
		t.Fatal("non-participant must not participate")
	}
	if isRequestParticipant(nil, nodes, 1) {
		t.Fatal("nil request must not participate")
	}
}

// TestValidateFormValuesAttachments verifies attachment-type dynamic fields
// accept only http(s) URL lists and honor the required flag.
func TestValidateFormValuesAttachments(t *testing.T) {
	fields := []FieldConfig{
		{Key: "invoices", Label: "Invoices", Type: FieldTypeAttachment, Required: true},
	}
	if _, err := validateFormValues(fields, map[string]any{}); err == nil {
		t.Fatal("expected missing required attachment to fail")
	}
	if _, err := validateFormValues(fields, map[string]any{
		"invoices": []any{"ftp://example.com/a.pdf"},
	}); err == nil {
		t.Fatal("expected non-http attachment to fail")
	}
	stored, err := validateFormValues(fields, map[string]any{
		"invoices": []any{"https://example.com/a.pdf"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stored, "https://example.com/a.pdf") {
		t.Fatalf("expected stored JSON to contain the URL, got %q", stored)
	}
}
