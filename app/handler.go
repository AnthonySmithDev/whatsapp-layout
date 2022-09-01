package app

import (
	"encoding/json"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func eventHandler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		fmt.Println("Received a message conv!", evt.Message.GetConversation())
	case *events.HistorySync:
		switch evt.Data.GetSyncType() {
		case proto.HistorySync_INITIAL_BOOTSTRAP:
			fileName := fmt.Sprintf("history.json")
			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic("Failed to open file to write history sync: " + err.Error())
			}
			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			err = enc.Encode(evt.Data.GetConversations())
			if err != nil {
				panic("Failed to write history sync: " + err.Error())
			}
			_ = file.Close()
		}
	}
}

func Handler(client *whatsmeow.Client) {
	client.AddEventHandler(eventHandler)
}
