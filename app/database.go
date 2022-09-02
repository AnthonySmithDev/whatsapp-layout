package app

import (
	// "fmt"
	"strings"

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

func (c *Conversation) Title() string {
	if c.IsGroup() {
		return *c.Name
	} else {
		contact, err := Store.GetContact(c.GetJID())
		if err != nil {
			panic(err)
		}
		return contact.FullName
	}
}

func (c *Conversation) Desc() string {
	message := c.GetMessages()[0].GetMessage()
	return getTypeString(message)
}

func getTypeString(message *proto.WebMessageInfo) string {
	if message.GetMessage().GetStickerMessage() != nil {
		return stickerStyle.Render("Sticker")
	}
	if message.GetMessage().GetImageMessage() != nil {
		return imageStyle.Render("Image")
	}
	if message.GetMessage().GetAudioMessage() != nil {
		return audioStyle.Render("Audio")
	}
	if message.GetMessage().GetVideoMessage() != nil {
		return videoStyle.Render("Video")
	}
	if message.GetMessage().GetExtendedTextMessage() != nil {
		return linkStyle.Render("Link")
	}
	if message.GetMessageStubType() == proto.WebMessageInfo_GROUP_CREATE {
		return groupStyle.Render("Group created")
	}
	if message.GetMessageStubType() == proto.WebMessageInfo_GROUP_PARTICIPANT_ADD {
		return groupStyle.Render("You were added to the group")
	}
	if message.GetMessageStubType() == proto.WebMessageInfo_GROUP_PARTICIPANT_REMOVE {
		numbers := message.GetMessageStubParameters()
		return groupStyle.Render("removed from group: ") + strings.Join(numbers, ", ")
	}
	return message.GetMessage().GetConversation()
}

func (c *Conversation) StylesMessages() []string {
	var messages []string
	for _, message := range c.GetReverseMessages() {
		text := getTypeString(message.GetMessage())

		if message.GetMessage().GetKey().GetFromMe() {
			messages = append(messages, youStyle.Render("You: ")+text)
		} else {
			var name string
			if c.IsGroup() {
				jid := message.GetMessage().GetParticipant()
				if jid, err := types.ParseJID(jid); err != nil {
					contact, _ := Store.GetContact(jid)
					name = contact.PushName
				}
			} else {
				jid := message.GetMessage().GetKey().GetRemoteJid()
				if jid, err := types.ParseJID(jid); err != nil {
					contact, _ := Store.GetContact(jid)
					name = contact.FullName
				}
			}
			messages = append(messages, defaulStyle.Render(name+": ")+text)
		}
	}
	return messages
}

func createConversation(conversations []*proto.Conversation) {
	for _, conversation := range conversations {
		Driver.Insert(Conversation{conversation})
	}
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

func (c *Conversation) GetReverseMessages() []*proto.HistorySyncMsg {
	return reverseMessage(c.GetMessages())
}

func reverseMessage(msg []*proto.HistorySyncMsg) []*proto.HistorySyncMsg {
	for i, j := 0, len(msg)-1; i < j; i, j = i+1, j-1 {
		msg[i], msg[j] = msg[j], msg[i]
	}
	return msg
}

type PushName struct {
	*proto.Pushname
}

func (c PushName) ID() (jsonField string, value interface{}) {
	value = c.GetId()
	jsonField = "id"
	return
}

func createGroupPushName(names []*proto.Pushname) {
	for _, name := range names {
		Driver.Insert(PushName{name})
	}
}

func groupFindById(id string) PushName {
	var pushName PushName
	err := Driver.Open(PushName{}).Where("id", "=", id).First().AsEntity(&pushName)
	if err != nil {
		return PushName{}
	}
	return pushName
}
