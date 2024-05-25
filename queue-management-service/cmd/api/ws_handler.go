package main

import (
	"errors"
	"fmt"
	"net/http"
	"queue-managemant-service/internal/data"

	"github.com/redis/go-redis/v9"
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
	req.RoomID = app.GetServiceRoom(id)
	service, err := app.GetServiceById(id)
	if err != nil {
		app.errorJson(w, err, http.StatusForbidden)
		return
	}

	
	user, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}
	var queue *data.Queue
	if user.Type == "client" {
		queue, err = app.Models.Queue.GetLastForClient(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				queue, err = app.Models.Queue.Insert(user.ID, service.InstId, service.ID)
				if err != nil {
					app.errorJson(w, err)
					return
				}
			} else {
				app.errorJson(w, err)
				return
			}
		}

	}
	users,err := app.Models.Queue.GetUsersFromQueue(service.ID)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	queueRes := &QueueRes{
		Queue:      queue,
		PeopleLeft: len(users),
	}
	
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.errorJson(w, err)
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
	sentMessage := false
	if user.Type == "client" {
		cl := &Client{
			Conn:    conn,
			Message: make(chan *Message, 10),
			ID:      user.ID,
			RoomID:  app.GetUserRoom(user.ID),
		}
		_, ok := app.hub.Rooms[cl.RoomID]
		if !ok {
			app.hub.Rooms[cl.RoomID] = &Room{
				ID:      cl.RoomID,
				Clients: make(map[int64]*Client),
			}
		}
		app.hub.Register <- cl
		if queue != nil {
			var messageForClient string
			command := app.Redis.Get(r.Context(), app.GetUserRoom(user.ID))
			messageForClient, err = command.Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					messageForClient = ""
				} else {
					app.hub.Unregister <- cl
				}
			}

			if messageForClient != "" {
				app.hub.Broadcast <- &Message{
					Service:          service,
					Status:           "connected",
					RoomID:           app.GetUserRoom(user.ID),
					UserID:           user.ID,
					Content:          queueRes,
					MessageForClient: messageForClient,
					Users: 		  users,
				}
				sentMessage = true
			}
		}

		go func() {
			go cl.writeMessage()
			cl.readMessage(app.hub)
		}()
	}

	m := &Message{
		Service: service,
		Status:  "connected",
		RoomID:  req.RoomID,
		UserID:  user.ID,
		Content: queueRes,
		Users: users,
	}

	serviceRoom := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      user.ID,
		RoomID:  req.RoomID,
	}

	app.hub.Register <- serviceRoom
	if !sentMessage {
		app.hub.Broadcast <- m
	}
	go serviceRoom.writeMessage()
	serviceRoom.readMessage(app.hub)
}

func (app *Config) JoinQueueForPeopleAmountRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomReq
	serviceId, err := app.readIntParam(r, "serviceId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	
	req.RoomID = app.GetServiceRoom(serviceId)
	service, err := app.GetServiceById(serviceId)
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
	users, err := app.Models.Queue.GetUsersFromQueue(serviceId)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if user.Type == "client" {
		queue, err = app.Models.Queue.GetLastForClient(user.ID)
		if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
			app.errorJson(w, err)
			return
		}
	}
	queueRes := &QueueRes{
		Queue:      queue,
		PeopleLeft: len(users),
	}

	sentMessage := false
	if user.Type == "client" {
		cl := &Client{
			Conn:    conn,
			Message: make(chan *Message, 10),
			ID:      user.ID,
			RoomID:  app.GetUserRoom(user.ID),
		}
		_, ok := app.hub.Rooms[cl.RoomID]
		if !ok {
			app.hub.Rooms[cl.RoomID] = &Room{
				ID:      cl.RoomID,
				Clients: make(map[int64]*Client),
			}
		}
		app.hub.Register <- cl
		if queue != nil {
			var messageForClient string
			command := app.Redis.Get(r.Context(), app.GetUserRoom(user.ID))
			messageForClient, err = command.Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					messageForClient = ""
				} else {
					app.hub.Unregister <- cl
				}
			}

			if messageForClient != "" {
				app.hub.Broadcast <- &Message{
					Service:          service,
					Status:           "connected",
					RoomID:           app.GetUserRoom(user.ID),
					UserID:           user.ID,
					Content:          queueRes,
					MessageForClient: messageForClient,
					Users: 		  users,
				}
				sentMessage = true
			}
		}

		go func() {
			go cl.writeMessage()
			cl.readMessage(app.hub)
		}()
	}

	m := &Message{
		Service: service,
		Status:  "connected",
		RoomID:  req.RoomID,
		UserID:  user.ID,
		Content: queueRes,
		Users: users,
	}

	serviceRoom := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      user.ID,
		RoomID:  req.RoomID,
	}

	app.hub.Register <- serviceRoom
	if !sentMessage {
		app.hub.Broadcast <- m
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
type QueueRes struct {
	Queue      *data.Queue `json:"queue,omitempty"`
	PeopleLeft int         `json:"peopleLeft"`
}

func (app *Config) GetUserRoom(userId int64) string {
	return fmt.Sprintf("%dUSER_ROOM", userId)
}

func (app *Config) GetServiceRoom(serviceId int64) string {
	return fmt.Sprintf("%dSERVICE_ROOM", serviceId)
}
