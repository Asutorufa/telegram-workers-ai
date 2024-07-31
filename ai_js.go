package ai

import (
	"io"
	"math/rand/v2"
	"syscall/js"
	_ "unsafe"

	_ "github.com/syumai/workers"
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
	obj := NewObject()
	obj.Set("prompt", opts.Prompt)
	obj.Set("seed", rand.IntN(999999))
	return obj
}

func (a *AI) Diffusion(opt DiffusionOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/stabilityai/stable-diffusion-xl-base-1.0", opt.toJS())

	t, err := AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	r := js.Global().Get("Response").New(t).Get("body")
	return ConvertReadableStreamToReadCloser(r)
}

type Llama2_7bChatOptions struct {
	Prompt string
}

func (opts *Llama2_7bChatOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	// obj := NewObject()
	// obj.Set("role", "user")
	// obj.Set("content", opts.Prompt)

	x := NewObject()
	// x.Set("messages", js.ValueOf([]any{obj}))
	x.Set("prompt", opts.Prompt)
	x.Set("max_tokens", 512)
	x.Set("seed", rand.IntN(999999))
	x.Set("stream", true)

	return x
}

func (a *AI) Llama3_8bInstruct(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	p := a.instance.Call("run", "@cf/meta/llama-3-8b-instruct", opt.toJS())

	t, err := AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return ConvertReadableStreamToReadCloser(t)
}

func (a *AI) Mistral7bInstructV02Lora(opt Llama2_7bChatOptions) (io.ReadCloser, error) {
	// mistral-7b-instruct-v0.2-lor
	p := a.instance.Call("run", "@cf/mistral/mistral-7b-instruct-v0.2-lora", opt.toJS())

	t, err := AwaitPromise(p)
	if err != nil {
		return nil, err
	}

	return ConvertReadableStreamToReadCloser(t)
}

func (a *AI) Translate(opts TranslateOptions) (string, error) {
	p := a.instance.Call("run", "@cf/meta/m2m100-1.2b", opts.toJS())

	t, err := AwaitPromise(p)
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
	obj := NewObject()
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

//go:linkname NewObject github.com/syumai/workers/internal/jsutil.NewObject
func NewObject() js.Value

//go:linkname AwaitPromise github.com/syumai/workers/internal/jsutil.AwaitPromise
func AwaitPromise(promiseVal js.Value) (js.Value, error)

//go:linkname ConvertReadableStreamToReadCloser github.com/syumai/workers/internal/jsutil.ConvertReadableStreamToReadCloser
func ConvertReadableStreamToReadCloser(stream js.Value) (io.ReadCloser, error)

//go:linkname ArrayFrom github.com/syumai/workers/internal/jsutil.ArrayFrom
func ArrayFrom(v js.Value) js.Value
