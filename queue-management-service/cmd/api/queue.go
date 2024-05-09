package main

import (
	"errors"
	"net/http"
	"queue-managemant-service/internal/data"
	"strconv"
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
	queue, err := app.Models.Queue.GetLastPositionedQueue(service.ID)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if queue.QueueStatus == "completed" {
		app.errorJson(w, data.ErrNoClientInQueue, http.StatusNotFound)
	}
	err = app.Models.Queue.CallFromQueue(queue)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	val, err := app.GetQueueLength(service.ID, r.Context())

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.hub.Broadcast <- &Message{
		RoomID:  strconv.FormatInt(input.ServiceID, 10),
		Service: service,
		Status:  "called",
		Content: map[string]any{
			"peopleLeft": val,
		},
		UserID: queue.ClientID,
	}
	app.hub.Broadcast <- &Message{
		RoomID:  strconv.FormatInt(queue.ClientID, 10),
		Service: service,
		Status:  "called",
		Content: map[string]any{
			"peopleLeft": val,
		},
		UserID:           queue.ClientID,
		MessageForClient: input.MessageForClient,
	}
	queue.QueueStatus = "called"
	app.writeJSON(w, http.StatusOK, map[string]any{"message": "Client called", "clientId": queue.ClientID, "queuePosition": queue.Position, "peopleInQueue": val, "queueId": queue.ID})
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

	if queue.QueueStatus == "completed" {
		app.hub.Unregister <- &Client{
			ID:     queue.ClientID,
			RoomID: strconv.FormatInt(queue.ServiceID, 10),
		}
		app.hub.Unregister <- &Client{
			ID:     queue.ClientID,
			RoomID: strconv.FormatInt(queue.ClientID, 10),
		}
	}

	if queue.QueueStatus == "completed" || queue.QueueStatus == "cancelled" {
		err = app.DecreaseQueueLength(queue.ServiceID, r.Context())
		if err != nil {
			app.errorJson(w, err)
			return
		}
	}
	val, err := app.GetQueueLength(queue.ServiceID, r.Context())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.hub.Broadcast <- &Message{
		RoomID: strconv.FormatInt(queue.ClientID, 10),
		Status: queue.QueueStatus,
		Content: map[string]any{
			"queue":      queue,
			"peopleLeft": val,
		},
		UserID: queue.ClientID,
	}

	app.writeJSON(w, http.StatusOK, queue)
}
