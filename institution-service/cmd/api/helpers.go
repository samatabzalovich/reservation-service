package main

import (
	"institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
)

func (instService *InstitutionService) getWorkHoursForRequest(workingHours []*inst.WorkingHours) []*data.WorkingHours {
	var workHours []*data.WorkingHours
	for _, wh := range workingHours {
		temp, err := data.NewWorkingHours(
			int(wh.GetDay()),
			wh.GetOpen(),
			wh.GetClose(),
		)
		if err != nil {
			return nil
		}
		workHours = append(workHours, temp)
	}
	return workHours
}

func (instService *InstitutionService) getWorkHoursForResponse(workingHours []*data.WorkingHours) []*inst.WorkingHours {
	var workHours []*inst.WorkingHours
	for _, wh := range workingHours {
		workHours = append(workHours, &inst.WorkingHours{
			Day:   int32(wh.DayOfWeek),
			Open:  wh.OpenTime.Format(TimeFormat),
			Close: wh.CloseTime.Format(TimeFormat),
		})
	}
	return workHours
}
