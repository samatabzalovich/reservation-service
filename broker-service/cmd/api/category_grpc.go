package main

import (
	"context"
	inst "institution-service/proto_files/institution_proto"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) GetCategoriesViaGRpc(w http.ResponseWriter) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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