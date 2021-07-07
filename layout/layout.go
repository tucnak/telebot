package layout

import (
	"bytes"
	"io/ioutil"
	"log"
	"sync"
	"text/template"

	"github.com/goccy/go-yaml"
	tele "gopkg.in/tucnak/telebot.v3"
)

type (
	// Layout provides an interface to interact with the layout,
	// parsed from the config file and locales.
	Layout struct {
		pref  *tele.Settings
		mu    sync.RWMutex // protects ctxs
		ctxs  map[tele.Context]string
		funcs template.FuncMap

		buttons map[string]Button
		markups map[string]Markup
		locales map[string]*template.Template

		*Config
	}

	// Button is a shortcut for tele.Btn.
	Button = tele.Btn

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
)

// New reads and parses the given layout file.
func New(path string) (*Layout, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lt := Layout{
		ctxs:  make(map[tele.Context]string),
		funcs: make(template.FuncMap),
	}

	for k, v := range funcs {
		lt.funcs[k] = v
	}

	return &lt, yaml.Unmarshal(data, &lt)
}

var funcs = template.FuncMap{
	// Built-in blank and helper functions.
	"locale": func() string { return "" },
	"config": func(string) string { return "" },
	"text":   func(string) string { return "" },
}

// AddFunc adds the given function to the template FuncMap.
// Note: to make it come into effect, always add functions before New().
func AddFunc(key string, fn interface{}) {
	funcs[key] = fn
}

// AddFuncs extends the template FuncMap with the given one.
// Note: to make it come into effect, always add functions before New().
func AddFuncs(fm template.FuncMap) {
	for k, v := range fm {
		funcs[k] = v
	}
}

// Settings returns built telebot Settings required for bot initialising.
//
//		settings:
//		  url: (custom url if needed)
//		  token: (not recommended)
//		  updates: (chan capacity)
//		  locales_dir: (optional)
//		  token_env: (token env var name, example: TOKEN)
// 		  parse_mode: (default parse mode)
// 		  long_poller: (long poller settings)
//		  webhook: (or webhook settings)
//
// Usage:
//		lt, err := layout.New("bot.yml")
//		b, err := tele.NewBot(lt.Settings())
//		// That's all!
//
func (lt *Layout) Settings() tele.Settings {
	if lt.pref == nil {
		panic("telebot/layout: settings is empty")
	}
	return *lt.pref
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

// Text returns a text, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
// Example of en.yml:
//		start: Hi, {{.FirstName}}!
//
// Usage:
//		func OnStart(c tele.Context) error {
//			return c.Send(lt.Text(c, "start", c.Sender()))
//		}
//
func (lt *Layout) Text(c tele.Context, k string, args ...interface{}) string {
	locale, ok := lt.Locale(c)
	if !ok {
		return ""
	}

	return lt.TextLocale(locale, k, args...)
}

// TextLocale returns a localised text processed with standard template engine.
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

// Callback returns casted to CallbackEndpoint button, which mostly
// useful for handlers registering.
//
// Example:
//
//		// Handling settings button
//		b.Handle(lt.Callback("settings"), OnSettings)
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
//		buttons:
//		  item:
//		    unique: item
//		    callback_data: {{.ID}}
//		    text: Item #{{.Number}}
//
// Usage:
//		btns := make([]tele.Btn, len(items))
//		for i, item := range items {
//			btns[i] = lt.Button(c, "item", struct {
//				Number int
//				Item   Item
//			}{
//				Number: i,
//				Item:   item,
//			})
//		}
//
//		m := b.NewMarkup()
//		m.Inline(m.Row(btns...))
//		// Your generated markup is ready.
//
func (lt *Layout) Button(c tele.Context, k string, args ...interface{}) *tele.Btn {
	locale, ok := lt.Locale(c)
	if !ok {
		return nil
	}

	return lt.ButtonLocale(locale, k, args...)
}

// ButtonLocale returns a localised button processed with standard template engine.
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

	tmpl, err := lt.template(template.New(k), locale).Parse(string(data))
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

	return &btn
}

// Markup returns a markup, which locale is dependent on the context.
// The given optional argument will be passed to the template engine.
//
//		buttons:
//		  settings: 'Settings'
//		markups:
//		  menu:
//		    - [settings]
//
// Usage:
//		func OnStart(c tele.Context) error {
//			return c.Send(
//				lt.Text(c, "start"),
//				lt.Markup(c, "menu"))
//		}
//
func (lt *Layout) Markup(c tele.Context, k string, args ...interface{}) *tele.ReplyMarkup {
	locale, ok := lt.Locale(c)
	if !ok {
		return nil
	}

	return lt.MarkupLocale(locale, k, args...)
}

// MarkupLocale returns a localised markup processed with standard template engine.
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

func (lt *Layout) template(tmpl *template.Template, locale string) *template.Template {
	funcs := make(template.FuncMap)

	// Redefining built-in blank functions
	funcs["config"] = lt.String
	funcs["text"] = func(k string) string { return lt.TextLocale(locale, k) }
	funcs["locale"] = func() string { return locale }

	return tmpl.Funcs(funcs)
}
