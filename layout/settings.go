package layout

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goccy/go-yaml"
	tele "gopkg.in/tucnak/telebot.v3"
)

type (
	Settings struct {
		URL     string
		Token   string
		Updates int

		LocalesDir string `json:"locales_dir"`
		TokenEnv   string `json:"token_env"`
		ParseMode  string `json:"parse_mode"`

		Webhook    *tele.Webhook    `json:"webhook"`
		LongPoller *tele.LongPoller `json:"long_poller"`
	}
)

func (lt *Layout) UnmarshalYAML(data []byte) error {
	var aux struct {
		Settings *Settings
		Config   map[string]interface{}
		Markups  yaml.MapSlice
		Locales  map[string]map[string]string
	}
	if err := yaml.Unmarshal(data, &aux); err != nil {
		return err
	}

	lt.Config = aux.Config

	if pref := aux.Settings; pref != nil {
		lt.pref = &tele.Settings{
			URL:       pref.URL,
			Token:     pref.Token,
			Updates:   pref.Updates,
			ParseMode: pref.ParseMode,
		}

		if pref.LocalesDir == "" {
			aux.Settings.LocalesDir = "locales"
		}
		if pref.TokenEnv != "" {
			lt.pref.Token = os.Getenv(pref.TokenEnv)
		}

		if pref.Webhook != nil {
			lt.pref.Poller = pref.Webhook
		} else if pref.LongPoller != nil {
			lt.pref.Poller = pref.LongPoller
		}
	}

	lt.Markups = make(map[string]Markup, len(aux.Markups))
	for _, item := range aux.Markups {
		k, v := item.Key.(string), item.Value

		data, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		// 1. Normal markup.

		var markup struct {
			Markup `yaml:",inline"`
			Resize *bool `json:"resize_keyboard"`
		}
		if yaml.Unmarshal(data, &markup) == nil {
			data, err := yaml.Marshal(markup.ReplyKeyboard)
			if err != nil {
				return err
			}

			tmpl, err := template.New(k).Parse(string(data))
			if err != nil {
				return err
			}

			markup.Markup.Keyboard = tmpl
			markup.ResizeReplyKeyboard = markup.Resize == nil || *markup.Resize

			lt.Markups[k] = markup.Markup
		}

		// 2. Shortened reply markup.

		var embeddedMarkup [][]string
		if yaml.Unmarshal(data, &embeddedMarkup) == nil {
			kb := make([][]tele.ReplyButton, len(embeddedMarkup))
			for i, btns := range embeddedMarkup {
				row := make([]tele.ReplyButton, len(btns))
				for j, btn := range btns {
					row[j] = tele.ReplyButton{Text: btn}
				}
				kb[i] = row
			}

			data, err := yaml.Marshal(kb)
			if err != nil {
				return err
			}

			tmpl, err := template.New(k).Parse(string(data))
			if err != nil {
				return err
			}

			markup := Markup{Keyboard: tmpl}
			markup.ResizeReplyKeyboard = true
			lt.Markups[k] = markup
		}

		// 3. Shortened inline markup.

		if yaml.Unmarshal(data, &[][]tele.InlineButton{}) == nil {
			tmpl, err := template.New(k).Parse(string(data))
			if err != nil {
				return err
			}

			lt.Markups[k] = Markup{
				inline:   true,
				Keyboard: tmpl,
			}
		}
	}

	if aux.Locales == nil {
		return lt.parseLocales(aux.Settings.LocalesDir)
	}

	return nil
}

func (lt *Layout) parseLocales(dir string) error {
	lt.Locales = make(map[string]*template.Template)

	return filepath.Walk(dir, func(path string, fi os.FileInfo, _ error) error {
		if fi == nil || fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var texts map[string]string
		if err := yaml.Unmarshal(data, &texts); err != nil {
			return err
		}

		name := filepath.Base(fi.Name())
		tmpl := template.New(name)

		for _, text := range texts {
			t, err := tmpl.Parse(text)
			if err != nil {
				return err
			}
			tmpl = t
		}

		lt.Locales[name] = tmpl
		return nil
	})
}
