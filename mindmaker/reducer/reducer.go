package reducer

import (
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
		if update.Message == nil {
			return
		}

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
		case "draw":
			r.handleDrawCommand(update)
		case "takeout":
			r.handleTakeoutCommand(update)
		case "buckets":
			r.handleBucketsCommand(update)
		case "start":
			r.handleStartCommand(update)
		default:
			r.handleAllOthersUpdate(update)
		}
	}
}
