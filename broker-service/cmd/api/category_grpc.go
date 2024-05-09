package main

import (
	inst "broker-service/proto_files/institution_proto"
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) GetCategoriesViaGRpc(w http.ResponseWriter) {
	conn, err := grpc.Dial(app.instHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewCategoryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitutionCategories(ctx, &inst.GetInstitutionCategoriesRequest{})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"categories": res.Category,
			"error":      false,
		})
}

func (app *Config) CreateCategoryViaGRpc(w http.ResponseWriter, ctx context.Context, requestPayload RequestPayload) {
	conn, err := grpc.Dial(app.instHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := inst.NewCategoryServiceClient(conn)

	res, err := c.CreateCategory(ctx, &inst.InstitutionCategory{
		Name:        requestPayload.Category.Name,
		Description: requestPayload.Category.Description,
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"category_id": res.Id,
			"error":       false,
		})
}
