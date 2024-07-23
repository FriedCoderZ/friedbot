package onebot

import (
	"fmt"
	"time"

	"github.com/FriedCoderZ/friedbot"
)

type API struct{}

func (a API) Install(bot *friedbot.Bot) error {
	bot.API = a
	return nil
}

func (API) SendPrivateMsg(userID int64, message string) (messageID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (API) SendGroupMsg(groupID int64, message string) (messageID int64, err error) {
	fmt.Printf("Send meesage \"%s\" to %d\n", message, groupID)
	return 0, nil
}

func (API) DeleteMsg(messageID int32) (err error) {
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

func (e event) GetMsg() *friedbot.Message {
	if !e.IsMsg() {
		return nil
	}
	sender := e.data["sender"].(map[string]any)
	segments := e.data["message"].([]any)
	msg := &friedbot.Message{
		ID:      int64(e.data["message_id"].(float64)),
		GroupID: int64(e.data["group_id"].(float64)),
		User: &friedbot.User{
			ID:       int64(sender["user_id"].(float64)),
			Nickname: sender["nickname"].(string),
			Sex:      sender["sex"].(string),
			Age:      int(sender["age"].(float64)),
		},
		Time: int64(e.data["time"].(float64)),
	}
	switch e.data["message_type"].(string) {
	case "private":
		msg.Type = friedbot.MsgTypePrivate
	case "group":
		msg.Type = friedbot.MsgTypeGroup
	default:
		msg.Type = friedbot.MsgTypeOther
	}
	for _, segment := range segments {
		data := segment.(map[string]any)["data"].(map[string]any)
		switch segment.(map[string]any)["type"].(string) {
		case "text":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeText,
				Data: friedbot.SegmentText{
					Text: data["text"].(string),
				},
			})
		case "face":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeFace,
				Data: friedbot.SegmentFace{
					ID: data["id"].(string),
				},
			})
		case "image":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeImage,
				Data: friedbot.SegmentImage{
					File: data["file"].(string),
				},
			})
		case "record":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeRecord,
				Data: friedbot.SegmentRecord{
					File: data["file"].(string),
				},
			})
		case "video":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeVideo,
				Data: friedbot.SegmentVideo{
					File: data["file"].(string),
				},
			})
		case "at":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeAt,
				Data: friedbot.SegmentAt{
					QQ: data["qq"].(string),
				},
			})
		case "rps":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeRps,
				Data: friedbot.SegmentRPS{},
			})
		case "dice":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeDice,
				Data: friedbot.SegmentDice{},
			})
		case "shake":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeShake,
				Data: friedbot.SegmentShake{},
			})
		case "poke":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypePoke,
				Data: friedbot.SegmentPoke{
					Type: data["type"].(string),
					ID:   data["id"].(string),
				},
			})
		case "anonymous":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeAnonymous,
				Data: friedbot.SegmentAnonymous{},
			})
		case "share":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeShare,
				Data: friedbot.SegmentShare{
					URL:   data["url"].(string),
					Title: data["title"].(string),
				},
			})
		case "contact":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeContact,
				Data: friedbot.SegmentContact{
					Type: data["type"].(string),
					ID:   data["id"].(string),
				},
			})
		case "location":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeLocation,
				Data: friedbot.SegmentLocation{
					Lat: data["lat"].(string),
					Lon: data["lon"].(string),
				},
			})
		case "music":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeMusic,
				Data: friedbot.SegmentMusic{
					Type:  data["type"].(string),
					ID:    data["id"].(string),
					URL:   data["url"].(string),
					Audio: data["audio"].(string),
					Title: data["title"].(string),
				},
			})
		case "reply":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeReply,
				Data: friedbot.SegmentReply{
					ID: data["id"].(string),
				},
			})
		case "forward":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeForward,
				Data: friedbot.SegmentForward{
					ID: data["id"].(string),
				},
			})
		case "node":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeNode,
				Data: friedbot.SegmentNode{
					ID: data["id"].(string),
				},
			})
		case "xml":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeXML,
				Data: friedbot.SegmentXML{
					Data: data["data"].(string),
				},
			})
		case "json":
			msg.Segments = append(msg.Segments, friedbot.Segment{
				Type: friedbot.SegmentTypeJSON,
				Data: friedbot.SegmentJSON{
					Data: data["data"].(string),
				},
			})
		}
	}
	return msg
}

func (e event) GetContent() string {
	if !e.IsMsg() {
		return ""
	}
	msg := e.data["message"].([]any)[0].(map[string]any)
	if msg["type"].(string) == "text" {
		return msg["data"].(map[string]any)["text"].(string)
	}
	return e.data["raw_message"].(string)
}

func (e event) GetUser() *friedbot.User {
	//TODO implement me
	panic("implement me")
}

func (e event) GetGroup() *friedbot.Group {

	return &friedbot.Group{
		ID: int64(e.data["group_id"].(float64)),
	}
}
