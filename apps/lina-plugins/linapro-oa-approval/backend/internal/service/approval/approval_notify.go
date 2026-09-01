// This file implements the plugin-owned approval notification fan-out through
// the host notify capability. Notification copy is localized through the host
// i18n capability with English fallback text, and send failures degrade to
// log entries so approval actions never roll back because of notifications.

package approval

import (
	"context"
	"strconv"

	usermsgv1 "lina-core/api/usermsg/v1"
	"lina-core/pkg/logger"
	"lina-core/pkg/plugin/capability/i18ncap"
	"lina-core/pkg/plugin/capability/notifycap"
)

// notifyTemplateKeys compose the localized notification copy. Each entry pairs
// one plugin i18n key with its English fallback text.
const (
	notifyKeySubmit  = "plugin.linapro-oa-approval.notify.submit"
	notifyKeyNext    = "plugin.linapro-oa-approval.notify.next"
	notifyKeyFinal   = "plugin.linapro-oa-approval.notify.final"
	notifyKeyReject  = "plugin.linapro-oa-approval.notify.reject"
	notifyKeyAppends = "plugin.linapro-oa-approval.notify.append"
	notifyKeyWit     = "plugin.linapro-oa-approval.notify.withdraw"
	notifyKeyAdvance = "plugin.linapro-oa-approval.notify.advance"
)

// notifyFallbacks stores the English fallback copy for every template key.
var notifyFallbacks = map[string]string{
	notifyKeySubmit:  "A new approval request is waiting for your approval",
	notifyKeyNext:    "An approval request advanced to your approval node",
	notifyKeyFinal:   "Your approval request has been approved",
	notifyKeyReject:  "Your approval request has been rejected",
	notifyKeyAppends: "You have been added as an approver to an approval request",
	notifyKeyWit:     "An approval request waiting for your approval was withdrawn",
	notifyKeyAdvance: "Your approval request advanced to the next approval node",
}

// approvalNotifier sends one localized approval notification through the host
// notify capability.
type approvalNotifier struct {
	notifySvc notifycap.Service
	i18nSvc   i18ncap.Service
}

// send delivers one notification to the target user with a localized title and
// the request summary as content. Failures degrade to warning logs.
func (n *approvalNotifier) send(
	ctx context.Context,
	templateKey string,
	recipientUserID int64,
	requestID int64,
	requestTitle string,
) {
	if n == nil || n.notifySvc == nil || recipientUserID <= 0 {
		return
	}
	fallback, ok := notifyFallbacks[templateKey]
	if !ok {
		return
	}
	title := fallback
	if n.i18nSvc != nil {
		title = n.i18nSvc.Translate(ctx, templateKey, fallback)
	}
	_, err := n.notifySvc.Send(ctx, notifycap.SendInput{
		SourceType: usermsgv1.SourceTypePlugin,
		SourceID:   strconv.FormatInt(requestID, 10),
		Recipients: []string{strconv.FormatInt(recipientUserID, 10)},
		Title:      title,
		Content:    requestTitle,
		Category:   notifycap.CategoryCodeOther,
	})
	if err != nil {
		logger.Warningf(
			ctx,
			"linapro-oa-approval notification send failed request=%d recipient=%d err=%v",
			requestID,
			recipientUserID,
			err,
		)
	}
}
