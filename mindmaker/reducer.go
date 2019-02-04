package mindmaker

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/stanleynguyen/mindmaker/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/persistence"
)

// Reducer reduce and handle updates accordingly
type Reducer struct {
	Bot         *tgbotapi.BotAPI
	Persistence persistence.Persistence
}

// NewReducer get a new Reducer
func NewReducer(bot *tgbotapi.BotAPI, db persistence.Persistence) *Reducer {
	return &Reducer{bot, db}
}

// HandleUpdates reduce updates to correct handler function
func (r *Reducer) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch update.Message.Command() {
		case "create":
			r.handleCreateCommand(update)
		case "set":
			r.handleSetCommand(update)
		case "scrap":
			r.handleScrapCommand(update)
		case "list":
			r.handleListCommand(update)
		case "add":
			r.handleAddCommand(update)
		}
	}
}

func (r *Reducer) handleCreateCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I'm confused 😵 My buckets cant have empty names\nPlease tell me in this format: /create <your bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := getBucketNameFromChatID(update.Message.Chat.ID, argStr)
	// TODO: check if bucket already exists
	err := r.Persistence.InsertBucket(bucketName)
	if err != nil {
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	err = r.Persistence.UpdateDefaultBucket(update.Message.Chat.ID, bucketName)
	if err != nil {
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bucket %v created for you, boss!", argStr))
	r.Bot.Send(msg)
}

func (r *Reducer) handleSetCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Boss 😱 Please gimme some name at least by telling me in this format: /set <bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := getBucketNameFromChatID(update.Message.Chat.ID, argStr)
	// TODO: check if bucket exists
	err := r.Persistence.UpdateDefaultBucket(update.Message.Chat.ID, bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You have selected bucket %v to draw decision from 😉", argStr))
	r.Bot.Send(msg)
}

func (r *Reducer) handleScrapCommand(update tgbotapi.Update) {
	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	if argStr == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Boss 😱 Please gimme some name at least by telling me in this format: /scrap <bucket name>")
		r.Bot.Send(msg)
		return
	}

	bucketName := getBucketNameFromChatID(update.Message.Chat.ID, argStr)
	// TODO check if bucket exists
	err := r.Persistence.DeleteBucket(bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bucket %v is out of the window 😶", argStr))
	r.Bot.Send(msg)
}

func (r *Reducer) handleListCommand(update tgbotapi.Update) {
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to read from 😕")
		r.Bot.Send(msg)
		return
	}

	options, err := r.Persistence.ReadAllOptions(bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFormattedListOfOptions(options))
	r.Bot.Send(msg)
}

func (r *Reducer) handleAddCommand(update tgbotapi.Update) {
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to add this decision to 😕")
		r.Bot.Send(msg)
		return
	}

	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	err = r.Persistence.InsertOption(bucketName, domain.Option(argStr))
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Decision %v addede to bucket %v", argStr, bucketName))
	r.Bot.Send(msg)
}

func (r *Reducer) sendErrMessage(chatID int64) {
	errMsg := tgbotapi.NewMessage(chatID, "Sorry I'm not feeeling very well :( Please try again later")
	r.Bot.Send(errMsg)
}

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func getBucketNameFromChatID(chatID int64, userGivenName string) string {
	return strconv.Itoa(int(chatID)) + " - " + userGivenName
}

func getFormattedListOfOptions(options []domain.Option) string {
	if len(options) == 0 {
		return "Opps! Your current bucket doesnt contain any options"
	}

	rtv := "Possible decisions in your current default bucket:\n"
	for i := 0; i < len(options); i++ {
		rtv += fmt.Sprintf("- %v", options[i])
		if i < len(options)-1 {
			rtv += "\n"
		}
	}

	return rtv
}
