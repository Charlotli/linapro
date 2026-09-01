// This file declares the equipment options request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment Options API

// OptionsReq defines the request for listing bounded equipment candidates.
type OptionsReq struct {
	g.Meta  `path:"/equipment/options" method:"get" tags:"Equipment" summary:"Get equipment options" dc:"Return the minimal equipment projection (id, code, name, status) for select controls, filtered by keyword across name and code. Scrapped equipment is excluded. At most 200 items are returned; excess results are not paged." permission:"equipment:query"`
	Keyword string `json:"keyword" dc:"Filter by equipment name or code (fuzzy match); empty means all" eg:"Projector"`
}

// OptionsRes Equipment options response
type OptionsRes struct {
	List []*OptionItem `json:"list" dc:"Equipment option list" eg:"[]"`
}

// OptionItem Equipment minimal option projection
type OptionItem struct {
	Id            int64  `json:"id" dc:"Equipment ID" eg:"1"`
	EquipmentCode string `json:"equipmentCode" dc:"Equipment code" eg:"EQ-PROJECTOR-01"`
	EquipmentName string `json:"equipmentName" dc:"Equipment name" eg:"Meeting room projector"`
	Status        int    `json:"status" dc:"Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped" eg:"1"`
}
