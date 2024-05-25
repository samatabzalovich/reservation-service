package main

import (
	"errors"
	"net/http"
	"queue-managemant-service/internal/data"
	"strconv"
	"time"
)

func (app *Config) CallNextClient(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ServiceID        int64  `json:"serviceId"`
		MessageForClient string `json:"messageForClient"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	service, err := app.GetServiceById(input.ServiceID)
	if err != nil {
		app.errorJson(w, err, http.StatusForbidden)
		return
	}
	user, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	
	queue, err := app.Models.Queue.GetLastPositionedQueue(service.ID)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if queue.QueueStatus == "completed" {
		app.errorJson(w, data.ErrNoClientInQueue, http.StatusNotFound)
	}
	employeeId, err := app.GetEmployeeForServiceAndUserID(service.ID, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorJson(w, data.ErrInvalidToken, http.StatusForbidden)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	queue.EmployeeID = &employeeId
	err = app.Models.Queue.CallFromQueue(queue)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	

	users, err := app.Models.Queue.GetUsersFromQueue(service.ID)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}


	app.hub.Broadcast <- &Message{
		RoomID:  app.GetServiceRoom(input.ServiceID),
		Service: service,
		Status:  "called",
		Content: map[string]any{
			"peopleLeft": len(users),
		},
		UserID: queue.ClientID,
		Users: users,
	}

	err = app.Redis.Set(r.Context(), app.GetUserRoom(queue.ClientID), input.MessageForClient, 10*60*time.Second).Err()
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.hub.Broadcast <- &Message{
		RoomID:  app.GetUserRoom(queue.ClientID),
		Service: service,
		Status:  "messageForClient",
		Content: map[string]any{
			"peopleLeft": len(users),
			"queue":      queue,
		},
		UserID:           queue.ClientID,
		MessageForClient: input.MessageForClient,
		Users:            users,
	}

	queue.QueueStatus = "called"
	app.writeJSON(w, http.StatusOK, map[string]any{"message": "Client called", "clientId": queue.ClientID, "queuePosition": queue.Position, "peopleInQueue": len(users), "queueId": queue.ID})
}

func (app *Config) GetAllForInstitution(w http.ResponseWriter, r *http.Request) {
	instId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(app.readString(r.URL.Query(), "pageSize", "10"))
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}
	page, err := strconv.Atoi(app.readString(r.URL.Query(), "page", "1"))
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}
	queues, metadata, err := app.Models.Queue.GetAllForInst(instId, data.Filters{Page: page, PageSize: pageSize})
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{
		"queues":   queues,
		"metadata": metadata,
	})
}

func (app *Config) DeleteAllForInst(w http.ResponseWriter, r *http.Request) {
	instId, err := app.readIntParam(r, "instId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	err = app.Models.Queue.DeleteAllForInst(instId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]string{"message": "All queues deleted"})
}

func (app *Config) DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queueId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	err = app.Models.Queue.Delete(queueId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Queue deleted"})
}

func (app *Config) GetQueue(w http.ResponseWriter, r *http.Request) {
	queueId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	queue, err := app.Models.Queue.GetById(queueId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, queue)
}

func (app *Config) UpdateQueueStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          int64  `json:"queueId"`
		QueueStatus string `json:"queueStatus"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	queue, err := app.Models.Queue.GetById(input.ID)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorJson(w, err, http.StatusNotFound)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	if queue.QueueStatus == "completed" || input.QueueStatus == "cancelled" {
		app.errorJson(w, data.ErrInvalidQueueInfo, http.StatusForbidden)
		return
	}

	queue.QueueStatus = input.QueueStatus
	err = app.Models.Queue.Update(queue)

	if err != nil {
		if errors.Is(err, data.ErrConcurrentUpdate) {
			app.errorJson(w, err, http.StatusConflict)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	users, err := app.Models.Queue.GetUsersFromQueue(queue.ServiceID)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	queueRes := &QueueRes{
		Queue:      queue,
		PeopleLeft: len(users),
	}
	app.hub.Broadcast <- &Message{
		RoomID:  app.GetServiceRoom(queue.ServiceID),
		Status:  queue.QueueStatus,
		Content: queueRes,
		UserID:  queue.ClientID,
		Users:   users,
	}

	if queue.QueueStatus == "completed" {
		// app.hub.Unregister <- &Client{ //TODO: check if this is needed
		// 	ID:     queue.ClientID,
		// 	RoomID: app.GetServiceRoom(queue.ServiceID),
		// }
		// app.hub.Unregister <- &Client{
		// 	ID:     queue.ClientID,
		// 	RoomID: app.GetUserRoom(queue.ClientID),
		// }
		_, err := app.Redis.Del(r.Context(), app.GetUserRoom(queue.ClientID)).Result()
		if err != nil {
			app.errorJson(w, err)
			return
		}
	}

	app.writeJSON(w, http.StatusOK, queue)
}

func (app *Config) GetQueueNumberForClientInInstitution(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "clientId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	instId, err := app.readIntFromUrl(r.URL.Query(), "instId", 0)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	employeeId, err := app.readIntFromUrl(r.URL.Query(), "employeeId", 0)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	count, err := app.Models.Queue.GetQueueForClientInInstitution(clientId, instId, employeeId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (app *Config) GetQueuesForClient(w http.ResponseWriter, r *http.Request) {
	client, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	queues, _, err := app.Models.Queue.GetForClient(client.ID)

	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"queues": queues})
}
