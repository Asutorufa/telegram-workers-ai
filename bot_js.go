package ai

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/syumai/tinyutil/httputil"
	"github.com/syumai/workers/cloudflare"
)

var Bot = &tgbotapi.BotAPI{
	Token:  cloudflare.Getenv("telegram_token"),
	Client: httputil.DefaultClient,
	Buffer: 100,
}

var MasterID func() int64 = func() int64 {
	for _, id := range strings.FieldsFunc(cloudflare.Getenv("telegram_ids"), func(r rune) bool { return r == ',' }) {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			continue
		}

		return i
	}

	return 0
}

func init() {
	Bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
}

func BotHandler() func(w http.ResponseWriter, r *http.Request) {
	idMap := make(map[int64]bool)
	for _, id := range strings.FieldsFunc(cloudflare.Getenv("telegram_ids"), func(r rune) bool { return r == ',' }) {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		idMap[i] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		update, err := Bot.HandleUpdate(r)
		if err != nil {
			enc := jsontext.NewEncoder(w)
			json.MarshalEncode(enc, []any{err.Error()})
			return
		}

		if update.Message == nil || (idMap != nil && !idMap[update.Message.From.ID]) {
			return
		}

		argument := update.Message.CommandArguments()

		// If we got a message
		log.Printf("[%s] %s argument: %s\n", update.Message.From.UserName, update.Message.Text, argument)

		var msg tgbotapi.Chattable
		switch update.Message.Command() {
		case "image":
			if argument == "" {
				return
			}
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
		case "dreamshaper":
			if argument == "" {
				return
			}
			data, err := NewAI().Dreamshaper(DiffusionOptions{Prompt: argument})
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

		case "phoenix":
			if argument == "" {
				return
			}
			data, err := NewAI().Phoenix(DiffusionOptions{Prompt: argument})
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

		case "flux1schnell":
			if argument == "" {
				return
			}
			data, err := NewAI().Flux1Schnell(DiffusionOptions{Prompt: argument})
			if err != nil {
				m := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			} else {
				m := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileReader{Name: "image.png", Reader: bytes.NewReader(data)})
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			}

		case "lucid":
			if argument == "" {
				return
			}
			data, err := NewAI().LucidOrigin(DiffusionOptions{Prompt: argument})
			if err != nil {
				m := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			} else {
				m := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileReader{Name: "image.png", Reader: bytes.NewReader(data)})
				m.ReplyToMessageID = update.Message.MessageID
				msg = m
			}

		case "user_id":
			m := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprint(update.Message.From.ID))
			m.ReplyToMessageID = update.Message.MessageID
			msg = m

		case "qwq32b":
			if argument == "" {
				return
			}
			data, err := NewAI().QwQ32B(Llama2_7bChatOptions{Prompt: argument})
			if err != nil {
				data = io.NopCloser(strings.NewReader(err.Error()))
			}
			defer data.Close()

			ReturnByEventSource(data, update)
			return

		case "gemma312b":
			if argument == "" {
				return
			}
			data, err := NewAI().Gemma3_12b(Llama2_7bChatOptions{Prompt: argument})
			if err != nil {
				data = io.NopCloser(strings.NewReader(err.Error()))
			}
			defer data.Close()

			ReturnByEventSource(data, update)
			return

		case "llama4scout17b16e":
			if argument == "" {
				return
			}
			data, err := NewAI().LLama4Scout17b16eInstruct(Llama2_7bChatOptions{Prompt: argument})
			if err != nil {
				data = io.NopCloser(strings.NewReader(err.Error()))
			}
			defer data.Close()
			ReturnByEventSource(data, update)
			return

		case "llama38b":
			if argument == "" {
				return
			}
			data, err := NewAI().Llama3_8bInstruct(Llama2_7bChatOptions{Prompt: argument})
			if err != nil {
				data = io.NopCloser(strings.NewReader(err.Error()))
			}
			defer data.Close()

			ReturnByEventSource(data, update)
			return

		case "mistral7b":
			if argument == "" {
				return
			}
			data, err := NewAI().Mistral7bInstructV02Lora(Llama2_7bChatOptions{Prompt: argument})
			if err != nil {
				data = io.NopCloser(strings.NewReader(err.Error()))
			}
			defer data.Close()

			ReturnByEventSource(data, update)
			return

		case "golden", "jpycny", "usdcny":
			var str string
			switch update.Message.Command() {
			case "golden":
				str = Golden()
			case "jpycny":
				str = JPYCNY()
				log.Println(str)
			case "usdcny":
				str = USDCNY()
			default:
				return
			}
			m := tgbotapi.NewMessage(update.Message.Chat.ID, str)
			m.ReplyToMessageID = update.Message.MessageID
			m.ParseMode = "MarkdownV2"
			msg = m

		default:
			return
		}

		_, err = Bot.Send(msg)
		if err != nil {
			log.Println("send msg error", err)
			return
		}
	}
}

func ReturnByEventSource(r io.ReadCloser, update *tgbotapi.Update) {
	br := NewLlamaStreamDecoder(r)

	text := strings.Builder{}
	last := 0
	msgId := 0
	count := 0
	for {
		e, err := br.Decode()
		if err != nil {
			if err != io.EOF {
				log.Println("decode error", err)
			}
			break
		}
		if e == "" {
			continue
		}

		text.WriteString(e)

		if count >= 12 || text.Len()-last <= 300 {
			continue
		}

		msg, err := SendText(update, msgId, text.String())
		if err != nil {
			log.Println("send msg error", "err", err)
			return
		}

		count++
		last = text.Len()

		if msgId == 0 {
			msgId = msg.MessageID
		}
	}

	SendText(update, msgId, text.String())
}

func SendText(update *tgbotapi.Update, msgId int, text string) (tgbotapi.Message, error) {
	var msg tgbotapi.Chattable
	if msgId == 0 {
		m := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		m.ReplyToMessageID = update.Message.MessageID
		msg = m
	} else {
		msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, msgId, text)
	}

	rm, err := Bot.Send(msg)
	if err != nil {
		return rm, err
	}

	return rm, nil
}
