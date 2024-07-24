package onebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FriedCoderZ/friedbot"
)

type API struct {
	port int
}

func NewAPI(port int) *API {
	return &API{
		port: port,
	}
}

func (a API) request(api string, data map[string]any) (map[string]any, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://localhost:%d%s", a.port, api)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a API) Install(bot *friedbot.Bot) error {
	bot.API = a
	return nil
}

func (a API) SendPrivateMsg(userID int64, message string) (int64, error) {
	data := map[string]any{
		"user_id": userID,
		"message": message,
	}
	res, err := a.request("/send_private_msg", data)
	if err != nil {
		return 0, err
	}
	messageID, ok := res["message_id"].(float64)
	if !ok {
		return 0, nil
	}
	return int64(messageID), nil
}

func (a API) SendGroupMsg(groupID int64, message string) (int64, error) {
	data := map[string]any{
		"group_id": groupID,
		"message":  message,
	}
	res, err := a.request("/send_group_msg", data)
	if err != nil {
		return 0, err
	}
	messageID, ok := res["message_id"].(float64)
	if !ok {
		return 0, nil
	}
	return int64(messageID), nil
}

func (a API) DeleteMsg(messageID int32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a API) GetMsg(messageID int32) (msg *friedbot.Message, err error) {
	//TODO implement me
	panic("implement me")
}

func (a API) ParseEvent(requestBody map[string]any) (friedbot.Event, error) {
	return newEvent(requestBody), nil
}

func (API) GetCache(event friedbot.Event) (events *friedbot.Ring, err error) {
	fmt.Println(event)
	cacheService := friedbot.GetService()
	data := event.GetData()
	switch data["post_type"] {
	case "message":
		switch data["message_type"] {
		case "group":
			return cacheService.GroupEventsPool.Key(fmt.Sprintf("%d", data["group_id"])), nil
		case "private":
			return cacheService.PrivateEventsPool.Key(fmt.Sprintf("%d", data["user_id"])), nil
		}
	case "notice":
		return cacheService.NoticeEventsPool.Key(data["notice_type"].(string)), nil
	case "request":
		return cacheService.RequestEventsPool.Key(data["request_type"].(string)), nil
	case "meta_event":
		return cacheService.MetaEventsPool.Key(data["meta_event_type"].(string)), nil
	}
	return nil, nil
}

type event struct {
	data map[string]any
	time time.Time
}

func newEvent(data map[string]any) *event {
	return &event{
		data: data,
		time: time.Now(),
	}
}

func (e event) IsMsg() bool {
	return e.data["post_type"].(string) == "message"
}

func (e event) GetData() map[string]any {
	return e.data
}

func (e event) GetTime() time.Time {
	return e.time
}

func parserUser(user map[string]any) *friedbot.User {
	id, _ := user["user_id"].(float64)
	nickname, _ := user["nickname"].(string)
	card, _ := user["card"].(string)
	sex, _ := user["sex"].(string)
	age, _ := user["age"].(float64)
	joinTime, _ := user["join_time"].(float64)
	lastSentTime, _ := user["last_sent_time"].(float64)
	level, _ := user["level"].(float64)
	role, _ := user["role"].(string)
	unFriendly, _ := user["unfriendly"].(bool)
	title, _ := user["title"].(string)
	titleExpireTime, _ := user["title_expire_time"].(float64)
	return &friedbot.User{
		ID:              int64(id),
		Nickname:        nickname,
		Card:            card,
		Sex:             sex,
		Age:             int(age),
		JoinTime:        int32(joinTime),
		LastSentTime:    int32(lastSentTime),
		Level:           int(level),
		Role:            role,
		Unfriendly:      unFriendly,
		Title:           title,
		TitleExpireTime: int32(titleExpireTime),
	}
}

func parserSegments(segment map[string]any) *friedbot.Segment {
	data, ok := segment["data"].(map[string]any)
	if !ok {
		return nil
	}
	typ, ok := data["type"].(string)
	if !ok {
		return nil
	}
	var s *friedbot.Segment
	switch typ {
	default:
		return nil
	case "text":
		text, _ := data["text"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeText,
			Data: friedbot.SegmentText{
				Text: text,
			},
		}
	case "face":
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeFace,
			Data: friedbot.SegmentFace{
				ID: id,
			},
		}
	case "image":
		file, _ := data["file"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeImage,
			Data: friedbot.SegmentImage{
				File: file,
			},
		}
	case "record":
		file, _ := data["file"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeRecord,
			Data: friedbot.SegmentRecord{
				File: file,
			},
		}
	case "video":
		file, _ := data["file"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeVideo,
			Data: friedbot.SegmentVideo{
				File: file,
			},
		}
	case "at":
		qq, _ := data["qq"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeAt,
			Data: friedbot.SegmentAt{
				QQ: qq,
			},
		}
	case "rps":
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeRps,
			Data: friedbot.SegmentRPS{},
		}
	case "dice":
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeDice,
			Data: friedbot.SegmentDice{},
		}
	case "shake":
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeShake,
			Data: friedbot.SegmentShake{},
		}
	case "poke":
		t, _ := data["type"].(string)
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypePoke,
			Data: friedbot.SegmentPoke{
				Type: t,
				ID:   id,
			},
		}
	case "anonymous":
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeAnonymous,
			Data: friedbot.SegmentAnonymous{},
		}
	case "share":
		url, _ := data["url"].(string)
		title, _ := data["title"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeShare,
			Data: friedbot.SegmentShare{
				URL:   url,
				Title: title,
			},
		}
	case "contact":
		t, _ := data["type"].(string)
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeContact,
			Data: friedbot.SegmentContact{
				Type: t,
				ID:   id,
			},
		}
	case "location":
		lat, _ := data["lat"].(string)
		lon, _ := data["lon"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeLocation,
			Data: friedbot.SegmentLocation{
				Lat: lat,
				Lon: lon,
			},
		}
	case "music":
		t, _ := data["type"].(string)
		id, _ := data["id"].(string)
		url, _ := data["url"].(string)
		audio, _ := data["audio"].(string)
		title, _ := data["title"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeMusic,
			Data: friedbot.SegmentMusic{
				Type:  t,
				ID:    id,
				URL:   url,
				Audio: audio,
				Title: title,
			},
		}
	case "reply":
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeReply,
			Data: friedbot.SegmentReply{
				ID: id,
			},
		}
	case "forward":
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeForward,
			Data: friedbot.SegmentForward{
				ID: id,
			},
		}
	case "node":
		id, _ := data["id"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeNode,
			Data: friedbot.SegmentNode{
				ID: id,
			},
		}
	case "xml":
		d, _ := data["data"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeXML,
			Data: friedbot.SegmentXML{
				Data: d,
			},
		}
	case "json":
		d, _ := data["data"].(string)
		s = &friedbot.Segment{
			Type: friedbot.SegmentTypeJSON,
			Data: friedbot.SegmentJSON{
				Data: d,
			},
		}
	}
	return s
}

func (e event) GetMsg() *friedbot.Message {
	if !e.IsMsg() {
		return nil
	}
	var sender map[string]any
	sender, _ = e.data["sender"].(map[string]any)
	messageID, _ := e.data["message_id"].(float64)
	messageType, _ := e.data["message_type"].(string)
	groupID, _ := e.data["group_id"].(float64)
	segments, _ := e.data["message"].([]any)
	content, _ := e.data["raw_message"].(string)
	timestamp, _ := e.data["time"].(float64)
	msg := &friedbot.Message{
		ID:      int64(messageID),
		GroupID: int64(groupID),
		User:    parserUser(sender),
		Content: content,
		Time:    int64(timestamp),
	}
	switch messageType {
	case "private":
		msg.Type = friedbot.MsgTypePrivate
	case "group":
		msg.Type = friedbot.MsgTypeGroup
	default:
		msg.Type = friedbot.MsgTypeOther
	}
	for _, segment := range segments {
		if s, ok := segment.(map[string]any); ok {
			res := parserSegments(s)
			if res != nil {
				msg.Segments = append(msg.Segments, *res)
			}
		}
	}
	return msg
}
