package app

import (
	// "fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type Conversation struct {
	*proto.Conversation
}

func (c Conversation) ID() (jsonField string, value interface{}) {
	value = c.GetId()
	jsonField = "id"
	return
}

func eventHandler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		// fmt.Println("Received a message conv!", evt.Message.GetConversation())
	case *events.HistorySync:
		switch evt.Data.GetSyncType() {
		case proto.HistorySync_INITIAL_BOOTSTRAP:
			for _, conversation := range evt.Data.GetConversations() {
				conv := Conversation{conversation}
				Driver.Insert(conv)
			}
		}
	}
}

func Handler(client *whatsmeow.Client) {
	client.AddEventHandler(eventHandler)
}
