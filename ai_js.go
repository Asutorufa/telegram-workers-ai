package ai

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand/v2"
	"syscall/js"

	"github.com/syumai/go-jsutil"
	"github.com/syumai/workers/cloudflare"
)

type AI struct {
	instance js.Value
}

func NewAI() *AI {
	return &AI{
		instance: cloudflare.GetBinding("AI"),
	}
}

/*
	const inputs = {
	  prompt: "cyberpunk cat",
	};
*/
type DiffusionOptions struct {
	Prompt string
}

func (opts *DiffusionOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsutil.NewObject()
	obj.Set("prompt", opts.Prompt)
	return obj
}

func (a *AI) Diffusion(opt DiffusionOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/stabilityai/stable-diffusion-xl-base-1.0", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	r := js.Global().Get("Response").New(t).Get("body")
	return jsutil.ConvertReadableStreamToReadCloser(r), nil
}

func (a *AI) Dreamshaper(opt DiffusionOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/lykon/dreamshaper-8-lcm", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	r := js.Global().Get("Response").New(t).Get("body")
	return jsutil.ConvertReadableStreamToReadCloser(r), nil
}

func (a *AI) Flux1Schnell(opt DiffusionOptions) ([]byte, error) {
	p := a.instance.Call("run", "@cf/black-forest-labs/flux-1-schnell", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	respString := js.Global().Get("JSON").Call("stringify", t).String()

	type Response struct {
		Image string `json:"image"`
	}

	var resp Response
	if err := json.Unmarshal([]byte(respString), &resp); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Image)
}

type Llama2_7bChatOptions struct {
	Prompt string
}

func (opts *Llama2_7bChatOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	system := jsutil.NewObject()
	system.Set("role", "system")
	system.Set("content", "You are a helpful assistant. Don't print with markdown format.")

	obj := jsutil.NewObject()
	obj.Set("role", "user")
	obj.Set("content", opts.Prompt)

	x := jsutil.NewObject()
	x.Set("messages", js.ValueOf([]any{system, obj}))
	// x.Set("prompt", opts.Prompt)
	x.Set("max_tokens", 1024)
	x.Set("seed", rand.IntN(999999))
	x.Set("stream", true)
	x.Set("temperature", 2.5)
	x.Set("top_k", 10)

	return x
}

// @cf/qwen/qwq-32b
func (a *AI) QwQ32B(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/qwen/qwq-32b", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return jsutil.ConvertReadableStreamToReadCloser(t), nil
}

// @cf/google/gemma-3-12b-it
func (a *AI) Gemma3_12b(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/google/gemma-3-12b-it", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return jsutil.ConvertReadableStreamToReadCloser(t), nil
}

// @cf/meta/llama-4-scout-17b-16e-instruct
func (a *AI) LLama4Scout17b16eInstruct(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/meta/llama-4-scout-17b-16e-instruct", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return jsutil.ConvertReadableStreamToReadCloser(t), nil
}

func (a *AI) Llama3_8bInstruct(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/meta/llama-3-8b-instruct", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return jsutil.ConvertReadableStreamToReadCloser(t), nil
}

func (a *AI) Mistral7bInstructV02Lora(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	// mistral-7b-instruct-v0.2-lor
	p := a.instance.Call("run", "@cf/mistral/mistral-7b-instruct-v0.2-lora", opt.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return jsutil.ConvertReadableStreamToReadCloser(t), nil
}

func (a *AI) Translate(opts TranslateOptions) (string, error) {
	p := a.instance.Call("run", "@cf/meta/m2m100-1.2b", opts.toJS())

	t, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}

	return t.Get("translated_text").String(), nil
}

/*
"@cf/meta/m2m100-1.2b",

	{
	  text: "I'll have an order of the moule frites",
	  source_lang: "english", // defaults to english
	  target_lang: "french",
	}


	​​Response

	{
	  "translated_text": "Je vais commander des moules frites"
	}
*/
type TranslateOptions struct {
	Text       string
	SourceLang string
	TargetLang string
}

func (opts *TranslateOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsutil.NewObject()
	if opts.Text != "" {
		obj.Set("text", opts.Text)
	}
	if opts.SourceLang != "" {
		obj.Set("source_lang", opts.SourceLang)
	}
	if opts.TargetLang != "" {
		obj.Set("target_lang", opts.TargetLang)
	}
	return obj
}
