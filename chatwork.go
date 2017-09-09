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

func (c *Chatwork) post(endpoint endpoint, vs *url.Values) *http.Response {
	reqBody := strings.NewReader(vs.Encode())
	request, requestError := http.NewRequest("POST", string(endpoint), reqBody)
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

func decodeBody(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

type Text string
type RoomId int64

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

func (m *CreateMessageRequest) values() *url.Values {
	vs := url.Values{}
	vs.Add("body", string(m.body))
	return &vs
}

type CreateMessageResponse struct {
	MessageId string `json:"message_id"`
}

func endpointFmt(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a)
}

func (c *Chatwork) CreateMessage(req *CreateMessageRequest) *CreateMessageResponse {
	endpoint := endpoint(baseURL + fmt.Sprintf("rooms/%d/messages", req.roomId))
	httpRes := c.post(endpoint, req.values())

	var res CreateMessageResponse
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
	endpoint := endpoint(baseURL + fmt.Sprintf("rooms/%d/tasks", req.roomId))
	httpRes := c.post(endpoint, req.values())

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
