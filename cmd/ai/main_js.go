package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall/js"

	ai "github.com/Asutorufa/telegram-workers-ai"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func main() {
	workers.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("new request", "method", r.Method, "path", r.URL.Path)

		enc := jsontext.NewEncoder(w)

		switch r.URL.Path {
		case "/tgbot":
			ai.BotHandler()(w, r)
		case "/tgbot/register":
			domain := "https://" + r.URL.Host
			if x := cloudflare.GetBinding("worker_url"); x.Type() == js.TypeString {
				domain = x.String()
			}

			wh, err := tgbotapi.NewWebhook(domain + "/tgbot")
			if err != nil {
				json.MarshalEncode(enc, []any{err.Error()})
				return
			}

			_, err = ai.Bot.Request(wh)
			if err != nil {
				json.MarshalEncode(enc, []any{err.Error()})
				return
			}

			resp, err := ai.Bot.Request(tgbotapi.NewSetMyCommands(
				tgbotapi.BotCommand{Command: "image", Description: "generate image by prompt"},
				tgbotapi.BotCommand{Command: "dreamshaper", Description: "generate image by @cf/lykon/dreamshaper-8-lcm"},
				tgbotapi.BotCommand{Command: "phoenix", Description: "generate image by @cf/leonardo/phoenix-1.0"},
				tgbotapi.BotCommand{Command: "llama38b", Description: "generate text by @cf/meta/llama-3-8b-instruct"},
				tgbotapi.BotCommand{Command: "qwq32b", Description: "generate text by @cf/qwen/qwq-32b"},
				tgbotapi.BotCommand{Command: "gemma312b", Description: "generate text by @cf/google/gemma-3-12b-it"},
				tgbotapi.BotCommand{Command: "llama4scout17b16e", Description: "generate text by @cf/meta/llama-4-scout-17b-16e-instruct"},
				tgbotapi.BotCommand{Command: "mistral7b", Description: "generate text by mistral-7b-instruct-v0.2-lora"},
				tgbotapi.BotCommand{Command: "flux1schnell", Description: "generate image by @cf/black-forest-labs/flux-1-schnell"},
				tgbotapi.BotCommand{Command: "lucid", Description: "generate image by @cf/leonardo/lucid-origin"},
				tgbotapi.BotCommand{Command: "golden", Description: "real time golden price"},
				tgbotapi.BotCommand{Command: "jpycny", Description: "real time jpy to cny price"},
				tgbotapi.BotCommand{Command: "usdcny", Description: "real time usd to cny price"},
			))

			if err != nil {
				enc := jsontext.NewEncoder(w)
				json.MarshalEncode(enc, []any{err.Error()})
				return
			}

			_, err = ai.Bot.Send(tgbotapi.NewMessage(ai.MasterID(), fmt.Sprintf("register webhook to %s", domain)))
			if err != nil {
				log.Println(err)
			}

			json.MarshalEncode(enc, resp)
		}
	})) // use http.DefaultServeMux
}
