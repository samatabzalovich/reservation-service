package data

import "time"

type EmployeeServices struct {
	Duration time.Duration `json:"duration"`
	Name     string `json:"name"`
}

type ScheduleForEmployeeTimeSlots struct {
	DayOfWeek      int    `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BreakStartTime string `json:"break_start_time"`
	BreakEndTime   string `json:"break_end_time"`
}

type ServiceForEmployeeTimeSlots struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}