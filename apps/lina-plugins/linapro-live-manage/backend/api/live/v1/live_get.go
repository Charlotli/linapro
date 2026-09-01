// This file declares the get-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content Get API

// GetReq defines the request for retrieving live-content details.
type GetReq struct {
	g.Meta `path:"/live/{id}" method:"get" tags:"LiveContents" summary:"Get live content details" dc:"Get live content details by ID, including the live room name and creator information. Records outside the current tenant are reported as not found." permission:"live:live:query"`
	Id     int64 `json:"id" v:"required" dc:"Live content ID" eg:"1"`
}

// GetRes Live content detail response
type GetRes struct {
	LiveItem
	CreatedByName string `json:"createdByName" dc:"Creator username" eg:"admin"`
}
