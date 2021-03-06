package chatwork

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.chatwork.com/v2/"

type ApiKey string

type Chatwork struct {
	apiKey ApiKey
}

func NewChatwork(apiKey string) *Chatwork {
	c := new(Chatwork)
	c.apiKey = ApiKey(apiKey)
	return c
}

type endpoint string

func newEndpoint(format string, a ...interface{}) endpoint {
	return endpoint(baseURL + fmt.Sprintf(format, a...))
}

type chatworkRequest interface {
	endpoint() endpoint
}

type chatworkPostRequest interface {
	chatworkRequest
	values() *url.Values
}

type chatworkGetRequest interface {
	chatworkRequest
	params() string
}

func (c *Chatwork) post(req chatworkPostRequest) *http.Response {
	reqBody := strings.NewReader(req.values().Encode())
	request, requestError := http.NewRequest("POST", string(req.endpoint()), reqBody)
	if requestError != nil {
		log.Fatal(requestError)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("X-ChatWorkToken", string(c.apiKey))

	httpClient := new(http.Client)
	res, error := httpClient.Do(request)
	if error != nil {
		log.Fatal(error)
	}
	return res
}

func (c *Chatwork) get(req chatworkGetRequest) *http.Response {
	url, params := string(req.endpoint()), req.params()
	if len(params) > 0 {
		url = url + "?" + params
	}
	request, requestError := http.NewRequest("", url, nil)
	if requestError != nil {
		log.Fatal(requestError)
	}

	request.Header.Add("X-ChatWorkToken", string(c.apiKey))

	httpClient := new(http.Client)
	res, error := httpClient.Do(request)
	if error != nil {
		log.Fatal(error)
	}
	return res
}

func decodeBody(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

type Text string
type RoomId int64
type MessageId string

type CreateMessageRequest struct {
	roomId RoomId
	body   Text
}

func NewCreateMessageRequest(roomId int64, body string) *CreateMessageRequest {
	m := new(CreateMessageRequest)
	m.roomId = RoomId(roomId)
	m.body = Text(body)
	return m
}

func (m *CreateMessageRequest) endpoint() endpoint {
	return newEndpoint("rooms/%d/messages", m.roomId)
}

func (m *CreateMessageRequest) values() *url.Values {
	vs := url.Values{}
	vs.Add("body", string(m.body))
	return &vs
}

type CreateMessageResponse struct {
	MessageId string `json:"message_id"`
}

func (c *Chatwork) CreateMessage(req *CreateMessageRequest) *CreateMessageResponse {
	httpRes := c.post(req)

	var res CreateMessageResponse
	if err := decodeBody(httpRes, &res); err != nil {
		log.Fatal(err)
	}
	return &res
}

type GetMessageRequest struct {
	roomId    RoomId
	messageId MessageId
}

func NewGetMessageRequest(roomId int64, messageId string) *GetMessageRequest {
	m := new(GetMessageRequest)
	m.roomId = RoomId(roomId)
	m.messageId = MessageId(messageId)
	return m
}

func (m *GetMessageRequest) endpoint() endpoint {
	return newEndpoint("rooms/%d/messages/%s", m.roomId, m.messageId)
}

func (m *GetMessageRequest) params() string {
	return ""
}

type GetMessageResponse struct {
	MessageId string `json:"message_id"`
	User      struct {
		UserId    int64  `json:"account_id"`
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_image_url"`
	} `json:"account"`
	Body     string `json:"body"`
	sendAt   int64  `json:"send_time"`
	updateAt int64  `json:"update_time"`
}

func (c *Chatwork) GetMessage(req *GetMessageRequest) *GetMessageResponse {
	httpRes := c.get(req)

	var res GetMessageResponse
	if err := decodeBody(httpRes, &res); err != nil {
		log.Fatal(err)
	}
	return &res
}

type UserId int64
type UserIds []UserId

type CreateTaskRequest struct {
	roomId    RoomId
	body      Text
	assignees UserIds
	due       *time.Time
}

func (t *CreateTaskRequest) endpoint() endpoint {
	return newEndpoint("rooms/%d/tasks", t.roomId)
}

func (t *CreateTaskRequest) values() *url.Values {
	vs := url.Values{}
	vs.Add("body", string(t.body))
	vs.Add("to_ids", t.assignees.toString(","))
	if t.due != nil {
		vs.Add("limit", strconv.FormatInt(t.due.Unix(), 10))
	}
	return &vs
}

type CreateTaskResponse struct {
	TaskIds []int64 `json:"task_ids"`
}

func NewCreateTaskRequest(roomId int64, body string, assignees []int64, due *time.Time) *CreateTaskRequest {
	t := new(CreateTaskRequest)
	t.roomId = RoomId(roomId)
	t.body = Text(body)
	t.assignees = make([]UserId, len(assignees))
	for i, a := range assignees {
		t.assignees[i] = UserId(a)
	}
	t.due = due
	return t
}

func (c *Chatwork) CreateTask(req *CreateTaskRequest) *CreateTaskResponse {
	httpRes := c.post(req)

	var res CreateTaskResponse
	if err := decodeBody(httpRes, &res); err != nil {
		log.Fatal(err)
	}
	return &res
}

func (ids UserIds) toString(sep string) string {
	buf := make([]string, len(ids))
	for i, id := range ids {
		buf[i] = strconv.FormatInt(int64(id), 10)
	}
	return strings.Join(buf, sep)
}
