package main

import (
	"io"
	"strings"

	"github.com/Akagi201/light"
	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"golang.org/x/net/websocket"
)

var opts struct {
	ListenAddr string `long:"listen" default:"0.0.0.0:8327" description:"WebSocket listen address and port"`
	Echo       bool   `long:"echo" description:"Whether to use echo mode"`
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

func handleWS(ws *websocket.Conn) {
	if opts.Echo {
		io.Copy(ws, ws)
		return
	}

	var b []byte
	err := websocket.Message.Receive(ws, &b)
	if err != nil {
		log.Printf("WebSocket Message receive failed, error: $v", err)
		return
	}
	log.Printf("WebSocket Message received: %v", string(b))
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

	app := light.New()

	app.Get("/", websocket.Handler(handleWS))

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	app.Listen(opts.ListenAddr)
}
