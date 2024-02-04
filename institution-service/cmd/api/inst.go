package main

import (
	"context"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
)

const (
	// TimeFormat is the format for the time
	TimeFormat = "15:04:00"
)

func (instService *InstitutionService) CreateInstitution(ctx context.Context, req *inst.CreateInstitutionRequest) (*inst.CreateInstitutionResponse, error) {
	workHours := instService.getWorkHoursForRequest(req.GetInstitution().WorkingHours)
	if workHours == nil {
		return &inst.CreateInstitutionResponse{Id: 0}, data.ErrInvalidWorkingHours
	}
	institution, err := data.NewInstitution(
		1,
		req.GetInstitution().Name,
		req.GetInstitution().Description,
		req.GetInstitution().Website,
		req.GetInstitution().OwnerId,
		req.GetInstitution().Latitude,
		req.GetInstitution().Longitude,
		req.GetInstitution().Country,
		req.GetInstitution().City,
		req.GetInstitution().Categories,
		req.GetInstitution().Phone,
		req.GetInstitution().Address,
		workHours,
	)
	if err != nil {
		return &inst.CreateInstitutionResponse{Id: 0}, err
	}
	id, err := instService.Models.Institutions.Insert(institution)
	if err != nil {
		return nil, err
	}
	return &inst.CreateInstitutionResponse{Id: id}, nil
}

func (instService *InstitutionService) GetInstitution(ctx context.Context, req *inst.GetInstitutionsByIdRequest) (*inst.Institution, error) {
	institution, err := instService.Models.Institutions.GetById(req.GetId())
	if err != nil {
		return nil, err
	}
	workHours := instService.getWorkHoursForResponse(institution.WorkingHours)
	return &inst.Institution{
		Id:           institution.ID,
		Name:         institution.Name,
		Description:  institution.Description,
		Website:      institution.Website,
		OwnerId:      institution.OwnerId,
		Latitude:     institution.Latitude,
		Longitude:    institution.Longitude,
		Address:      institution.Address,
		Phone:        institution.Phone,
		Country:      institution.Country,
		City:         institution.City,
		Categories:   institution.Categories,
		WorkingHours: workHours,
	}, nil
}

func (instService *InstitutionService) UpdateInstitution(ctx context.Context, req *inst.UpdateInstitutionRequest) (*inst.UpdateInstitutionResponse, error) {
	workHours := instService.getWorkHoursForRequest(req.GetInstitution().WorkingHours)
	if workHours == nil {
		return &inst.UpdateInstitutionResponse{Id: 0}, data.ErrInvalidWorkingHours
	}
	institution := &data.Institution{
		ID:           req.GetInstitution().Id,
		Name:         req.GetInstitution().Name,
		Description:  req.GetInstitution().Description,
		Website:      req.GetInstitution().Website,
		OwnerId:      req.GetInstitution().OwnerId,
		Latitude:     req.GetInstitution().Latitude,
		Longitude:    req.GetInstitution().Longitude,
		Address:      req.GetInstitution().Address,
		Phone:        req.GetInstitution().Phone,
		Country:      req.GetInstitution().Country,
		City:         req.GetInstitution().City,
		Categories:   req.GetInstitution().Categories,
		WorkingHours: workHours,
	}
	version, err := instService.Models.Institutions.GetVersionByIdForOwner(institution.OwnerId, institution.ID)
	if err != nil {
		return nil, err
	}
	institution.Version = version
	err = instService.Models.Institutions.Update(institution)
	if err != nil {
		return nil, err
	}
	return &inst.UpdateInstitutionResponse{Id: req.Institution.Id}, nil
}

func (instService *InstitutionService) DeleteInstitution(ctx context.Context, req *inst.DeleteInstitutionRequest) (*inst.DeleteInstitutionResponse, error) {
	err := instService.Models.Institutions.Delete(req.GetId())
	if err != nil {
		return nil, err
	}
	return &inst.DeleteInstitutionResponse{Id: req.GetId()}, nil
}

func (instService *InstitutionService) SearchInstitutions(ctx context.Context, req *inst.SearchInstitutionsRequest) (*inst.InstitutionsResponse, error) {
	sortSafeList := []string{"id", "rating", "appointments", "employees", "-id", "-rating", "-appointments", "-employees"}
	filter, err := data.NewFilters(int(req.GetPageNumber()), int(req.GetPageSize()), req.GetSort(), sortSafeList)
	if err != nil {
		return nil, err
	}

	institutions, metadata, err := instService.Models.Institutions.Search(req.Categories, req.SearchText, filter)
	if err != nil {
		return nil, err
	}
	var institutionsResponse []*inst.Institution
	for _, institution := range institutions {
		workHoursResponse := instService.getWorkHoursForResponse(institution.WorkingHours)
		institutionsResponse = append(institutionsResponse, &inst.Institution{
			Id:           institution.ID,
			Name:         institution.Name,
			Description:  institution.Description,
			Website:      institution.Website,
			OwnerId:      institution.OwnerId,
			Latitude:     institution.Latitude,
			Longitude:    institution.Longitude,
			Address:      institution.Address,
			Phone:        institution.Phone,
			Country:      institution.Country,
			City:         institution.City,
			Categories:   institution.Categories,
			WorkingHours: workHoursResponse,
		})
	}
	metadataResponse := &inst.Metadata{
		TotalRecords: int32(metadata.TotalRecords),
		CurrentPage:  int32(metadata.CurrentPage),
		PageSize:     int32(metadata.PageSize),
		FirstPage:    int32(metadata.FirstPage),
		LastPage:     int32(metadata.LastPage),
	}
	return &inst.InstitutionsResponse{Institution: institutionsResponse, Metadata: metadataResponse}, nil
}

func (instService *InstitutionService) isCategoriesIdValid(categories []int64) bool {
	for _, id := range categories {
		if id < 1 {
			return false
		}
	}
	return true
}
