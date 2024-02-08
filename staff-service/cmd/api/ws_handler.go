package main

import (
	"net/http"
)

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JoinRoomReq struct {
	RoomID string `json:"roomId"`
	UserID int64  `json:"userId"`
}

func (app *Config) JoinRegisterEmployeeRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomReq
	id, err := app.readStringParam(r, "token")
	if err != nil {
		app.errorJson(w, err, http.StatusForbidden)
		return
	}
	req.RoomID = id
	inst, err := app.GetInstitutionForToken(req.RoomID)
	if err != nil {
		app.errorJson(w, err, http.StatusForbidden)
		return
	}

	// check if room exists
	_,ok := app.hub.Rooms[req.RoomID]
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
	cl := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      user.ID,
		RoomID:  req.RoomID,
	}

	m := &Message{
		Institution: inst,
		Status:     "connected",
		RoomID:      req.RoomID,
		UserID:      user.ID,
	}

	app.hub.Register <- cl
	app.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(app.hub)
}

type RoomRes struct {
	ID string `json:"id"`
}

type ClientRes struct {
	ID string `json:"id"`
}
