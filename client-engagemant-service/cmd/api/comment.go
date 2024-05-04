package main

import (
	"net/http"
	"client-engagemant-service/internal/data"
)

// Insert(c *Comment) error
// 		GetById(id int64) (*Comment, error)
// 		GetByInstitutionId(instId int64) ([]*Comment, error)
// 		GetByUserId(userId int64) ([]*Comment, error)
// 		Update(c *Comment) error
// 		Delete(id int64) error
// 		DeleteByInstitutionId(instId int64) error
// 		DeleteByUserId(userId int64) error

func(app *Config) LeaveComment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Comment string `json:"comment"`
		InstID  int64  `json:"inst_id"`
		UserID  int64  `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	comment := &data.Comment{
		Comment: input.Comment,
		InstitutionId: input.InstID,
		UserId: input.UserID,
	}

	err = app.Models.Comment.Insert(comment)
	if err != nil {
		if err == data.ErrInvalidField {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		app.errorJson(w, err)
	}
}

func(app *Config) GetCommentsForInstitution(w http.ResponseWriter, r *http.Request) {
	instID, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	comments, err := app.Models.Comment.GetByInstitutionId(instID)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, comments)
}

func(app *Config) GetCommentsForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	comments, err := app.Models.Comment.GetByUserId(userID)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, comments)
}

func(app *Config) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      int64  `json:"id"`
		Comment string `json:"comment"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	comment := &data.Comment{
		ID:      input.ID,
		Comment: input.Comment,
	}

	err = app.Models.Comment.Update(comment)
	if err != nil {
		app.errorJson(w, err)
	}
}

func(app *Config) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.Models.Comment.Delete(id)
	if err != nil {
		app.errorJson(w, err)
	}
}

func(app *Config) DeleteCommentsForInstitution(w http.ResponseWriter, r *http.Request) {
	instID, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.Models.Comment.DeleteByInstitutionId(instID)
	if err != nil {
		app.errorJson(w, err)
	}
}

func(app *Config) DeleteCommentsForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.Models.Comment.DeleteByUserId(userID)
	if err != nil {
		app.errorJson(w, err)
	}
}


func(app *Config) GetComment(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	comment, err := app.Models.Comment.GetById(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, comment)
}
