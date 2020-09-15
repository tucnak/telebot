package layout

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	tele "gopkg.in/tucnak/telebot.v3"
)

type (
	Layout struct {
		pref   *tele.Settings
		ctxs   map[tele.Context]string
		locale string

		Config  map[string]interface{}
		Markups map[string]Markup
		Locales map[string]*template.Template
	}

	Markup struct {
		tele.ReplyMarkup `yaml:",inline"`
		Keyboard         *template.Template `yaml:"-"`
		inline           bool
	}

	Button struct {
		tele.ReplyButton  `yaml:",inline"`
		tele.InlineButton `yaml:",inline"`
	}

	LocaleFunc func(tele.Recipient) string
)

func New(path string) (*Layout, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lt := Layout{ctxs: make(map[tele.Context]string)}
	return &lt, yaml.Unmarshal(data, &lt)
}

func (lt *Layout) Settings() tele.Settings {
	if lt.pref == nil {
		panic("telebot/layout: settings is empty")
	}
	return *lt.pref
}

func (lt *Layout) With(c tele.Context) *Layout {
	cp := *lt
	cp.locale = lt.ctxs[c]
	return &cp
}

func (lt *Layout) Get(k string) string {
	return fmt.Sprint(lt.Config[k])
}

func (lt *Layout) Int(k string) int {
	return cast.ToInt(lt.Config[k])
}

func (lt *Layout) Int64(k string) int64 {
	return cast.ToInt64(lt.Config[k])
}

func (lt *Layout) Float(k string) float64 {
	return cast.ToFloat64(lt.Config[k])
}

func (lt *Layout) Duration(k string) time.Duration {
	return cast.ToDuration(lt.Config[k])
}

func (lt *Layout) Text(k string, args ...interface{}) string {
	if len(lt.Locales) == 0 {
		return ""
	}

	tmpl, ok := lt.Locales[lt.locale]
	if !ok {
		return ""
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, k, arg); err != nil {
		// TODO: Log.
	}
	return buf.String()
}

func (lt *Layout) Markup(k string, args ...interface{}) *tele.ReplyMarkup {
	if len(lt.Markups) == 0 {
		return nil
	}

	markup, ok := lt.Markups[k]
	if !ok {
		return nil
	}

	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	var buf bytes.Buffer
	if err := markup.Keyboard.Execute(&buf, arg); err != nil {
		// TODO: Log.
	}

	r := tele.ReplyMarkup{
		ForceReply:          markup.ForceReply,
		ResizeReplyKeyboard: markup.ResizeReplyKeyboard,
		OneTimeKeyboard:     markup.OneTimeKeyboard,
		ReplyKeyboardRemove: markup.ReplyKeyboardRemove,
		Selective:           markup.Selective,
	}

	if markup.inline {
		if err := yaml.Unmarshal(buf.Bytes(), &r.InlineKeyboard); err != nil {
			// TODO: Log.
		}
	} else {
		if err := yaml.Unmarshal(buf.Bytes(), &r.ReplyKeyboard); err != nil {
			// TODO: Log.
		}
	}

	return &r
}
