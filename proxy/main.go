package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"github.com/yhat/wsutil"
)

var opts struct {
	ListenAddr  string `long:"listen" default:"0.0.0.0:8327" description:"WebSocket listen address and port"`
	UpstreamURL string `long:"upstream" default:"" description:"upstream WebSocket url. e.g.: ws://127.0.0.1:8328/ws"`
	Echo        bool   `long:"echo" description:"Whether to use echo mode"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.InfoLevel
	f := new(logrus.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	log.Formatter = f
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("error: %v", err)
		} else {
			return
		}
	}

	upstream, err := url.Parse(opts.UpstreamURL)
	if err != nil {
		log.Fatalf("Parse upstream url failed, err: %v", err)
	}

	proxy := wsutil.NewSingleHostReverseProxy(upstream)

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	http.ListenAndServe(opts.ListenAddr, proxy)
}
