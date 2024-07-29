package chatgpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/FriedCoderZ/friedbot/internal/util"
)

const (
	roleTypeUser      = "user"
	roleTypeAssistant = "assistant"
	roleTypeSystem    = "system"
	roleTypeError     = "error"
)

var (
	links map[int64]*link
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type chatError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type chatResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []choice  `json:"choices"`
	Usage   usage     `json:"usage"`
	Message string    `json:"message"`
	Error   chatError `json:"error"`
}

type link struct {
	messages []message
	expTime  time.Time
	*sync.RWMutex
}

func newLink() *link {
	m := make([]message, len(basicMessages))
	copy(m, basicMessages)
	return &link{
		messages: m,
		expTime:  time.Now().Add(expDuration),
		RWMutex:  &sync.RWMutex{},
	}
}

func (l *link) refreshExp() {
	currentTime := time.Now()
	l.expTime = currentTime.Add(expDuration)
}

func (l *link) isValid() bool {
	if len(l.messages) == 0 {
		return false
	}
	currentTime := time.Now()
	if currentTime.After(l.expTime) {
		return false
	}
	return true
}

func (l *link) addMsg(msg *message) {
	l.Lock()
	defer l.Unlock()
	l.messages = append(l.messages, *msg)
	l.refreshExp()
}

func (l *link) withdrawMsg() {
	l.Lock()
	defer l.Unlock()
	if len(l.messages) < 2 {
		return
	}
	l.messages = l.messages[:len(l.messages)-2]
	l.refreshExp()
}

func (l *link) generateMsg() (*message, error) {
	l.RLock()
	defer l.RUnlock()
	request := util.Request{
		URL: api,
		Data: map[string]any{
			"model":    model,
			"messages": l.messages,
		},
		Headers: map[string]string{
			"Authorization": "Bearer " + key,
		},
	}
	r, err := request.POST()
	if err != nil {
		return nil, err
	}
	resp := &chatResponse{}
	err = json.Unmarshal(r, resp)
	if err != nil {
		return nil, err
	}
	if resp.Message != "" {
		return nil, errors.New(resp.Message)
	}
	if resp.Error.Code != "" {
		err = fmt.Errorf("[%s]\n%s", resp.Error.Code, resp.Error.Message)
		switch resp.Error.Code {
		case "context_length_exceeded":
			c := fmt.Sprintf("消息长度达到了上限，请使用 %s 来创建一条新的会话", triggerWordCreate[0])
			m := &message{Role: roleTypeError, Content: c}
			return m, err
		default:
			return nil, err
		}
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("连接ChatGPT失败")
	}
	return &resp.Choices[0].Message, nil
}
