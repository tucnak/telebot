package layout

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	tele "gopkg.in/tucnak/telebot.v3"
)

type (
	Layout struct {
		pref  *tele.Settings
		mu    sync.RWMutex // protects ctxs
		ctxs  map[tele.Context]string
		funcs template.FuncMap

		config  map[string]interface{}
		buttons map[string]Button
		markups map[string]Markup
		locales map[string]*template.Template
	}

	Button = tele.Btn

	Markup struct {
		inline          *bool
		keyboard        *template.Template
		ResizeKeyboard  *bool `json:"resize_keyboard,omitempty"` // nil == true
		ForceReply      bool  `json:"force_reply,omitempty"`
		OneTimeKeyboard bool  `json:"one_time_keyboard,omitempty"`
		RemoveKeyboard  bool  `json:"remove_keyboard,omitempty"`
		Selective       bool  `json:"selective,omitempty"`
	}
)

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

	// Built-in blank and helper functions
	lt.funcs["config"] = lt.Get
	lt.funcs["locale"] = func() string { return "" }
	lt.funcs["text"] = func(k string) string { return "" }

	return &lt, yaml.Unmarshal(data, &lt)
}

var funcs = make(template.FuncMap)

func AddFunc(key string, fn interface{}) {
	funcs[key] = fn
}

func AddFuncs(fm template.FuncMap) {
	for k, v := range fm {
		funcs[k] = v
	}
}

func (lt *Layout) Settings() tele.Settings {
	if lt.pref == nil {
		panic("telebot/layout: settings is empty")
	}
	return *lt.pref
}

func (lt *Layout) Get(k string) string {
	return fmt.Sprint(lt.config[k])
}

func (lt *Layout) Int(k string) int {
	return cast.ToInt(lt.config[k])
}

func (lt *Layout) Int64(k string) int64 {
	return cast.ToInt64(lt.config[k])
}

func (lt *Layout) Float(k string) float64 {
	return cast.ToFloat64(lt.config[k])
}

func (lt *Layout) Duration(k string) time.Duration {
	return cast.ToDuration(lt.config[k])
}

func (lt *Layout) Text(c tele.Context, k string, args ...interface{}) string {
	locale, ok := lt.Locale(c)
	if !ok {
		return ""
	}

	return lt.text(locale, k, args...)
}

func (lt *Layout) text(locale, k string, args ...interface{}) string {
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
		// TODO: Log.
	}

	return buf.String()
}

func (lt *Layout) Button(k string) tele.CallbackEndpoint {
	btn, ok := lt.buttons[k]
	if !ok {
		return nil
	}
	return &btn
}

func (lt *Layout) Markup(c tele.Context, k string, args ...interface{}) *tele.ReplyMarkup {
	markup, ok := lt.markups[k]
	if !ok {
		return nil
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	locale, _ := lt.Locale(c)
	if err := lt.template(markup.keyboard, locale).Execute(&buf, arg); err != nil {
		// TODO: Log.
	}

	r := &tele.ReplyMarkup{}
	if *markup.inline {
		if err := yaml.Unmarshal(buf.Bytes(), &r.InlineKeyboard); err != nil {
			// TODO: Log.
		}
	} else {
		r.ResizeKeyboard = markup.ResizeKeyboard == nil || *markup.ResizeKeyboard
		r.ForceReply = markup.ForceReply
		r.OneTimeKeyboard = markup.OneTimeKeyboard
		r.RemoveKeyboard = markup.RemoveKeyboard
		r.Selective = markup.Selective

		if err := yaml.Unmarshal(buf.Bytes(), &r.ReplyKeyboard); err != nil {
			// TODO: Log.
		}
	}

	return r
}

func (lt *Layout) template(tmpl *template.Template, locale string) *template.Template {
	funcs := make(template.FuncMap)

	// Redefining built-in blank functions
	funcs["text"] = func(k string) string { return lt.text(locale, k) }
	funcs["locale"] = func() string { return locale }

	return tmpl.Funcs(funcs)
}

func (lt *Layout) SetLocale(c tele.Context, locale string) {
	lt.mu.Lock()
	lt.ctxs[c] = locale
	lt.mu.Unlock()
}

func (lt *Layout) Locale(c tele.Context) (string, bool) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	locale, ok := lt.ctxs[c]
	return locale, ok
}
