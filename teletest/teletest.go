// Package teletest is a fake Telegram Bot API server for offline tests.
//
// It records every request a bot makes and answers each method with a reply
// you choose, so a test can drive a bot without a real token or network and
// then check exactly what the bot sent:
//
//	srv := teletest.NewServer(t)
//	srv.Reply("sendMessage", map[string]interface{}{
//		"message_id": 1,
//		"chat":       map[string]interface{}{"id": 42, "type": "private"},
//	})
//
//	b, _ := tele.NewBot(tele.Settings{
//		URL:     srv.URL,
//		Token:   srv.Token,
//		Offline: true,
//	})
//	b.Send(tele.ChatID(42), "hello")
//
//	req := srv.Last()
//	// req.Method == "sendMessage", req.Params["text"] == "hello"
//
// The package does not import telebot, so it works with any client that
// speaks the Bot API over HTTP.
package teletest

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// DefaultToken is the bot token NewServer accepts.
const DefaultToken = "123456:TEST-token"

// TB is the part of testing.TB the server uses.
type TB interface {
	Helper()
	Cleanup(func())
	Errorf(format string, args ...interface{})
}

// Server is a fake Bot API server. Its zero value is not usable; create one
// with NewServer.
type Server struct {
	// URL is the base URL to pass as telebot's Settings.URL.
	URL string

	// Token is the bot token the server accepts. Requests with another
	// token are answered with 401 Unauthorized, as Telegram does.
	Token string

	srv *httptest.Server

	mu       sync.Mutex
	requests []Request
	replies  map[string][]reply
}

// Request is one recorded Bot API call.
type Request struct {
	// Method is the Bot API method, such as "sendMessage".
	Method string

	// Params holds the call's parameters. Strings are kept as they are;
	// numbers, booleans and objects are kept as their JSON encoding, which
	// is also how telebot sends them in multipart requests.
	Params map[string]string

	// Files holds files uploaded with a multipart request, by field name.
	Files map[string]File
}

// File is an uploaded file.
type File struct {
	Name string
	Data []byte
}

// JSON decodes the JSON-encoded parameter name into v, failing the test if
// it cannot.
func (r Request) JSON(t TB, name string, v interface{}) {
	t.Helper()
	if err := json.Unmarshal([]byte(r.Params[name]), v); err != nil {
		t.Errorf("teletest: %s: parameter %s: %v", r.Method, name, err)
	}
}

type reply struct {
	status int
	body   []byte
}

// NewServer starts a fake Bot API server that is closed when the test ends.
func NewServer(t TB) *Server {
	t.Helper()
	s := &Server{Token: DefaultToken, replies: map[string][]reply{}}
	s.srv = httptest.NewServer(http.HandlerFunc(s.serve))
	s.URL = s.srv.URL
	t.Cleanup(s.srv.Close)
	return s
}

// Client returns an HTTP client for the server, to pass as Settings.Client.
func (s *Server) Client() *http.Client {
	return s.srv.Client()
}

// Reply sets the result Telegram returns for method. The result is encoded
// to JSON, or used as it is when it is a json.RawMessage or []byte. Replies
// queue up: each call to the method takes the next reply, and the last one
// keeps being returned. Methods without a reply return true.
func (s *Server) Reply(method string, result interface{}) {
	var raw []byte
	switch r := result.(type) {
	case json.RawMessage:
		raw = r
	case []byte:
		raw = r
	default:
		var err error
		if raw, err = json.Marshal(result); err != nil {
			panic(fmt.Sprintf("teletest: reply for %s: %v", method, err))
		}
	}
	body, _ := json.Marshal(map[string]interface{}{"ok": true, "result": json.RawMessage(raw)})
	s.push(method, reply{status: http.StatusOK, body: body})
}

// ReplyError makes method fail with the Bot API error code and description,
// such as 400 and "Bad Request: chat not found".
func (s *Server) ReplyError(method string, code int, description string) {
	body, _ := json.Marshal(map[string]interface{}{
		"ok":          false,
		"error_code":  code,
		"description": description,
	})
	s.push(method, reply{status: code, body: body})
}

func (s *Server) push(method string, r reply) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replies[method] = append(s.replies[method], r)
}

// Requests returns every request recorded so far, oldest first.
func (s *Server) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Request(nil), s.requests...)
}

// Calls returns the recorded requests for method, oldest first.
func (s *Server) Calls(method string) []Request {
	var out []Request
	for _, r := range s.Requests() {
		if r.Method == method {
			out = append(out, r)
		}
	}
	return out
}

// Last returns the most recent request, or a zero Request if there is none.
func (s *Server) Last() Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) == 0 {
		return Request{}
	}
	return s.requests[len(s.requests)-1]
}

// Reset forgets the recorded requests and the queued replies.
func (s *Server) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = nil
	s.replies = map[string][]reply{}
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// The path is /bot<token>/<method>.
	path := strings.TrimPrefix(r.URL.Path, "/bot")
	slash := strings.LastIndexByte(path, '/')
	if slash < 0 || path == r.URL.Path {
		http.NotFound(w, r)
		return
	}
	token, method := path[:slash], path[slash+1:]
	if token != s.Token {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"ok":false,"error_code":401,"description":"Unauthorized"}`)
		return
	}

	req, err := decode(method, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]interface{}{
			"ok": false, "error_code": 400, "description": "teletest: " + err.Error(),
		})
		w.Write(body)
		return
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	rep := reply{status: http.StatusOK, body: []byte(`{"ok":true,"result":true}`)}
	if q := s.replies[method]; len(q) > 0 {
		rep = q[0]
		if len(q) > 1 {
			s.replies[method] = q[1:]
		}
	}
	s.mu.Unlock()

	w.WriteHeader(rep.status)
	w.Write(rep.body)
}

func decode(method string, r *http.Request) (Request, error) {
	req := Request{Method: method, Params: map[string]string{}}

	ctype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch ctype {
	case "multipart/form-data":
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return req, err
		}
		for k, v := range r.MultipartForm.Value {
			req.Params[k] = v[0]
		}
		for k, fhs := range r.MultipartForm.File {
			data, err := readFile(fhs[0])
			if err != nil {
				return req, err
			}
			if req.Files == nil {
				req.Files = map[string]File{}
			}
			req.Files[k] = File{Name: fhs[0].Filename, Data: data}
		}

	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return req, err
		}
		for k, v := range r.PostForm {
			req.Params[k] = v[0]
		}

	default: // JSON body, or none
		for k, v := range r.URL.Query() {
			req.Params[k] = v[0]
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return req, err
		}
		if len(strings.TrimSpace(string(body))) == 0 || strings.TrimSpace(string(body)) == "null" {
			return req, nil
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(body, &fields); err != nil {
			return req, err
		}
		for k, v := range fields {
			var s string
			if json.Unmarshal(v, &s) == nil {
				req.Params[k] = s
			} else {
				req.Params[k] = string(v)
			}
		}
	}
	return req, nil
}

func readFile(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
