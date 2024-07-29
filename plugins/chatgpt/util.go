package chatgpt

import (
	"github.com/FriedCoderZ/friedbot"
)

func getLinkByMsg(msg *friedbot.Message) *link {
	id := msg.User.ID
	if msg.Type == friedbot.MsgTypeGroup {
		id = -msg.GroupID
	}
	l, ok := links[id]
	if !ok {
		l = newLink()
		links[id] = l
	}
	return l
}
