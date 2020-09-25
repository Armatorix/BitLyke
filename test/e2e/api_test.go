package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	api = model.NewAPIClient(&model.Configuration{
		Host:   "localhost:8081",
		Scheme: "http",
		HTTPClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}).DefaultApi
	ctx = context.Background()
)

func TestAPI(t *testing.T) {

	resp, err := api.PublicHealthCheckGet(ctx)
	if err != nil || resp.StatusCode != 200 {
		t.Fatal("API health check failed")
	}
	t.Run("remove shorts", cleanupDB)
	t.Run("basic flow", testBasicFlow)
	t.Run("duplicated inserts/deletions", testDuplicatedShorts)
	t.Run("forbidden special cases", testForbiddenCases)
	t.Run("remove shorts", cleanupDB)
}

func testBasicFlow(t *testing.T) {
	tts := []model.ShortLink{
		{
			ShortPath: "test1",
			RealUrl:   "https://www.google.com",
		},
		{
			ShortPath: "test2",
			RealUrl:   "https://www.yandex.com",
		},
		{
			ShortPath: "dummytest",
			RealUrl:   "https://duckduckgo.com/",
		},
		{
			ShortPath: "someurlwithparams",
			RealUrl:   "https://www.youtube.com/watch?v=ngf1KF2_kPI",
		},
	}

	for i, tt := range tts {
		sl, resp, err := api.ApiPost(ctx, tt)
		assert.NoError(t, err, tt)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, tt)
		assert.Equal(t, tt, sl)

		for _, alreadyShorted := range tts[:i+1] {
			resp, err := api.LidGet(ctx, alreadyShorted.ShortPath)
			assert.EqualError(t, err, "307 Temporary Redirect")
			assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		}
	}

	for i, tt := range tts {
		resp, err := api.ApiLidDelete(ctx, tt.ShortPath)
		assert.NoError(t, err, tt)
		assert.Equal(t, http.StatusOK, resp.StatusCode, tt)

		for _, deleted := range tts[:i+1] {
			resp, err := api.LidGet(ctx, deleted.ShortPath)
			assert.Error(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		}

		for _, stillShorted := range tts[i+1:] {
			resp, err := api.LidGet(ctx, stillShorted.ShortPath)
			assert.EqualError(t, err, "307 Temporary Redirect")
			assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		}
	}
}

func testDuplicatedShorts(t *testing.T) {
	tts := []model.ShortLink{
		{
			ShortPath: "test1",
			RealUrl:   "https://www.google.com",
		},
		{
			ShortPath: "test2",
			RealUrl:   "https://www.yandex.com",
		},
		{
			ShortPath: "dummytest",
			RealUrl:   "https://duckduckgo.com/",
		},
		{
			ShortPath: "someurlwithparams",
			RealUrl:   "https://www.youtube.com/watch?v=ngf1KF2_kPI",
		},
	}

	for _, tt := range tts {
		sl, resp, err := api.ApiPost(ctx, tt)
		assert.NoError(t, err, tt)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, tt)
		assert.Equal(t, tt, sl)

		resp, err = api.LidGet(ctx, tt.ShortPath)
		assert.EqualError(t, err, "307 Temporary Redirect")
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

		_, resp, err = api.ApiPost(ctx, tt)
		assert.Error(t, err, tt)
		assert.Equal(t, http.StatusConflict, resp.StatusCode, tt)

		resp, err = api.LidGet(ctx, tt.ShortPath)
		assert.EqualError(t, err, "307 Temporary Redirect")
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

		resp, err = api.ApiLidDelete(ctx, tt.ShortPath)
		assert.NoError(t, err, tt)
		assert.Equal(t, http.StatusOK, resp.StatusCode, tt)

		resp, err = api.LidGet(ctx, tt.ShortPath)
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		resp, err = api.ApiLidDelete(ctx, tt.ShortPath)
		assert.Error(t, err, tt)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, tt)

		resp, err = api.LidGet(ctx, tt.ShortPath)
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}
func testForbiddenCases(t *testing.T) {

}
func cleanupDB(t *testing.T) {
	sls, resp, err := api.ApiGet(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	for _, ls := range sls {
		resp, err = api.ApiLidDelete(ctx, ls.ShortPath)
		assert.NoError(t, err, ls)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	sls, resp, err = api.ApiGet(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, sls, 0)
}
