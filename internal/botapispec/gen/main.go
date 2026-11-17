// Command gen turns the Bot API documentation page into api.json, the
// machine-readable spec vendored next to it.
//
//	curl -sL https://core.telegram.org/bots/api -o /tmp/api.html
//	go run ./internal/botapispec/gen /tmp/api.html > internal/botapispec/api.json
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"gopkg.in/telebot.v4/internal/botapispec"
)

var (
	reVersion = regexp.MustCompile(`<strong>(Bot API [0-9.]+)</strong>`)
	reDate    = regexp.MustCompile(`(?s)<h4><a class="anchor" name="[^"]*" href="[^"]*"><i class="anchor-icon"></i></a>([^<]+)</h4>\s*<p><strong>Bot API`)
	reH3      = regexp.MustCompile(`<h3><a class="anchor" name="[^"]*" href="[^"]*"><i class="anchor-icon"></i></a>([^<]+)</h3>`)
	reH4      = regexp.MustCompile(`<h4><a class="anchor" name="[^"]*" href="[^"]*"><i class="anchor-icon"></i></a>([^<]+)</h4>`)
	reTable   = regexp.MustCompile(`(?s)<table class="table">(.*?)</table>`)
	reRow     = regexp.MustCompile(`(?s)<tr>(.*?)</tr>`)
	reCell    = regexp.MustCompile(`(?s)<td>(.*?)</td>`)
	reHead    = regexp.MustCompile(`(?s)<th>(.*?)</th>`)
	reList    = regexp.MustCompile(`(?s)^\s*<p>[^\n]*</p>\s*<ul>(.*?)</ul>`)
	reItem    = regexp.MustCompile(`(?s)<li><a href="#[^"]*">([A-Za-z]+)</a></li>`)
	reDesc    = regexp.MustCompile(`(?s)<p>(.*?)</p>`)
	reTag     = regexp.MustCompile(`<[^>]+>`)
	reSpace   = regexp.MustCompile(`\s+`)
)

func text(s string) string {
	s = strings.ReplaceAll(s, "<br>", " ")
	s = reTag.ReplaceAllString(s, "")
	return strings.TrimSpace(reSpace.ReplaceAllString(html.UnescapeString(s), " "))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gen api.html")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	doc := string(raw)

	spec := botapispec.Spec{}
	if m := reVersion.FindStringSubmatch(doc); m != nil {
		spec.Version = m[1]
	}
	if m := reDate.FindStringSubmatch(doc); m != nil {
		spec.ReleaseDate = m[1]
	}

	start := strings.Index(doc, `name="getting-updates"`)
	if start < 0 || spec.Version == "" {
		fmt.Fprintln(os.Stderr, "gen: unexpected page layout")
		os.Exit(1)
	}

	// Walk h3 sections and the h4 entries inside them.
	h3s := reH3.FindAllStringSubmatchIndex(doc, -1)
	section := func(at int) string {
		name := ""
		for _, h := range h3s {
			if h[0] > at {
				break
			}
			name = doc[h[2]:h[3]]
		}
		return html.UnescapeString(name)
	}

	h4s := reH4.FindAllStringSubmatchIndex(doc, -1)
	for i, h := range h4s {
		if h[0] < start {
			continue
		}
		name := doc[h[2]:h[3]]
		if strings.ContainsAny(name, " ") {
			continue
		}
		end := len(doc)
		if i+1 < len(h4s) {
			end = h4s[i+1][0]
		}
		for _, h3 := range h3s {
			if h3[0] > h[1] && h3[0] < end {
				end = h3[0]
			}
		}
		body := doc[h[1]:end]

		var desc []string
		for _, p := range reDesc.FindAllStringSubmatch(body, -1) {
			desc = append(desc, text(p[1]))
		}

		var fields []botapispec.Field
		if t := reTable.FindStringSubmatch(body); t != nil {
			heads := reHead.FindAllStringSubmatch(t[1], -1)
			for _, row := range reRow.FindAllStringSubmatch(t[1], -1) {
				cells := reCell.FindAllStringSubmatch(row[1], -1)
				if len(cells) == 0 {
					continue
				}
				f := botapispec.Field{Name: text(cells[0][1]), Type: text(cells[1][1])}
				if len(heads) == 4 { // Parameter, Type, Required, Description
					f.Required = text(cells[2][1]) == "Yes"
					f.Description = text(cells[3][1])
				} else { // Field, Type, Description
					f.Description = text(cells[2][1])
					f.Required = !strings.HasPrefix(f.Description, "Optional")
				}
				fields = append(fields, f)
			}
		}

		if first := name[0]; first >= 'a' && first <= 'z' {
			spec.Methods = append(spec.Methods, botapispec.Method{
				Name:        name,
				Section:     section(h[0]),
				Description: strings.Join(desc, "\n"),
				Params:      fields,
			})
			continue
		}

		typ := botapispec.Type{
			Name:        name,
			Section:     section(h[0]),
			Description: strings.Join(desc, "\n"),
			Fields:      fields,
		}
		if fields == nil {
			if l := reList.FindStringSubmatch(body); l != nil {
				for _, it := range reItem.FindAllStringSubmatch(l[1], -1) {
					typ.Subtypes = append(typ.Subtypes, it[1])
				}
			}
		}
		spec.Types = append(spec.Types, typ)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(spec); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
