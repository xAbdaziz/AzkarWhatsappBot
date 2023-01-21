package libs

import (
	"context"
	"fmt"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"time"
)

type Bot struct {
	Client *whatsmeow.Client
}

func BotClient(client *whatsmeow.Client) *Bot {
	return &Bot{
		Client: client,
	}

}

func (bot *Bot) sendMessage(message string) {
	groups, _ := bot.Client.GetJoinedGroups()
	for _, group := range groups {
		group := group
		go func() {
			_, err := bot.Client.SendMessage(context.Background(), group.JID.ToNonAD(), "", &waProto.Message{
				Conversation: proto.String(message),
			})
			if err != nil {
				println(fmt.Sprintf("Couldn't send a Message to %s Due to %s", group.Name, err.Error()))
			}
		}()
	}
}

func MinuteTicker() *time.Ticker {
	// Return new ticker that triggers on the minute
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
