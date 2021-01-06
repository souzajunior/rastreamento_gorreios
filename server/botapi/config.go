package botapi

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TrackBot is the base struct to storage information about the BOT and connections
type TrackBot struct {
	Bot                 *tgbotapi.BotAPI
	ConnectionsUpdating map[int]map[string]bool
	Updates             tgbotapi.UpdateConfig
	TrackInterval       int64
}

const botKeyName = "BOT_KEY_RASTGORREIOS"
const botKeyTrackInterval = "RASTGORREIOS_TRACK_INTERVAL"

// defaultTrack in minutes
const defaultInterval = 5

// InitBot is the bot initializer
func InitBot() (tBot *TrackBot, err error) {
	// Looking for env key
	var (
		TokenAPI, interval string
		found              bool
		intervalTime       int64
	)

	if TokenAPI, found = os.LookupEnv(botKeyName); !found {
		return nil, errors.New(fmt.Sprintf("failed to find BOT TOKEN {%s} in environment variables", botKeyName))
	}

	if interval, found = os.LookupEnv(botKeyTrackInterval); !found {
		log.Println(fmt.Sprintf("Was not found the environment variable to interval track {%s}. Using default track interval %d minute(s)",
			botKeyTrackInterval, defaultInterval))
	} else {
		intervalTime, _ = strconv.ParseInt(interval, 10, 64)
	}

	var bot *tgbotapi.BotAPI
	if bot, err = tgbotapi.NewBotAPI(TokenAPI); err != nil {
		return
	}

	tBot = &TrackBot{
		Bot:                 bot,
		ConnectionsUpdating: make(map[int]map[string]bool),
		Updates:             tgbotapi.NewUpdate(0),
		TrackInterval:       intervalTime,
	}
	tBot.Updates.Timeout = 60

	log.Printf("Bot %s was initialized\n", bot.Self.UserName)

	return
}

// Send is responsible to send a message in the chat
// params:
// msg - message to send
// rep - reply the message with the message passed in param
func (t *TrackBot) Send(u tgbotapi.Update, msg string, rep bool) {
	newMsg := tgbotapi.NewMessage(u.Message.Chat.ID, msg)

	if rep {
		newMsg.ReplyToMessageID = u.Message.MessageID
	}

	_, _ = t.Bot.Send(newMsg)
}

// SendChat is responsible to send a message in a specific chat
// params:
// chatID - chat ID
// msg - message to send
func (t *TrackBot) SendChat(chatID int64, msg string) {
	newMsg := tgbotapi.NewMessage(chatID, msg)
	_, _ = t.Bot.Send(newMsg)
}
