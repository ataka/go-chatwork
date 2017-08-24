package chatwork

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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

type text string
type roomId int64

type Message struct {
	roomId roomId
	body   text
}

func NewMessage(roomId roomId, body text) *Message {
	m := new(Message)
	m.roomId = roomId
	m.body = body
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
