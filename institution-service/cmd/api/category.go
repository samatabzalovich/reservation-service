package main

import (
	"context"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func(categoryService *CategoryService) CreateCategory(ctx context.Context, req *inst.InstitutionCategory) (*inst.InstitutionCategory, error) {
	user, err := contextGetUser(ctx)
	if err != nil || user.Type != "admin"{
		return nil, status.Error(codes.PermissionDenied, "You do not have permission to create a category")
	}
	photoUrl := req.GetPhotoUrl()
	
	category := &data.Category{
		Name: req.GetName(),
		Description: req.GetDescription(),
		PhotoUrl: &photoUrl,
	}
	id, err := categoryService.Models.Categories.Insert(category)
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	return &inst.InstitutionCategory{
		Id: id,
		Name: category.Name,
		Description: category.Description,

	}, nil
}

func(categoryService *CategoryService) GetInstitutionCategories(ctx context.Context, req *inst.GetInstitutionCategoriesRequest) (*inst.CategoryResponse, error) {
	categories, err := categoryService.Models.Categories.GetAll()
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	var categoriesResponse []*inst.InstitutionCategory
	for _, category := range categories {
		var catphotoUrl string
	if category.PhotoUrl == nil {
		catphotoUrl = ""
	} else {
		catphotoUrl = *category.PhotoUrl
	}
		categoriesResponse = append(categoriesResponse, &inst.InstitutionCategory{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
			PhotoUrl: catphotoUrl,
		})
	}
	return &inst.CategoryResponse{Category: categoriesResponse}, nil
}

func(categoryService *CategoryService) UpdateCategory(ctx context.Context, req *inst.InstitutionCategory) (*inst.InstitutionCategory, error) {
	user, err := contextGetUser(ctx)
	if err != nil || user.Type != "admin"{
		return nil, status.Error(codes.PermissionDenied, "You do not have permission to update a category")
	}
	photoUrl := req.GetPhotoUrl()
	category, err := data.NewCategory(req.GetId(), req.GetName(), req.GetDescription(), photoUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, data.ErrInvalidCategoryId.Error())
	}
	err = categoryService.Models.Categories.Update(category)
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	var catphotoUrl string
	if category.PhotoUrl == nil {
		catphotoUrl = ""
	} else {
		catphotoUrl = *category.PhotoUrl
	}
	return &inst.InstitutionCategory{
		Id: category.ID,
		Name: category.Name,
		Description: category.Description,
		PhotoUrl: catphotoUrl,
	}, nil
}

func(categoryService *CategoryService) DeleteCategory(ctx context.Context, req *inst.InstitutionCategory) (*inst.InstitutionCategory, error) {
	user, err := contextGetUser(ctx)
	if err != nil || user.Type != "admin"{
		return nil, status.Error(codes.PermissionDenied, "You do not have permission to delete a category")
	}
	category, err := categoryService.Models.Categories.GetById(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	err = categoryService.Models.Categories.Delete(category.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	var photoUrl string
	if category.PhotoUrl == nil {
		photoUrl = ""
	} else {
		photoUrl = *category.PhotoUrl
	}
	return &inst.InstitutionCategory{
		Id: category.ID,
		Name: category.Name,
		Description: category.Description,
		PhotoUrl: photoUrl,
	}, nil
}

func(categoryService *CategoryService) GetCategoriesForInstitution(ctx context.Context, req *inst.GetInstitutionsByIdRequest) (*inst.CategoryResponse, error) {
	categories, err := categoryService.Models.Categories.GetByInstitution(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, InvalidServerErr)
	}
	var categoriesResponse []*inst.InstitutionCategory
	for _, category := range categories {
		var photoUrl string
		if category.PhotoUrl == nil {
			photoUrl = ""
		} else {
			photoUrl = *category.PhotoUrl
		}
		categoriesResponse = append(categoriesResponse, &inst.InstitutionCategory{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
			PhotoUrl: photoUrl,
		})
	}
	return &inst.CategoryResponse{Category: categoriesResponse}, nil
}