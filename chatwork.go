package chatwork

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.1"
	apiVersion     = "v1"
	defaultBaseURL = "https://api.chatwork.com/"
	userAgent      = "go-chatwork/" + libraryVersion
)

type Transport struct {
	Token     string
	Transport http.RoundTripper
}

func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}
func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.Header.Set("X-ChatWorkToken", t.Token)
	return t.transport().RoundTrip(req)
}

type ChatworkResult struct {
	MessageId int      `json:"message_id,int"`
	Errors    []string `json:"errors"`
}
type Chatwork struct {
	client  *http.Client
	url     *url.URL
	version string
	token   string
}

func NewChatwork(httpClient *http.Client) *Chatwork {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseUrl, _ := url.Parse(defaultBaseURL)

	client := &Chatwork{client: httpClient, url: baseUrl, version: apiVersion}
	return client
}

func (chatwork *Chatwork) SendMessage(room_id, body string) (*ChatworkResult, error) {
	data := url.Values{}
	data.Set("body", body)

	u := chatwork.url
	u.Path = fmt.Sprintf("/%s/rooms/%s/messages", chatwork.version, room_id)
	url := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := chatwork.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	result := &ChatworkResult{}
	err = json.Unmarshal(contents, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
