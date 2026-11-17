package telebot

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	gotoken "go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/telebot.v4/internal/botapispec"
)

var updateCoverage = flag.Bool("update", false, "rewrite COVERAGE.md")

// reBotAPI matches the doc comment line that names the Bot API methods a
// function calls, or the Bot API objects a type mirrors:
//
//	// Bot API: sendMessage, sendPhoto
var reBotAPI = regexp.MustCompile(`(?m)^Bot API: (.+?)\.?$`)

// goDecl is an exported declaration of the package.
type goDecl struct {
	name   string   // Bot.Send, Message
	botAPI []string // names from the "Bot API:" doc line
	file   string
}

type goPackage struct {
	funcs     map[string]*goDecl
	types     map[string]*goDecl
	calls     map[string]map[string]bool // func name -> called names
	tags      map[string][]string        // type name -> json tags of its fields
	codecTags map[string][]string        // type name -> json tags in its (Un)MarshalJSON
	embeds    map[string][]string        // type name -> embedded type names
	strings   map[string]bool            // every string literal in the package
}

func parsePackage(t *testing.T) *goPackage {
	fset := gotoken.NewFileSet()
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var files []*ast.File
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, f)
	}

	p := &goPackage{
		funcs:     map[string]*goDecl{},
		types:     map[string]*goDecl{},
		calls:     map[string]map[string]bool{},
		tags:      map[string][]string{},
		codecTags: map[string][]string{},
		embeds:    map[string][]string{},
		strings:   map[string]bool{},
	}

	botAPI := func(doc *ast.CommentGroup) []string {
		if doc == nil {
			return nil
		}
		m := reBotAPI.FindStringSubmatch(doc.Text())
		if m == nil {
			return nil
		}
		var names []string
		for _, n := range strings.Split(m[1], ",") {
			names = append(names, strings.TrimSpace(n))
		}
		return names
	}

	for _, f := range files {
		file := fset.Position(f.Pos()).Filename

		ast.Inspect(f, func(n ast.Node) bool {
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == gotoken.STRING {
				if s, err := strconv.Unquote(lit.Value); err == nil {
					p.strings[s] = true
				}
			}
			return true
		})

		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				name := d.Name.Name
				if d.Recv != nil {
					typ := d.Recv.List[0].Type
					if star, ok := typ.(*ast.StarExpr); ok {
						typ = star.X
					}
					if id, ok := typ.(*ast.Ident); ok {
						name = id.Name + "." + name
					}
				}
				calls := map[string]bool{}
				if d.Body != nil {
					ast.Inspect(d.Body, func(n ast.Node) bool {
						if c, ok := n.(*ast.CallExpr); ok {
							switch fn := c.Fun.(type) {
							case *ast.Ident:
								calls[fn.Name] = true
							case *ast.SelectorExpr:
								calls[fn.Sel.Name] = true
							}
						}
						return true
					})
				}
				p.calls[name] = calls

				// Fields decoded or encoded by hand are declared in the
				// structs of the type's (Un)MarshalJSON.
				if d.Recv != nil && d.Body != nil && (d.Name.Name == "UnmarshalJSON" || d.Name.Name == "MarshalJSON") {
					typ := name[:strings.Index(name, ".")]
					ast.Inspect(d.Body, func(n ast.Node) bool {
						if st, ok := n.(*ast.StructType); ok {
							p.codecTags[typ] = append(p.codecTags[typ], jsonTags(st)...)
						}
						return true
					})
				}
				if ast.IsExported(d.Name.Name) {
					p.funcs[name] = &goDecl{name: name, botAPI: botAPI(d.Doc), file: file}
				}

			case *ast.GenDecl:
				if d.Tok != gotoken.TYPE {
					continue
				}
				for _, s := range d.Specs {
					ts := s.(*ast.TypeSpec)
					doc := ts.Doc
					if doc == nil && len(d.Specs) == 1 {
						doc = d.Doc
					}
					if ast.IsExported(ts.Name.Name) {
						p.types[ts.Name.Name] = &goDecl{name: ts.Name.Name, botAPI: botAPI(doc), file: file}
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					for _, fld := range st.Fields.List {
						if len(fld.Names) != 0 {
							continue
						}
						typ := fld.Type
						if star, ok := typ.(*ast.StarExpr); ok {
							typ = star.X
						}
						if id, ok := typ.(*ast.Ident); ok {
							p.embeds[ts.Name.Name] = append(p.embeds[ts.Name.Name], id.Name)
						}
					}
					p.tags[ts.Name.Name] = append(p.tags[ts.Name.Name], jsonTags(st)...)
				}
			}
		}
	}
	return p
}

// jsonTags returns the json names of the fields of st.
func jsonTags(st *ast.StructType) []string {
	var names []string
	for _, fld := range st.Fields.List {
		if fld.Tag == nil {
			continue
		}
		tag, _ := strconv.Unquote(fld.Tag.Value)
		name := strings.Split(reflect.StructTag(tag).Get("json"), ",")[0]
		if name != "" && name != "-" {
			names = append(names, name)
		}
	}
	return names
}

// reachesAPI reports which functions reach Bot.Raw or Bot.sendFiles, by
// name. Calls through interfaces and same-named methods are conflated, so it
// over-approximates, which is what a "must be documented" check wants.
func (p *goPackage) reachesAPI() map[string]bool {
	byName := map[string][]string{}
	for full := range p.calls {
		short := full[strings.LastIndex(full, ".")+1:]
		byName[short] = append(byName[short], full)
	}
	reach := map[string]bool{"Bot.Raw": true, "Bot.sendFiles": true}
	for changed := true; changed; {
		changed = false
		for fn, calls := range p.calls {
			if reach[fn] {
				continue
			}
			for c := range calls {
				for _, target := range byName[c] {
					if reach[target] && !reach[fn] {
						reach[fn] = true
						changed = true
					}
				}
			}
		}
	}
	return reach
}

// allTags returns the json tags of a type, following embedded structs.
func (p *goPackage) allTags(name string, seen map[string]bool) []string {
	if seen[name] {
		return nil
	}
	seen[name] = true
	tags := append(append([]string(nil), p.tags[name]...), p.codecTags[name]...)
	for _, e := range p.embeds[name] {
		tags = append(tags, p.allTags(e, seen)...)
	}
	return tags
}

func TestCoverage(t *testing.T) {
	spec, err := botapispec.Load()
	if err != nil {
		t.Fatal(err)
	}
	pkg := parsePackage(t)

	methods := map[string]botapispec.Method{}
	for _, m := range spec.Methods {
		methods[m.Name] = m
	}
	types := map[string]botapispec.Type{}
	for _, ty := range spec.Types {
		types[ty.Name] = ty
	}

	// Every name in a "Bot API:" line exists in the spec.
	byMethod := map[string][]string{}
	for _, f := range sortedDecls(pkg.funcs) {
		for _, n := range f.botAPI {
			if _, ok := methods[n]; !ok {
				t.Errorf("%s: %s names %q, which is not a Bot API method", f.file, f.name, n)
			}
			byMethod[n] = append(byMethod[n], f.name)
		}
	}
	byType := map[string][]string{}
	for _, ty := range sortedDecls(pkg.types) {
		for _, n := range ty.botAPI {
			if _, ok := types[n]; !ok {
				t.Errorf("%s: type %s names %q, which is not a Bot API type", ty.file, ty.name, n)
			}
			byType[n] = append(byType[n], ty.name)
		}
	}

	// Every exported Bot method that calls the Bot API says which method(s).
	reach := pkg.reachesAPI()
	for _, f := range sortedDecls(pkg.funcs) {
		if strings.HasPrefix(f.name, "Bot.") && reach[f.name] && f.botAPI == nil && f.name != "Bot.Raw" {
			t.Errorf("%s: %s calls the Bot API but its doc comment has no \"Bot API:\" line", f.file, f.name)
		}
	}

	// A method telebot calls by name is documented somewhere.
	for _, m := range spec.Methods {
		if pkg.strings[m.Name] && byMethod[m.Name] == nil {
			t.Errorf("%s is called but no doc comment names it in a \"Bot API:\" line", m.Name)
		}
	}

	// Every exported struct that carries json fields mirrors a Bot API type,
	// unless it is listed in notBotAPITypes.
	for _, ty := range sortedDecls(pkg.types) {
		if ty.botAPI == nil && len(pkg.tags[ty.name]) > 0 && !notBotAPITypes[ty.name] {
			t.Errorf("%s: type %s has json fields but its doc comment has no \"Bot API:\" line", ty.file, ty.name)
		}
	}

	got := renderCoverage(spec, pkg, byMethod, byType)
	if *updateCoverage {
		if err := os.WriteFile("COVERAGE.md", got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile("COVERAGE.md")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Error("COVERAGE.md is out of date, run: go test -run TestCoverage -update")
	}
}

// notBotAPITypes are exported structs with json tags that are telebot's own,
// not mirrors of a Bot API object.
var notBotAPITypes = map[string]bool{
	"Btn":              true, // keyboard builder
	"CallbackResponse": true, // answerCallbackQuery parameters
	"CommandParams":    true, // setMyCommands parameters
	"Currency":         true, // payments currency table
	"GiftQuery":        true, // gift list parameters
	"QueryResponse":    true, // answerInlineQuery parameters
	"Reactions":        true, // setMessageReaction parameters
	"StoredMessage":    true, // Editable for stored messages
	"WebhookEndpoint":  true, // Webhook settings
	"WebhookTLS":       true, // Webhook settings

	// Sent by VerifyUser and VerifyChat, but not part of the Bot API.
	"VerificationInfo":         true,
	"VerificationRequirements": true,
}

func sortedDecls(m map[string]*goDecl) []*goDecl {
	var out []*goDecl
	for _, d := range m {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func renderCoverage(spec *botapispec.Spec, pkg *goPackage, byMethod, byType map[string][]string) []byte {
	code := func(names []string) string {
		var out []string
		for _, n := range names {
			out = append(out, "`"+n+"`")
		}
		return strings.Join(out, ", ")
	}

	var methodsDone, typesDone int
	var mrows, trows bytes.Buffer
	for _, m := range spec.Methods {
		impl := "**missing**"
		var missing []string
		if names := byMethod[m.Name]; names != nil {
			impl = code(names)
			methodsDone++
			for _, prm := range m.Params {
				if !pkg.strings[prm.Name] && !anyTag(pkg, prm.Name) {
					missing = append(missing, prm.Name)
				}
			}
		}
		fmt.Fprintf(&mrows, "| [%s](https://core.telegram.org/bots/api#%s) | %s | %s |\n",
			m.Name, strings.ToLower(m.Name), impl, strings.Join(missing, ", "))
	}
	for _, ty := range spec.Types {
		impl := "**missing**"
		var missing []string
		if names := byType[ty.Name]; names != nil {
			impl = code(names)
			typesDone++
			have := map[string]bool{}
			for _, n := range names {
				for _, tag := range pkg.allTags(n, map[string]bool{}) {
					have[tag] = true
				}
			}
			for _, f := range ty.Fields {
				// The "type" discriminator of a union member is set by telebot.
				discriminator := f.Name == "type" && (strings.Contains(f.Description, "always") || strings.Contains(f.Description, "must be"))
				if !have[f.Name] && !discriminator {
					missing = append(missing, f.Name)
				}
			}
		}
		fmt.Fprintf(&trows, "| [%s](https://core.telegram.org/bots/api#%s) | %s | %s |\n",
			ty.Name, strings.ToLower(ty.Name), impl, strings.Join(missing, ", "))
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# Bot API coverage\n\n")
	fmt.Fprintf(&buf, "Generated by `go test -run TestCoverage -update` from the `Bot API:` lines in doc comments and the spec in [internal/botapispec/api.json](internal/botapispec/api.json) (%s, %s). Do not edit by hand.\n\n", spec.Version, spec.ReleaseDate)
	fmt.Fprintf(&buf, "Methods: %d of %d. Types: %d of %d.\n\n", methodsDone, len(spec.Methods), typesDone, len(spec.Types))
	fmt.Fprintf(&buf, "A method's \"unseen parameters\" are parameters whose name appears nowhere in telebot's code; a type's \"missing fields\" are fields none of its telebot types decode.\n\n")
	fmt.Fprintf(&buf, "## Methods\n\n| Bot API | telebot | unseen parameters |\n|---|---|---|\n")
	buf.Write(mrows.Bytes())
	fmt.Fprintf(&buf, "\n## Types\n\n| Bot API | telebot | missing fields |\n|---|---|---|\n")
	buf.Write(trows.Bytes())
	return buf.Bytes()
}

func anyTag(pkg *goPackage, name string) bool {
	for _, tags := range pkg.tags {
		for _, t := range tags {
			if t == name {
				return true
			}
		}
	}
	return false
}
