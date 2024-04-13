package main

import (
	"context"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func(categoryService *CategoryService) CreateCategory(ctx context.Context, req *inst.InstitutionCategory) (*inst.InstitutionCategory, error) {
	user, err := categoryService.contextGetUser(ctx)
	if err != nil || user.Type != "admin"{
		return nil, status.Error(codes.PermissionDenied, "You do not have permission to create a category")
	}
	category := &data.Category{
		Name: req.GetName(),
		Description: req.GetDescription(),
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
		categoriesResponse = append(categoriesResponse, &inst.InstitutionCategory{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
		})
	}
	return &inst.CategoryResponse{Category: categoriesResponse}, nil
}