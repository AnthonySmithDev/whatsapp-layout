package app

import (
	db "github.com/sonyarouje/simdb"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var Driver *db.Driver

func Database() {
	var err error
	Driver, err = db.New("data")
	if err != nil {
		panic(err)
	}
}

type Conversation struct {
	*proto.Conversation
}

type Message struct {
	*proto.HistorySyncMsg
}

func (c Conversation) ID() (jsonField string, value interface{}) {
	value = c.GetId()
	jsonField = "id"
	return
}

func (c *Conversation) IsGroup() bool {
	return c.Name != nil
}

func (c *Conversation) GetJID() types.JID {
	jid, err := types.ParseJID(c.GetId())
	if err != nil {
		panic(err)
	}
	return jid
}

func (c *Conversation) GetMessagesNew() []*proto.HistorySyncMsg {
	return reverseMessage(c.GetMessages())
}

func reverseMessage(msg []*proto.HistorySyncMsg) []*proto.HistorySyncMsg {
	for i, j := 0, len(msg)-1; i < j; i, j = i+1, j-1 {
		msg[i], msg[j] = msg[j], msg[i]
	}
	return msg
}

func findById(id string) *Conversation {
	var conversation *Conversation
	err := Driver.Open(Conversation{}).Where("id", "=", id).First().AsEntity(&conversation)
	if err != nil {
		panic(err)
	}
	return conversation
}

func findAll() []*Conversation {
	var conversations []*Conversation
	err := Driver.Open(Conversation{}).Get().AsEntity(&conversations)
	if err != nil {
		panic(err)
	}
	return conversations
}

func addMessage(conv *Conversation, messa proto.Message) {
	// conv.Messages = append(conv.Messages, messa)
	err := Driver.Update(conv)
	if err != nil {
		panic(err)
	}

}
