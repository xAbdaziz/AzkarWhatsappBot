package main

import (
	"AzkarWhatsappBot/libs"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
)

var botNum = ""

func registerHandler(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {

		// Send a welcome message when added to a group
		case *events.JoinedGroup:
			if len(v.Participants) > 1 {
				if BotIsAdded(v.Participants, botNum) {
					_, _ = client.SendMessage(context.Background(), v.JID.ToNonAD(), &waProto.Message{Conversation: proto.String(os.Getenv("WELCOME_MSG"))})
				}
			}
			break
		}
	}
}

func main() {

	_, err := os.Stat("config.env")
	if os.IsNotExist(err) {
		fmt.Println("Couldn't find config.env, did you fill sample_config.env and rename it to config.env?")
		os.Exit(1)
		return
	}

	// Load config.env
	_ = godotenv.Load("config.env")

	// Spoof the bot as Windows
	store.DeviceProps.Os = proto.String("Windows")
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()

	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:db.sqlite?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	eventHandler := registerHandler(client)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		botNum = client.Store.ID.ToNonAD().String()
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	libs.Start(client)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func BotIsAdded(participants []types.GroupParticipant, botNum string) bool {

	for _, participant := range participants {
		if participant.JID.ToNonAD().String() == botNum {
			return true
		}
	}
	return false
}
