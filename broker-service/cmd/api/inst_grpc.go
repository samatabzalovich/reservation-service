package main

import (
	inst "broker-service/proto_files/institution_proto"
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) CreateInstitutionViaGRpc(w http.ResponseWriter, ctx context.Context, requestPayload RequestPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := c.CreateInstitution(newCtx, &inst.CreateInstitutionRequest{
		Name:        requestPayload.Institution.Name,
		Description: requestPayload.Institution.Description,
		Website:     requestPayload.Institution.Website,
		OwnerId:     requestPayload.Institution.OwnerId,
		Latitude:    requestPayload.Institution.Latitude,
		Longitude:   requestPayload.Institution.Longitude,
		Country:     requestPayload.Institution.Country,
		City:        requestPayload.Institution.City,
		CategoryId:  requestPayload.Institution.CategoryId,
		Phone:       requestPayload.Institution.Phone,
		Address:     requestPayload.Institution.Address,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institution": res.Id,
			"error":       false,
		})
}

func (app *Config) UpdateInstitutionViaGRpc(w http.ResponseWriter, ctx context.Context, requestPayload RequestPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := c.UpdateInstitution(newCtx, &inst.UpdateInstitutionRequest{
		Id:          requestPayload.Institution.Id,
		Name:        requestPayload.Institution.Name,
		Description: requestPayload.Institution.Description,
		Website:     requestPayload.Institution.Website,
		OwnerId:     requestPayload.Institution.OwnerId,
		Latitude:    requestPayload.Institution.Latitude,
		Longitude:   requestPayload.Institution.Longitude,
		Country:     requestPayload.Institution.Country,
		City:        requestPayload.Institution.City,
		CategoryId:  requestPayload.Institution.CategoryId,
		Phone:       requestPayload.Institution.Phone,
		Address:     requestPayload.Institution.Address,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institution": res.Id,
			"error":       false,
		})
}

func (app *Config) DeleteInstitutionViaGRpc(w http.ResponseWriter, ctx context.Context, requestPayload RequestPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := c.DeleteInstitution(newCtx, &inst.DeleteInstitutionRequest{
		Id: requestPayload.Institution.Id,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institution": res.Id,
			"error":       false,
		})
}

func (app *Config) GetInstitutionViaGRpc(w http.ResponseWriter, requestPayload RequestPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitution(ctx, &inst.GetInstitutionsByIdRequest{
		Id: requestPayload.Institution.Id,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institution": res,
			"error":       false,
		})
}

func (app *Config) SearchInstitutionsViaGRpc(w http.ResponseWriter, filterPayload FilterPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.SearchInstitutions(ctx, &inst.SearchInstitutionsRequest{
		PageSize:   int32(filterPayload.PageSize),
		PageNumber: int32(filterPayload.Page),
		SearchText: filterPayload.SearchText,
		Sort:       filterPayload.Sort,
		CategoryId: filterPayload.CategoryId,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institutions": res,
			"error":        false,
		})
}

func (app *Config) GetInstitutionsByCategoryViaGRpc(w http.ResponseWriter, filterPayload FilterPayload) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitutionsByCategory(ctx, &inst.GetInstitutionsByCategoryRequest{
		CategoryId: filterPayload.CategoryId,
		PageSize:   int32(filterPayload.PageSize),
		PageNumber: int32(filterPayload.Page),
		Sort:       filterPayload.Sort,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"institutions": res,
			"error":        false,
		})
}
