//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/url"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
)

var serverContent string

func main() {
	addr := "localhost:8080"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go read(c, done)
	clip := make(chan string)
	go detectClipChange(clip)
	write(c, done, clip)
}

func read(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		serverContent = string(message)
		err = clipboard.WriteAll(string(message))
		if err != nil {
			log.Println(err)
		}
		log.Println("read:", string(message))
	}
}

func detectClipChange(c chan string) {
	prevContent := ""
	for {
		content, err := clipboard.ReadAll()
		if err != nil {
			log.Println(err)
		}
		if content != prevContent && content != serverContent {
			c <- content
			prevContent = content
		}
		time.Sleep(time.Second)
	}
}

func write(c *websocket.Conn, done chan struct{}, clip chan string) {
	for {
		select {
		case <-done:
			return
		case content := <-clip:
			err := c.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				log.Println(err)
			}
			log.Println("send:", content)
		}
	}
}
