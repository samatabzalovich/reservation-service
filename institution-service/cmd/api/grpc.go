package main

import (
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
)

var (
	InvalidServerErr = "invalid server error"
)

type InstitutionService struct {
	inst.UnimplementedInstitutionServiceServer
	Models data.Models
}

type CategoryService struct {
	inst.UnimplementedCategoryServiceServer
	Models data.Models
}

