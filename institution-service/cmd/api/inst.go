package main

import (
	"context"
	"errors"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// TimeFormat is the format for the time
	TimeFormat = "15:04:00"
)

func (instService *InstitutionService) CreateInstitution(ctx context.Context, req *inst.CreateInstitutionRequest) (*inst.CreateInstitutionResponse, error) {
	workHours := instService.getWorkHoursForRequest(req.GetInstitution().WorkingHours)
	if workHours == nil {
		return &inst.CreateInstitutionResponse{Id: 0}, status.Error(codes.InvalidArgument, data.ErrInvalidWorkingHours.Error())
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
		return &inst.CreateInstitutionResponse{Id: 0}, status.Error(codes.InvalidArgument, err.Error())
	}
	id, err := instService.Models.Institutions.Insert(institution)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidCategoryId):
			return &inst.CreateInstitutionResponse{Id: 0}, status.Error(codes.InvalidArgument, err.Error())
		default:
			return &inst.CreateInstitutionResponse{Id: 0}, status.Error(codes.Internal, InvalidServerErr)
		}
	}
	return &inst.CreateInstitutionResponse{Id: id}, nil
}

func (instService *InstitutionService) GetInstitution(ctx context.Context, req *inst.GetInstitutionsByIdRequest) (*inst.Institution, error) {
	institution, err := instService.Models.Institutions.GetById(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get institution")
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
		return &inst.UpdateInstitutionResponse{Id: 0}, status.Error(codes.InvalidArgument, data.ErrInvalidWorkingHours.Error())
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
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	institution.Version = version
	err = instService.Models.Institutions.Update(institution)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, data.ErrInvalidCategoryId):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, InvalidServerErr)
		}
	}
	return &inst.UpdateInstitutionResponse{Id: req.Institution.Id}, nil
}

func (instService *InstitutionService) DeleteInstitution(ctx context.Context, req *inst.DeleteInstitutionRequest) (*inst.DeleteInstitutionResponse, error) {
	err := instService.Models.Institutions.Delete(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	return &inst.DeleteInstitutionResponse{Id: req.GetId()}, nil
}

func (instService *InstitutionService) SearchInstitutions(ctx context.Context, req *inst.SearchInstitutionsRequest) (*inst.InstitutionsResponse, error) {
	sortSafeList := []string{"id", "rating", "appointments", "employees", "-id", "-rating", "-appointments", "-employees"}
	filter, err := data.NewFilters(int(req.GetPageNumber()), int(req.GetPageSize()), req.GetSort(), sortSafeList)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	institutions, metadata, err := instService.Models.Institutions.Search(req.Categories, req.SearchText, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
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

func (instService *InstitutionService) GetForToken(ctx context.Context, req *inst.GetInstForTokenRequest) (*inst.Institution, error) {
	token := req.GetToken()
	if token == "" && len(token) != 26 {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	institution, err := instService.Models.Institutions.GetForToken(data.ScopeEmployeeReg, token)
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	return &inst.Institution{
		Id:          institution.ID,
		Name:        institution.Name,
		Description: institution.Description,
		Website:     institution.Website,
		OwnerId:     institution.OwnerId,
		Latitude:    institution.Latitude,
		Longitude:   institution.Longitude,
		Address:     institution.Address,
		Phone:       institution.Phone,
		Country:     institution.Country,
		City:        institution.City,
	}, nil
}

func (instService *InstitutionService) GetInstitutionsForOwner(ctx context.Context, req *inst.GetInstitutionsByIdRequest) (*inst.InstitutionsResponse, error) {
	institutions, metadata, err := instService.Models.Institutions.GetForOwner(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
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
	metadataRes := &inst.Metadata{
		TotalRecords: int32(metadata.TotalRecords),
		CurrentPage:  int32(metadata.CurrentPage),

		PageSize:  int32(metadata.PageSize),
		FirstPage: int32(metadata.FirstPage),
		LastPage:  int32(metadata.LastPage),
	}
	return &inst.InstitutionsResponse{Institution: institutionsResponse, Metadata: metadataRes}, nil
}
