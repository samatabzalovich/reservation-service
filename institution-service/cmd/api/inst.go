package main

import (
	"context"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
)

func (instService *InstitutionService) CreateInstitution(ctx context.Context, req *inst.CreateInstitutionRequest) (*inst.CreateInstitutionResponse, error) {
	institution := &data.Institution{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Website:     req.GetWebsite(),
		OwnerId:     req.GetOwnerId(),
		Latitude:    req.GetLatitude(),
		Longitude:   req.GetLongitude(),
		Address:     req.GetAddress(),
		Phone:       req.GetPhone(),
		Country:     req.GetCountry(),
		City:        req.GetCity(),
		CategoryId:  req.GetCategoryId(),
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
		CategoryId:  institution.CategoryId,
	}, nil
}

func (instService *InstitutionService) GetInstitutionsByCategory(ctx context.Context, req *inst.GetInstitutionsByCategoryRequest) (*inst.InstitutionsResponse, error) {
	filters, err := data.NewFilters(int(req.GetPageNumber()), int(req.GetPageSize()), req.GetSort(), []string{"id", "rating", "appointments", "employees", "-id", "-rating", "-appointments", "-employees"})
	if err != nil {
		return nil, err
	}
	institutions, metadata, err := instService.Models.Institutions.GetAll(req.CategoryId, filters)
	if err != nil {
		return nil, err
	}
	var institutionsResponse []*inst.Institution
	for _, institution := range institutions {
		institutionsResponse = append(institutionsResponse, &inst.Institution{
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
			CategoryId:  institution.CategoryId,
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

func (instService *InstitutionService) UpdateInstitution(ctx context.Context, req *inst.UpdateInstitutionRequest) (*inst.UpdateInstitutionResponse, error) {
	institution := &data.Institution{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Website:     req.GetWebsite(),
		OwnerId:     req.GetOwnerId(),
		Latitude:    req.GetLatitude(),
		Longitude:   req.GetLongitude(),
		Address:     req.GetAddress(),
		Phone:       req.GetPhone(),
		Country:     req.GetCountry(),
		City:        req.GetCity(),
		CategoryId:  req.GetCategoryId(),
	}
	err := instService.Models.Institutions.Update(institution)
	if err != nil {
		return nil, err
	}
	return &inst.UpdateInstitutionResponse{Id: req.Id}, nil
}

func (instService *InstitutionService) DeleteInstitution(ctx context.Context, req *inst.DeleteInstitutionRequest) (*inst.DeleteInstitutionResponse, error) {
	err := instService.Models.Institutions.Delete(req.GetId())
	if err != nil {
		return nil, err
	}
	return &inst.DeleteInstitutionResponse{Id: req.GetId()}, nil
}

func (instService *InstitutionService) SearchInstitution(ctx context.Context, req *inst.SearchInstitutionsRequest) (*inst.InstitutionsResponse, error) {
	sortSafeList := []string{"id", "rating", "appointments", "employees", "-id", "-rating", "-appointments", "-employees"}
	filter, err := data.NewFilters(int(req.GetPageNumber()), int(req.GetPageSize()), req.GetSort(), sortSafeList)
	if err != nil {
		return nil, err
	}

	institutions, metadata,err := instService.Models.Institutions.Search(req.CategoryId, req.SearchText, filter)
	if err != nil {
		return nil, err
	}
	var institutionsResponse []*inst.Institution
	for _, institution := range institutions {
		institutionsResponse = append(institutionsResponse, &inst.Institution{
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
			CategoryId:  institution.CategoryId,
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
