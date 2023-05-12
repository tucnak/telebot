package layout

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"text/template"

	"github.com/goccy/go-yaml"
	tele "gopkg.in/telebot.v3"
)

type (
	// Layout provides an interface to interact with the layout,
	// parsed from the config file and locales.
	Layout struct {
		pref  *tele.Settings
		mu    sync.RWMutex // protects ctxs
		ctxs  map[tele.Context]string
		funcs template.FuncMap

		commands map[string]string
		buttons  map[string]Button
		markups  map[string]Markup
		results  map[string]Result
		locales  map[string]*template.Template

		Config
	}

	// Button is a shortcut for tele.Btn.
	Button struct {
		tele.Btn `yaml:",inline"`
		Data     interface{} `yaml:"data"`
		IsReply  bool        `yaml:"reply"`
	}

	// Markup represents layout-specific markup to be parsed.
	Markup struct {
		inline          *bool
		keyboard        *template.Template
		ResizeKeyboard  *bool `yaml:"resize_keyboard,omitempty"` // nil == true
		ForceReply      bool  `yaml:"force_reply,omitempty"`
		OneTimeKeyboard bool  `yaml:"one_time_keyboard,omitempty"`
		RemoveKeyboard  bool  `yaml:"remove_keyboard,omitempty"`
		Selective       bool  `yaml:"selective,omitempty"`
	}

	// Result represents layout-specific result to be parsed.
	Result struct {
		result          *template.Template
		tele.ResultBase `yaml:",inline"`
		Content         ResultContent `yaml:"content"`
		Markup          string        `yaml:"markup"`
	}

	// ResultBase represents layout-specific result's base to be parsed.
	ResultBase struct {
		tele.ResultBase `yaml:",inline"`
		Content         ResultContent `yaml:"content"`
	}

	// ResultContent represents any kind of InputMessageContent and implements it.
	ResultContent map[string]interface{}
)

// New parses the given layout file.
func New(path string, funcs ...template.FuncMap) (*Layout, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lt := Layout{
		ctxs:  make(map[tele.Context]string),
		funcs: make(template.FuncMap),
	}

	for k, v := range builtinFuncs {
		lt.funcs[k] = v
	}
	for i := range funcs {
		for k, v := range funcs[i] {
			lt.funcs[k] = v
		}
	}

	return &lt, yaml.Unmarshal(data, &lt)
}

// NewDefault parses the given layout file without localization features.
// See Layout.Default for more details.
func NewDefault(path, locale string, funcs ...template.FuncMap) (*DefaultLayout, error) {
	lt, err := New(path, funcs...)
	if err != nil {
		return nil, err
	}
	return lt.Default(locale), nil
}

var builtinFuncs = template.FuncMap{
	// Built-in blank and helper functions.
	"locale": func() string { return "" },
	"config": func(string) string { return "" },
	"text":   func(string, ...interface{}) string { return "" },
}

// Settings returns built telebot Settings required for bot initializing.
//
//	settings:
//		url: (custom url if needed)
//		token: (not recommended)
//		updates: (chan capacity)
//		locales_dir: (optional)
//		token_env: (token env var name, example: TOKEN)
//		parse_mode: (default parse mode)
//		long_poller: (long poller settings)
//		webhook: (or webhook settings)
//
// Usage:
//	lt, err := layout.New("bot.yml")
//	b, err := tele.NewBot(lt.Settings())
//	// That's all!
//
func (lt *Layout) Settings() tele.Settings {
	if lt.pref == nil {
		panic("telebot/layout: settings is empty")
	}
	return *lt.pref
}

// Default returns a simplified layout instance with the pre-defined locale.
// It's useful when you have no need for localization and don't want to pass
// context each time you use layout functions.
func (lt *Layout) Default(locale string) *DefaultLayout {
	return &DefaultLayout{
		locale: locale,
		lt:     lt,
		Config: lt.Config,
	}
}

// Locales returns all presented locales.
func (lt *Layout) Locales() []string {
	var keys []string
	for k := range lt.locales {
		keys = append(keys, k)
	}
	return keys
}

// Locale returns the context locale.
func (lt *Layout) Locale(c tele.Context) (string, bool) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	locale, ok := lt.ctxs[c]
	return locale, ok
}

// SetLocale allows you to change a locale for the passed context.
func (lt *Layout) SetLocale(c tele.Context, locale string) {
	lt.mu.Lock()
	lt.ctxs[c] = locale
	lt.mu.Unlock()
}

// Commands returns a list of telebot commands, which can be
// used in b.SetCommands later.
func (lt *Layout) Commands() (cmds []tele.Command) {
	for k, v := range lt.commands {
		cmds = append(cmds, tele.Command{
			Text:        strings.TrimLeft(k, "/"),
			Description: v,
		})
	}
	return
}

// CommandsLocale returns a list of telebot commands and localized description, which can be
// used in b.SetCommands later.
//
// Example of bot.yml:
//	commands:
//	  /start: '{{ text `cmdStart` }}'
//
// en.yml:
//	cmdStart: Start the bot
//
// ru.yml:
//	cmdStart: Запуск бота
//
// Usage:
//	b.SetCommands(lt.CommandsLocale("en"), "en")
//	b.SetCommands(lt.CommandsLocale("ru"), "ru")
func (lt *Layout) CommandsLocale(locale string, args ...interface{}) (cmds []tele.Command) {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	for k, v := range lt.commands {
		tmpl, err := lt.template(template.New(k).Funcs(lt.funcs), locale).Parse(v)
		if err != nil {
			log.Println("telebot/layout:", err)
			return nil
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, arg); err != nil {
			log.Println("telebot/layout:", err)
			return nil
		}

		cmds = append(cmds, tele.Command{
			Text:        strings.TrimLeft(k, "/"),
			Description: buf.String(),
		})
	}
	return
}

// Text returns a text, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
// Example of en.yml:
//	start: Hi, {{.FirstName}}!
//
// Usage:
//	func onStart(c tele.Context) error {
//		return c.Send(lt.Text(c, "start", c.Sender()))
//	}
//
func (lt *Layout) Text(c tele.Context, k string, args ...interface{}) string {
	locale, ok := lt.Locale(c)
	if !ok {
		return ""
	}

	return lt.TextLocale(locale, k, args...)
}

// TextLocale returns a localized text processed with text/template engine.
// See Text for more details.
func (lt *Layout) TextLocale(locale, k string, args ...interface{}) string {
	tmpl, ok := lt.locales[locale]
	if !ok {
		return ""
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	if err := lt.template(tmpl, locale).ExecuteTemplate(&buf, k, arg); err != nil {
		log.Println("telebot/layout:", err)
	}

	return buf.String()
}

// Callback returns a callback endpoint used to handle buttons.
//
// Example:
//	// Handling settings button
//	b.Handle(lt.Callback("settings"), onSettings)
//
func (lt *Layout) Callback(k string) tele.CallbackEndpoint {
	btn, ok := lt.buttons[k]
	if !ok {
		return nil
	}
	return &btn
}

// Button returns a button, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
//	buttons:
//		item:
//			unique: item
//			callback_data: {{.ID}}
//			text: Item #{{.Number}}
//
// Usage:
//	btns := make([]tele.Btn, len(items))
//	for i, item := range items {
//		btns[i] = lt.Button(c, "item", struct {
//			Number int
//			Item   Item
//		}{
//			Number: i,
//			Item:   item,
//		})
//	}
//
//	m := b.NewMarkup()
//	m.Inline(m.Row(btns...))
//	// Your generated markup is ready.
//
func (lt *Layout) Button(c tele.Context, k string, args ...interface{}) *tele.Btn {
	locale, ok := lt.Locale(c)
	if !ok {
		return nil
	}

	return lt.ButtonLocale(locale, k, args...)
}

// ButtonLocale returns a localized button processed with text/template engine.
// See Button for more details.
func (lt *Layout) ButtonLocale(locale, k string, args ...interface{}) *tele.Btn {
	btn, ok := lt.buttons[k]
	if !ok {
		return nil
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	data, err := yaml.Marshal(btn)
	if err != nil {
		log.Println("telebot/layout:", err)
		return nil
	}

	tmpl, err := lt.template(template.New(k).Funcs(lt.funcs), locale).Parse(string(data))
	if err != nil {
		log.Println("telebot/layout:", err)
		return nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, arg); err != nil {
		log.Println("telebot/layout:", err)
		return nil
	}

	if err := yaml.Unmarshal(buf.Bytes(), &btn); err != nil {
		log.Println("telebot/layout:", err)
		return nil
	}

	return &btn.Btn
}

// Markup returns a markup, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
//	buttons:
//		settings: 'Settings'
//	markups:
//		menu:
//		- [settings]
//
// Usage:
//	func onStart(c tele.Context) error {
//		return c.Send(
//			lt.Text(c, "start"),
//			lt.Markup(c, "menu"),
//		)
//	}
//
func (lt *Layout) Markup(c tele.Context, k string, args ...interface{}) *tele.ReplyMarkup {
	locale, ok := lt.Locale(c)
	if !ok {
		return nil
	}

	return lt.MarkupLocale(locale, k, args...)
}

// MarkupLocale returns a localized markup processed with text/template engine.
// See Markup for more details.
func (lt *Layout) MarkupLocale(locale, k string, args ...interface{}) *tele.ReplyMarkup {
	markup, ok := lt.markups[k]
	if !ok {
		return nil
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	if err := lt.template(markup.keyboard, locale).Execute(&buf, arg); err != nil {
		log.Println("telebot/layout:", err)
	}

	r := &tele.ReplyMarkup{}
	if *markup.inline {
		if err := yaml.Unmarshal(buf.Bytes(), &r.InlineKeyboard); err != nil {
			log.Println("telebot/layout:", err)
		}
	} else {
		r.ResizeKeyboard = markup.ResizeKeyboard == nil || *markup.ResizeKeyboard
		r.ForceReply = markup.ForceReply
		r.OneTimeKeyboard = markup.OneTimeKeyboard
		r.RemoveKeyboard = markup.RemoveKeyboard
		r.Selective = markup.Selective

		if err := yaml.Unmarshal(buf.Bytes(), &r.ReplyKeyboard); err != nil {
			log.Println("telebot/layout:", err)
		}
	}

	return r
}

// Result returns an inline result, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
//	results:
//		article:
//			type: article
//			id: '{{ .ID }}'
//			title: '{{ .Title }}'
//			description: '{{ .Description }}'
//			message_text: '{{ .Content }}'
//			thumb_url: '{{ .PreviewURL }}'
//
// Usage:
//	func onQuery(c tele.Context) error {
//		results := make(tele.Results, len(articles))
//		for i, article := range articles {
//			results[i] = lt.Result(c, "article", article)
//		}
//		return c.Answer(&tele.QueryResponse{
//			Results:   results,
//			CacheTime: 100,
//		})
//	}
//
func (lt *Layout) Result(c tele.Context, k string, args ...interface{}) tele.Result {
	locale, ok := lt.Locale(c)
	if !ok {
		return nil
	}

	return lt.ResultLocale(locale, k, args...)
}

// ResultLocale returns a localized result processed with text/template engine.
// See Result for more details.
func (lt *Layout) ResultLocale(locale, k string, args ...interface{}) tele.Result {
	result, ok := lt.results[k]
	if !ok {
		return nil
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	if err := lt.template(result.result, locale).Execute(&buf, arg); err != nil {
		log.Println("telebot/layout:", err)
	}

	var (
		data = buf.Bytes()
		base Result
		r    tele.Result
	)

	if err := yaml.Unmarshal(data, &base); err != nil {
		log.Println("telebot/layout:", err)
	}

	switch base.Type {
	case "article":
		r = &tele.ArticleResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "audio":
		r = &tele.AudioResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "contact":
		r = &tele.ContactResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "document":
		r = &tele.DocumentResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "gif":
		r = &tele.GifResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "location":
		r = &tele.LocationResult{ResultBase: base.ResultBase}
		if err := json.Unmarshal(data, &r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "mpeg4_gif":
		r = &tele.Mpeg4GifResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "photo":
		r = &tele.PhotoResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "venue":
		r = &tele.VenueResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "video":
		r = &tele.VideoResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "voice":
		r = &tele.VoiceResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	case "sticker":
		r = &tele.StickerResult{ResultBase: base.ResultBase}
		if err := yaml.Unmarshal(data, r); err != nil {
			log.Println("telebot/layout:", err)
		}
	default:
		log.Println("telebot/layout: unsupported inline result type")
		return nil
	}

	if base.Content != nil {
		r.SetContent(base.Content)
	}

	if result.Markup != "" {
		markup := lt.MarkupLocale(locale, result.Markup, args...)
		if markup == nil {
			log.Printf("telebot/layout: markup with name %s was not found\n", result.Markup)
		} else {
			r.SetReplyMarkup(markup)
		}
	}

	return r
}

func (lt *Layout) template(tmpl *template.Template, locale string) *template.Template {
	funcs := make(template.FuncMap)

	// Redefining built-in blank functions
	funcs["config"] = lt.String
	funcs["text"] = func(k string, args ...interface{}) string { return lt.TextLocale(locale, k, args...) }
	funcs["locale"] = func() string { return locale }

	return tmpl.Funcs(funcs)
}

// IsInputMessageContent implements telebot.InputMessageContent.
func (ResultContent) IsInputMessageContent() bool {
	return true
}
