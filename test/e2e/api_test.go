package main

import (
	"context"
	"testing"

	"github.com/Armatorix/BitLyke/pkg/model"
)

func Test_API(t *testing.T) {
	api := model.NewAPIClient(&model.Configuration{
		Host:   "localhost:8081",
		Scheme: "http",
	}).DefaultApi
	ctx := context.Background()

	resp, err := api.PublicHealthCheckGet(ctx)
	if err != nil || resp.StatusCode != 200 {
		t.Skip("API healthcheck failed")
	}

}
