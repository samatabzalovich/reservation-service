package main

import (
	"errors"
	"net/http"
	"queue-managemant-service/internal/data"
	"strconv"
)

type JoinRoomReq struct {
	RoomID string `json:"roomId"`
	UserID int64  `json:"userId"`
}

func (app *Config) JoinQueueForServiceRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomReq
	id, err := app.readIntParam(r, "serviceId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	req.RoomID = strconv.FormatInt(id, 10)
	service, err := app.GetServiceById(id)
	if err != nil {
		app.errorJson(w, err, http.StatusForbidden)
		return
	}

	// check if room exists
	_, ok := app.hub.Rooms[req.RoomID]
	if !ok {
		app.hub.Rooms[req.RoomID] = &Room{
			ID:      req.RoomID,
			Clients: make(map[int64]*Client),
		}
	}
	user, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	var queue *data.Queue
	length, err := app.GetQueueLength(id, r.Context())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if user.Type == "client" {
		queue, err = app.Models.Queue.GetLastForClient(id)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				queue, err = app.Models.Queue.Insert(user.ID, service.InstId, service.ID)
				if err != nil {
					app.errorJson(w, err)
					return
				}
				err = app.IncreaseQueueLength(id, r.Context())
				if err != nil {
					app.errorJson(w, err)
					return
				}
				length++
			} else {
				app.errorJson(w, err)
				return
			}
		}

	}
	serviceRoom := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      user.ID,
		RoomID:  req.RoomID,
	}

	m := &Message{
		Service: service,
		Status:  "connected",
		RoomID:  req.RoomID,
		UserID:  user.ID,
		Content: map[string]any{
			"queue":      queue,
			"peopleLeft": length,
		},
	}

	app.hub.Register <- serviceRoom
	app.hub.Broadcast <- m
	if user.Type == "client" {
		cl := &Client{
			Conn:    conn,
			Message: make(chan *Message, 10),
			ID:      user.ID,
			RoomID:  strconv.FormatInt(user.ID, 10),
		}
		app.hub.Register <- cl
		go func() {
			go cl.writeMessage()
			cl.readMessage(app.hub)
		}()
	}
	go serviceRoom.writeMessage()
	serviceRoom.readMessage(app.hub)
}

type RoomRes struct {
	ID string `json:"id"`
}

type ClientRes struct {
	ID string `json:"id"`
}
