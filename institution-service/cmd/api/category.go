package main

import (
	"context"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
)


func(categoryService *CategoryService) CreateCategory(ctx context.Context, req *inst.InstitutionCategory) (*inst.InstitutionCategory, error) {
	category := &data.Category{
		Name: req.GetName(),
		Description: req.GetDescription(),
	}
	id, err := categoryService.Models.Categories.Insert(category)
	if err != nil {
		return nil, err
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
		return nil, err
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