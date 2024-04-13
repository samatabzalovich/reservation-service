package data

import (
	"database/sql"
	"time"
)

type Appointment struct {
	ID           int64
	ClientName   string
	EmployeeName string
	DeviceToken  string
	PhotoUrl     string
	StartTime    time.Time
}

type AppointmentModel struct {
	DB *sql.DB
}

func (m AppointmentModel) GetUpcomingAppointments() ([]*Appointment, error) {
	//get appointments that starts in 2 hours
	rows, err := m.DB.Query(`SELECT a.id, u.username, e.name, a.start_time, ud.token, s.photo_url
FROM appointments a join users u on a.user_id = u.id join employee e on a.employee_id = e.id 
    join user_devices ud on u.id = ud.user_id join services s on a.service_id = s.id
WHERE a.start_time 
    BETWEEN NOW() AND NOW() + INTERVAL '2 hours' 
  AND a.is_canceled = false AND a.is_notified = false;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(&appointment.ID, &appointment.ClientName, &appointment.EmployeeName, &appointment.StartTime, &appointment.DeviceToken, &appointment.PhotoUrl)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	return appointments, nil
}

func (m AppointmentModel) MarkAsNotified(appointmentId int64) error {
	_, err := m.DB.Exec("UPDATE appointments SET is_notified = true WHERE id = $1", appointmentId)
	if err != nil {
		return err
	}
	return nil
}
