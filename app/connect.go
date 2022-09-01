package app

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	// waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/mdp/qrterminal/v3"
)

var Client *whatsmeow.Client
var Store *sqlstore.SQLStore

func Connect() {
	// dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:store.db?_foreign_keys=on", nil)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	Store = sqlstore.NewSQLStore(container, *deviceStore.ID)
	// clientLog := waLog.Stdout("Client", "INFO", true)
	Client = whatsmeow.NewClient(deviceStore, nil)

	if Client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := Client.GetQRChannel(context.Background())
		err = Client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = Client.Connect()
		if err != nil {
			panic(err)
		}
	}
}
