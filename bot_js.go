package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/syumai/tinyutil/httputil"
	"github.com/syumai/workers/cloudflare"
)

var Bot = &tgbotapi.BotAPI{
	Token:  cloudflare.Getenv("telegram_token"),
	Client: httputil.DefaultClient,
	Buffer: 100,
}

func init() {
	Bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
}

func BotHandler() func(w http.ResponseWriter, r *http.Request) {
	idMap := make(map[int64]bool)
	for _, id := range strings.FieldsFunc(cloudflare.Getenv("telegram_ids"), func(r rune) bool { return r == ',' }) {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		idMap[i] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		update, err := Bot.HandleUpdate(r)
		if err != nil {
			json.NewEncoder(w).Encode([]any{
				err.Error(),
			})
			return
		}

		if update.Message == nil || (idMap != nil && !idMap[update.Message.From.ID]) {
			return
		}

		// If we got a message
		fmt.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		argument := update.Message.CommandArguments()

		if argument == "" {
			return
		}

		var msg tgbotapi.Chattable
		switch update.Message.Command() {
		case "image":
			data, err := NewAI().Diffusion(DiffusionOptions{Prompt: argument})
			if err != nil {
				m := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			} else {
				defer data.Close()
				m := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileReader{Name: "image.png", Reader: data})
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			}
		case "user_id":
			m := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprint(update.Message.From.ID))
			m.ReplyToMessageID = update.Message.MessageID
			msg = m
		default:
			return
		}

		Bot.Send(msg)
	}
}
