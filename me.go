package chatwork

import (
	"log"
)

//
// GetMe
//

type GetMeRequest struct{}

func NewGetMeRequest() *GetMeRequest {
	return &GetMeRequest{}
}

func (m *GetMeRequest) endpoint() endpoint {
	return newEndpoint("me")
}

func (m *GetMeRequest) params() string {
	return ""
}

type GetMeResponse struct {
	UserId                   int64  `json:"account_id"`
	RoomId                   int64  `json:"room_id"`
	Name                     string `json:"name"`
	ChatWorkId               string `json:"chatwork_id"`
	OrganizationId           int64  `json:"organization_id"`
	OrganizationName         string `json:"organization_name"`
	OrganizationDepartment   string `json:"department"`
	OrganizationTitle        string `json:"title"`
	Url                      string `json:"url"`
	Indtroduction            string `json:"introduction"`
	Email                    string `json:"mail"`
	TelOrganization          string `json:"tel_organization"`
	TelOrganizationExtension string `json:"tel_extension"`
	TelMobile                string `json:"tel_mobile"`
	SnsSkype                 string `json:"skype"`
	SnsFacebook              string `json:"facebook"`
	SnsTwitter               string `json:"twitter"`
	AvatarUrl                string `json:"avatar_image_url"`
}

func (c *Chatwork) GetMe(req *GetMeRequest) *GetMeResponse {
	httpRes := c.get(req)

	var res GetMeResponse
	if err := decodeBody(httpRes, &res); err != nil {
		log.Fatal(err)
	}
	return &res
}
