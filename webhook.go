package telebot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// A WebhookTLS specifies the path to a key and a cert so the poller can open
// a TLS listener
type WebhookTLS struct {
	Key  string
	Cert string
}

// A WebhookEndpoint describes the endpoint to which telegram will send its requests.
// This must be a public URL and can be a loadbalancer or something similar. If the
// endpoint uses TLS and the certificate is selfsigned you have to add the certificate
// path of this certificate so telegram will trust it. This field can be ignored if you
// have a trusted certifcate (letsencrypt, ...).
type WebhookEndpoint struct {
	PublicURL string
	Cert      string
}

// A Webhook configures the poller for webhooks. It opens a port on the given
// listen adress. If TLS is filled, the listener will use the key and cert to open
// a secure port. Otherwise it will use plain HTTP.
// If you have a loadbalancer ore other infrastructure in front of your service, you
// must fill the Endpoint structure so this poller will send this data to telegram. If
// you leave these values empty, your local adress will be sent to telegram which is mostly
// not what you want (at least while developing). If you have a single instance of your
// bot you should consider to use the LongPoller instead of a WebHook.
// You can also leave the Listen field empty. In this case it is up to the caller to
// add the Webhook to a http-mux.
type Webhook struct {
	Listen   string
	TLS      *WebhookTLS
	Endpoint *WebhookEndpoint
	dest     chan<- Update
	bot      *Bot
}

type registerResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (h *Webhook) getFiles() map[string]File {
	m := make(map[string]File)

	if h.TLS != nil {
		m["certificate"] = FromDisk(h.TLS.Cert)
	}
	// check if it is overwritten by an endpoint
	if h.Endpoint != nil {
		if h.Endpoint.Cert == "" {
			// this can be the case if there is a loadbalancer or reverseproxy in
			// front with a public cert. in this case we do not need to upload it
			// to telegram. we delete the certificate from the map, because someone
			// can have an internal TLS listener with a private cert
			delete(m, "certificate")
		} else {
			// someone configured a certificate
			m["certificate"] = FromDisk(h.Endpoint.Cert)
		}
	}
	return m
}

func (h *Webhook) getParams() map[string]string {
	param := make(map[string]string)
	if h.TLS != nil {
		param["url"] = fmt.Sprintf("https://%s", h.Listen)
	} else {
		// this will not work with telegram, they want TLS
		// but i allow this because telegram will send an error
		// when you register this hook. in their docs they write
		// that port 80/http is allowed ...
		param["url"] = fmt.Sprintf("http://%s", h.Listen)
	}
	if h.Endpoint != nil {
		param["url"] = h.Endpoint.PublicURL
	}
	return param
}

func (h *Webhook) Poll(b *Bot, dest chan Update, stop chan struct{}) {
	res, err := b.sendFiles("setWebhook", h.getFiles(), h.getParams())
	if err != nil {
		b.debug(fmt.Errorf("setWebhook failed %q: %v", string(res), err))
		close(stop)
		return
	}
	var result registerResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		b.debug(fmt.Errorf("bad json data %q: %v", string(res), err))
		close(stop)
		return
	}
	if !result.Ok {
		b.debug(fmt.Errorf("cannot register webhook: %s", result.Description))
		close(stop)
		return
	}
	// store the variables so the HTTP-handler can use 'em
	h.dest = dest
	h.bot = b

	if h.Listen == "" {
		h.waitForStop(stop)
		return
	}

	s := &http.Server{
		Addr:    h.Listen,
		Handler: h,
	}

	go func(stop chan struct{}) {
		h.waitForStop(stop)
		s.Shutdown(context.Background())
	}(stop)

	if h.TLS != nil {
		s.ListenAndServeTLS(h.TLS.Cert, h.TLS.Key)
	} else {
		s.ListenAndServe()
	}
}

func (h *Webhook) waitForStop(stop chan struct{}) {
	<-stop
	close(stop)
}

// The handler simply reads the update from the body of the requests
// and writes them to the update channel.
func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		h.bot.debug(fmt.Errorf("cannot decode update: %v", err))
		return
	}
	h.dest <- update
}
