package main

import (
	"encoding/json"
	"net/http"

	ai "github.com/Asutorufa/telegram-workers-ai"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func main() {
	http.HandleFunc("/tgbot", ai.BotHandler())

	http.HandleFunc("/tgbot/register", func(w http.ResponseWriter, r *http.Request) {
		wh, err := tgbotapi.NewWebhook(cloudflare.Getenv("worker_url") + "/tgbot")
		if err != nil {
			json.NewEncoder(w).Encode([]any{
				err.Error(),
			})
			return
		}

		_, err = ai.Bot.Request(wh)
		if err != nil {
			json.NewEncoder(w).Encode([]any{
				err.Error(),
			})
			return
		}

		resp, err := ai.Bot.Request(tgbotapi.NewSetMyCommands(
			tgbotapi.BotCommand{Command: "image", Description: "generate image by prompt"},
			tgbotapi.BotCommand{Command: "dreamshaper", Description: "generate image by @cf/lykon/dreamshaper-8-lcm"},
			tgbotapi.BotCommand{Command: "llama38binstruct", Description: "generate text by @cf/meta/llama-3-8b-instruct"},
			tgbotapi.BotCommand{Command: "mistral7binstruct", Description: "generate text by mistral-7b-instruct-v0.2-lora"},
			tgbotapi.BotCommand{Command: "golden", Description: "real time golden price"},
		))

		if err != nil {
			json.NewEncoder(w).Encode([]any{
				err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {	})
	workers.Serve(nil) // use http.DefaultServeMux
}
