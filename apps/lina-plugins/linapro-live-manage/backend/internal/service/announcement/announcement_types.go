// announcement_types.go carries the local type aliases and enumeration
// constants of the announcement service, keeping the impl signatures short
// while staying anchored to the host capability packages.

package announcement

import (
	"lina-core/pkg/plugin/capability/tenantcap"
)

// Announcement visibility values for the enabled column.
const (
	// announcementDisabled marks a hidden announcement row.
	announcementDisabled = false
	// announcementEnabled marks a viewer-visible announcement row.
	announcementEnabled = true
)

// tenantcapFilter aliases the tenant table-filter service type.
type tenantcapFilter = tenantcap.FilterService
