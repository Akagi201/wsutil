package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"golang.org/x/net/websocket"
)

var opts struct {
	OriginURL  string   `long:"origin" default:"http://localhost/" description:"Origin url of the WebSocket Client"`
	WSURL      string   `long:"ws" default:"ws://localhost:8327/ws" description:"WebSocket Server URL to connect to"`
	Buffer     int      `long:"buffer" default:"1024" description:"WebSocket receive buffer size"`
	Protocols  []string `long:"protocols" description:"WebSocket subprotocols."`
	SkipVerify bool     `long:"skipverify" description:"Skip TLS certificate verification"`
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

func dial(wsurl string, origin string, protocols []string) (*websocket.Conn, error) {
	config, err := websocket.NewConfig(wsurl, origin)
	if err != nil {
		return nil, err
	}

	if len(protocols) != 0 {
		config.Protocol = protocols
	}

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: opts.SkipVerify,
	}

	return websocket.DialConfig(config)
}

func inLoop(ws *websocket.Conn, errors chan<- error, in chan<- []byte) {
	msg := make([]byte, opts.Buffer)

	for {
		n, err := ws.Read(msg)
		if err != nil {
			errors <- err
			continue
		}

		in <- msg[:n]
	}
}

func outLoop(ws *websocket.Conn, out <-chan []byte, errors chan<- error) {
	for msg := range out {
		_, err := ws.Write(msg)
		if err != nil {
			errors <- err
		}
	}
}

func printErrors(errors <-chan error) {
	for err := range errors {
		if err == io.EOF {
			log.Fatalf("\rConnection closed by remote, err: %v", err)
		} else {
			fmt.Printf("\rGot err: %v\n> ", err)
		}
	}
}

func printMsgs(in <-chan []byte) {
	for msg := range in {
		fmt.Printf("\r< %v\n> ", string(msg))
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("Cli flags parse error: %v", err)
		} else {
			return
		}
	}

	ws, err := dial(opts.WSURL, opts.OriginURL, opts.Protocols)
	if err != nil {
		log.Fatalf("WebSocket dial error: %v", err)
	}
	defer ws.Close()

	if len(opts.Protocols) != 0 {
		log.Printf("Connected to %v via %v from %v...", opts.WSURL, opts.Protocols, opts.OriginURL)
	} else {
		log.Printf("Connected to %v from %v...", opts.WSURL, opts.OriginURL)
	}

	var wg sync.WaitGroup

	wg.Add(3)

	errors := make(chan error)
	in := make(chan []byte)
	out := make(chan []byte)

	defer close(errors)
	defer close(in)
	defer close(out)

	go inLoop(ws, errors, in)
	go printMsgs(in)
	go printErrors(errors)
	go outLoop(ws, out, errors)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() {
		out <- []byte(scanner.Text())
		fmt.Print("> ")
	}

	wg.Wait()
}
