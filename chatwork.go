package chatwork

import (
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

func (c *Chatwork) post(endpoint endpoint, vs url.Values) {
	body := strings.NewReader(vs.Encode())
	request, requestError := http.NewRequest("POST", string(endpoint), body)
	if requestError != nil {
		log.Fatal(requestError)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("X-ChatWorkToken", string(c.apiKey))

	httpClient := new(http.Client)
	res, error := httpClient.Do(request)
	defer res.Body.Close()
	if error != nil {
		log.Fatal(error)
	}
}

type Text string
type RoomId int64

type Message struct {
	roomId RoomId
	body   Text
}

func NewMessage(roomId int64, body string) *Message {
	m := new(Message)
	m.roomId = RoomId(roomId)
	m.body = Text(body)
	return m
}

func endpointFmt(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a)
}

func (c *Chatwork) CreateMessage(message *Message) {
	endpoint := endpoint(baseURL + fmt.Sprintf("rooms/%d/messages", message.roomId))
	vs := url.Values{}
	vs.Add("body", string(message.body))
	c.post(endpoint, vs)
}

type UserId int64
type UserIds []UserId

type Task struct {
	roomId    RoomId
	body      Text
	assignees UserIds
	due       time.Time
}

func NewTask(roomId int64, body string, assignees []int64, due time.Time) *Task {
	t := new(Task)
	t.roomId = RoomId(roomId)
	t.body = Text(body)
	t.assignees = make([]UserId, len(assignees))
	for i, a := range assignees {
		t.assignees[i] = UserId(a)
	}
	t.due = due
	return t
}

func (c *Chatwork) CreateTask(task *Task) {
	endpoint := endpoint(baseURL + fmt.Sprintf("rooms/%d/tasks", task.roomId))
	vs := url.Values{}
	vs.Add("body", string(task.body))
	vs.Add("to_ids", task.assignees.toString(","))
	c.post(endpoint, vs)
}

func (ids UserIds) toString(sep string) string {
	buf := make([]string, len(ids))
	for i, id := range ids {
		buf[i] = strconv.FormatInt(int64(id), 10)
	}
	return strings.Join(buf, sep)
}
