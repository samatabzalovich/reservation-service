package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *Message
	ID      int64  `json:"id"`
	RoomID  string `json:"roomId"`
}

type Message struct {
	Institution any `json:"institution,omitempty"`
	Status string `json:"status"`
	Content string `json:"content,omitempty"`
	RoomID  string `json:"roomId"`
	UserID  int64  `json:"userId"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				hub.Broadcast <- &Message{
					Status: "disconnected",
					RoomID:  c.RoomID,
					UserID:  c.ID,
				}
			}
			break
		}

		msg := &Message{
			Status: "received",
			Content: string(m),
			RoomID:  c.RoomID,
		}
		hub.Broadcast <- msg
	}
}
