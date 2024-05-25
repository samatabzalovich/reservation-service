package main

import (
	"client-engagemant-service/internal/data"
	"errors"
	"net/http"
)

func (app *Config) LeaveFeedbackForAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AppointmentId int64  `json:"appointmentId"`
		Rating        int    `json:"rating"`
		Comment       string `json:"comment"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	user, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	ids, err := app.Models.Rating.GetServiceIdAndInstIdToRatingAppointment(input.AppointmentId, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrInvalidAppointmentId) {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}

		app.errorJson(w, err)
		return
	}

	feedback, err := data.NewRating(input.AppointmentId, ids.EmployeeId, user.ID, ids.InstitutionId, input.Rating, input.Comment, ids.ServiceId)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Rating.Insert(feedback)

	if err != nil {
		if errors.Is(err, data.ErrInvalidField) {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		if errors.Is(err, data.ErrRatingAlreadyExists) {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, feedback)
}

func (app *Config) GetFeedbackForAppointment(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	feedback, err := app.Models.Rating.GetRatingForAppointment(id)

	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorJson(w, err, http.StatusNotFound)
			return
		}
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, feedback)
}

func (app *Config) GetFeedbacksForEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	feedbacks, err := app.Models.Rating.GetRatingsForEmployee(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string][]*data.Rating{"feedbacks": feedbacks})
}

func (app *Config) GetFeedbacksForClient(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	feedbacks, err := app.Models.Rating.GetRatingsForClient(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string][]*data.Rating{"feedbacks": feedbacks})
}

func (app *Config) GetFeedbacksForInstitution(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	feedbacks, err := app.Models.Rating.GetRatingsForInstitution(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string][]*data.Rating{"feedbacks": feedbacks})
}

func (app *Config) GetAverageRatingForEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForEmployee(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForInstitution(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForInstitution(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForService(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForService(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForClient(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForClient(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForEmployeeService(w http.ResponseWriter, r *http.Request) {
	employeeId, err := app.readIntParam(r, "employee_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	serviceId, err := app.readIntParam(r, "service_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForEmployeeService(employeeId, serviceId)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForClientService(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "client_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	serviceId, err := app.readIntParam(r, "service_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForClientService(clientId, serviceId)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForClientInstitution(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "client_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	instId, err := app.readIntParam(r, "institution_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForClientInstitution(clientId, instId)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) GetAverageRatingForClientEmployee(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "client_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	employeeId, err := app.readIntParam(r, "employee_id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	avgRating, err := app.Models.Rating.GetAverageRatingForClientEmployee(clientId, employeeId)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]float64{"average_rating": avgRating})
}

func (app *Config) UpdateFeedback(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      int64  `json:"id"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	user, err := app.contextGetUser(r)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	if input.Rating < 1 || input.Rating > 10 {
		app.errorJson(w, data.ErrRatingMustBe1, http.StatusBadRequest)
		return
	}

	feedback := &data.Rating{
		ID:       input.ID,
		Rating:   input.Rating,
		Comment:  input.Comment,
		ClientId: user.ID,
	}

	err = app.Models.Rating.Update(feedback)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Feedback updated successfully"})

}

func (app *Config) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.Models.Rating.Delete(id)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Feedback deleted successfully"})
}
