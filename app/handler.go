package app

import (
	"fmt"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func eventHandler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppState:
		fmt.Println("AppState event")
	case *events.AppStateSyncComplete:
		fmt.Println("AppStateSyncComplete event")
	case *events.Message:
		// conv := findById(evt.Info.Chat.String())
		// addMessage(conv, evt.Message)
		fmt.Println("Received a message conv!", evt.Message.GetConversation())
	case *events.Receipt:
		if evt.Type == events.ReceiptTypeRead || evt.Type == events.ReceiptTypeReadSelf {
			fmt.Println(fmt.Sprintf("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp))
		} else if evt.Type == events.ReceiptTypeDelivered {
			fmt.Println(fmt.Sprintf("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp))
		}
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

func Handler() {
	Client.AddEventHandler(eventHandler)
}
