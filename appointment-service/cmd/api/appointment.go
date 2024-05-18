package main

import (
	"appointment-service/internal/data"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *Config) GetAppointmentById(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointment, err := app.Models.Appointments.GetById(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"appointment": appointment})
}

func (app *Config) GetAppointmentsForInstitution(w http.ResponseWriter, r *http.Request) {
	institutionId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointments, err := app.Models.Appointments.GetAllForInst(institutionId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (app *Config) GetAvailableTimeSlots(w http.ResponseWriter, r *http.Request) {
	var (
		employeeID int64
		serviceID  int64
	)

	var selectedDay time.Time

	selectedDay, err := app.readTimeParam("selected_day", r.URL.Query())
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	employeeID, err = app.readIntParam(r, "employee_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	serviceID, err = strconv.ParseInt(app.readString(r.URL.Query(), "service_id", ""), 10, 64)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	slots, err := app.ReturnAvailableTimeSlots(employeeID, serviceID, selectedDay)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			app.errorJson(w, ErrNotFound, http.StatusNotFound)
			return
		}
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"slots": slots})
}

func (app *Config) ReturnAvailableTimeSlots(employeeID, serviceID int64, selectedDay time.Time) ([]*data.SlotResponse, error) {
	var response struct {
		Schedule data.ScheduleForEmployeeTimeSlots `json:"schedule"`
		Service  data.ServiceForEmployeeTimeSlots  `json:"service"`
	}
	endpoint := fmt.Sprintf("http://%s%s%d/%d?selected_day=%s", app.staffServiceHost, "schedule/", employeeID, serviceID, selectedDay.Format(time.RFC3339))
	err := app.sendGetRequest(endpoint, &response)
	if err != nil {
		return nil, err
	}
	appointments, err := app.Models.Appointments.GetAllForEmployee(employeeID)

	if err != nil {
		return nil, err
	}

	var output struct {
		Schedule data.EmployeeSchedule
		Service  data.EmployeeServices
	}

	err = output.Schedule.ParseToTimeStruct(response.Schedule.DayOfWeek, response.Schedule.StartTime, response.Schedule.EndTime, response.Schedule.BreakStartTime, response.Schedule.BreakEndTime)
	if err != nil {
		return nil, err
	}

	output.Service.Duration, err = time.ParseDuration(response.Service.Duration + "s")
	if err != nil {
		return nil, err
	}

	
	availableTimeSlots := []*data.SlotResponse{}
	// add all time slots of the day for the employee
	allTimeSlots := app.distributeTimeSlots(output.Schedule.StartTime, output.Schedule.EndTime, output.Schedule.BreakStartTime, output.Schedule.BreakEndTime, output.Service.Duration)

	// take only the time slots that are not already taken and add them to availableTimeSlots
	for _, timeSlot := range allTimeSlots {
		isAvailable := true
		for _, appointment := range appointments {
			if appointment.StartTime.Time.Local().Hour() == (timeSlot.StartTime.Hour()) && appointment.StartTime.Time.Local().Minute() == (timeSlot.StartTime.Minute()) {
				isAvailable = false
				break
			}
		}
		if isAvailable {
			availableTimeSlots = append(availableTimeSlots, timeSlot)
		}
	}
	return availableTimeSlots, nil
}
func(app *Config) distributeTimeSlots(start, end, breakStart, breakEnd time.Time, serviceDuration time.Duration) []*data.SlotResponse {
    totalDuration := end.Sub(start) - breakEnd.Sub(breakStart)
    totalMinutes := int(totalDuration.Minutes())
    numSlots := totalMinutes / int(serviceDuration.Minutes())
    gapTime := totalMinutes % int(serviceDuration.Minutes())
    var gapPerSlot int 
	if numSlots == 1 {
		gapPerSlot = 0
	} else {
		gapPerSlot = gapTime / numSlots
	}

    slots := make([]*data.SlotResponse, 0, numSlots) 
    currentTime := start
    for i := 0; i < numSlots; i++ {
		//check if the current time is in the break time or equal to the break time
		var slotStart time.Time
		if (app.isValidSlotStart(currentTime, breakStart, breakEnd, serviceDuration)) {
			slotStart = currentTime
		} else {
			currentTime = currentTime.Add(breakEnd.Sub(breakStart))
			slotStart = currentTime
		}
        currentTime = currentTime.Add(serviceDuration)
        slotEnd := currentTime
        slot := &data.SlotResponse{ // Create a pointer to a new SlotResponse
			StartTime: slotStart,
			EndTime:   slotEnd,
		}
		currentTime = currentTime.Add(time.Duration(gapPerSlot) * time.Minute)
		slots = append(slots, slot)
    }

    return slots
}

func (app *Config) isValidSlotStart(currentTime time.Time, breakStart, breakEnd time.Time, serviceDuration time.Duration) bool {
	

	if (currentTime.After(breakStart) && currentTime.Before(breakEnd)) || currentTime.Equal(breakStart) || currentTime.Equal(breakEnd)  {
		return false
	}
	slotEnd := currentTime.Add(serviceDuration)
	if (slotEnd.After(breakStart) && slotEnd.Before(breakEnd))  || slotEnd.Equal(breakEnd) {
		return false
	}
	return true
}


func (app *Config) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClientID   int           `json:"client_id"`
		InstId     int           `json:"inst_id"`
		EmployeeID int           `json:"employee_id"`
		ServiceID  int           `json:"service_id"`
		StartTime  data.DateTime `json:"start_time"`
		EndTime    data.DateTime `json:"end_time"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	slots, err := app.ReturnAvailableTimeSlots(int64(input.EmployeeID), int64(input.ServiceID), input.StartTime.Time)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			slots = nil
			err = nil
		} else {
			app.errorJson(w, err)
			return
		}

	}

	appointment, err := data.NewAppointment(
		input.ClientID,
		input.InstId,
		input.EmployeeID,
		input.ServiceID,
		input.StartTime,
		input.EndTime,
		false,
		slots,
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Appointments.Insert(appointment)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	// TODO: send notification test
	go func() { 
		
		if err := app.SendAppointmentNotification(*appointment); err != nil {
			log.Printf("Error sending appointment notification: %v", err) 
		}
	}()
	
	app.writeJSON(w, http.StatusCreated, map[string]any{"id": appointment.ID})
}

func (app *Config) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment struct {
		ID          int64    `json:"id"`
		StartTime   data.DateTime `json:"start_time"`
		EndTime     data.DateTime `json:"end_time"`
		IsCancelled bool     `json:"is_cancelled"`
	}
	err := app.readJSON(w, r, &appointment)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	appointmentData, err := app.Models.Appointments.GetById(appointment.ID)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if appointmentData.IsCancelled {
		app.errorJson(w, data.ErrInvalidAppointment, http.StatusForbidden)
		return
	}

	// check if appointment starts in 2 hours
	if time.Now().Add(time.Hour * 2).After(appointmentData.StartTime.Time.Local()) {
		app.errorJson(w, data.ErrAppointmentStartsIn2Hours, http.StatusForbidden)
		return
	}
	if !((appointmentData.StartTime.Time.Equal(appointment.StartTime.Time) && appointmentData.EndTime.Time.Equal(appointment.EndTime.Time)) || (appointment.StartTime.Time.String() == "0001-01-01 00:00:00 +0000 UTC")) {
		slots,err := app.ReturnAvailableTimeSlots( int64(appointmentData.EmployeeID), int64(appointmentData.ServiceID), appointment.StartTime.Time)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				app.errorJson(w, data.ErrInvalidAppointment, http.StatusNotFound)
				return
			} else {
				app.errorJson(w, err)
				return
			}
		}
		isValid := false
		for _, slot := range slots {
			if app.EqualTime(slot.StartTime, appointment.StartTime.Time) && app.EqualTime(slot.EndTime, appointment.EndTime.Time) {
				isValid = true
				break
			}
		}
		if !isValid {
			app.errorJson(w, data.ErrInvalidAppointment, http.StatusForbidden)
			return
		}
	} else {
		appointment.StartTime = appointmentData.StartTime
		appointment.EndTime = appointmentData.EndTime
	}
	err = app.Models.Appointments.Update(
		&data.Appointment{
			ID: appointment.ID,
			StartTime: appointment.StartTime,
			EndTime: appointment.EndTime,
			IsCancelled: appointment.IsCancelled,
			ClientID: appointmentData.ClientID,
			EmployeeID: appointmentData.EmployeeID,
			ServiceID: appointmentData.ServiceID,

		},
	)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "appointment updated"})
}

func (app *Config) GetAppointmentsForEmployee(w http.ResponseWriter, r *http.Request) {
	employeeId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointments, err := app.Models.Appointments.GetAllForEmployee(employeeId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (app *Config) GetAppointmentsForClient(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointments, err := app.Models.Appointments.GetAllForClient(clientId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (app *Config) EqualTime(t1, t2 time.Time) bool {
	return t1.Hour() == t2.Hour() && t1.Minute() == t2.Minute()
}


func (app *Config) GetNumberOfCompletedAppointmentsForUser(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "clientId")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	var employeeIdInt int64
	employeeIdString := app.readString(r.URL.Query(), "employeeId", "")
	if employeeIdString == "" {
		employeeIdInt = 0
	} else {
		employeeIdInt, err = strconv.ParseInt(employeeIdString, 10, 64)
		if err != nil {
			app.errorJson(w, err)
			return
		}
	}

	instIdString := app.readString(r.URL.Query(), "instId", "")

	if instIdString == "" {
		app.errorJson(w, errors.New("missing instId"), http.StatusBadRequest)
		return
	}

	instId, err := strconv.ParseInt(instIdString, 10, 64)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	count, err := app.Models.Appointments.GetNumberOfCompletedAppointmentsForUser(clientId, instId, employeeIdInt)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"count": count})
}