package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/stretchr/testify/assert"
)

const redirectErrorMsg = "307 Temporary Redirect"

var (
	api = model.NewAPIClient(&model.Configuration{
		Host:   "localhost:8080",
		Scheme: "http",
		HTTPClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}).DefaultApi
	ctx = context.TODO()
	tts = []model.ShortLink{
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
		{
			ShortPath: "apitest",
			RealUrl:   "http://letmegooglethat.com/?q=bitlyke",
		},
	}
	forbiddenTts = []model.ShortLink{
		{
			ShortPath: "api",
			RealUrl:   "https://www.google.com/short-path-blocked-by-internal-path",
		},
		{
			ShortPath: "counts",
			RealUrl:   "http://localhost/short-path-blocked-by-internal-path",
		},
		{
			ShortPath: "some/dummy/path",
			RealUrl:   "https://hackerrank.com/short-path-contains-permited-sign",
		},
	}
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
	t.Run("remove shorts", testSuitCleanUpDB)
	t.Run("basic flow", testSuitBasicFlow)
	t.Run("duplicated inserts/deletions", testSuitDuplicatedShorts)
	t.Run("forbidden special cases", testSuitForbiddenCases)
	t.Run("remove shorts", testSuitCleanUpDB)
}

func testSuitBasicFlow(t *testing.T) {
	for i, tt := range tts {
		testCorrectPathCreate(t, tt)
		for _, alreadyShorted := range tts[:i+1] {
			testCorrectRedirection(t, alreadyShorted)
		}
	}

	for i, tt := range tts {
		testCorrectDelete(t, tt)
		for _, deleted := range tts[:i+1] {
			testNotFoundRedirection(t, deleted)
		}
		for _, stillShorted := range tts[i+1:] {
			testCorrectRedirection(t, stillShorted)
		}
	}
}

func testSuitDuplicatedShorts(t *testing.T) {
	for _, tt := range tts {
		testCorrectPathCreate(t, tt)
		testCorrectRedirection(t, tt)

		testDuplicatedPathCreate(t, tt)
		testCorrectRedirection(t, tt)

		testCorrectDelete(t, tt)
		testNotFoundRedirection(t, tt)
		testNotFoundDelete(t, tt)
	}
}

func testSuitForbiddenCases(t *testing.T) {
	for _, tt := range forbiddenTts {
		testBadPathCreate(t, tt)
		testNotRedirected(t, tt)
	}
}

func testSuitCleanUpDB(t *testing.T) {
	sls := testCorrectGetAll(t)

	for _, ls := range sls {
		testCorrectDelete(t, ls)
	}

	sls = testCorrectGetAll(t)
	assert.Len(t, sls, 0)
}

func testCorrectRedirection(t *testing.T, tt model.ShortLink) {
	resp, err := api.LinkIdGet(ctx, tt.ShortPath)
	resp.Body.Close()
	assert.EqualError(t, err, redirectErrorMsg)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	redirURL, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, tt.RealUrl, redirURL.String())
}

func testNotFoundRedirection(t *testing.T, tt model.ShortLink) {
	resp, err := api.LinkIdGet(ctx, tt.ShortPath)
	resp.Body.Close()
	assert.Error(t, err, tt)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, tt)
}

func testNotRedirected(t *testing.T, tt model.ShortLink) {
	resp, _ := api.LinkIdGet(ctx, tt.ShortPath)
	resp.Body.Close()
	assert.NotEqual(t, http.StatusTemporaryRedirect, resp.StatusCode, tt)
}

func testCorrectPathCreate(t *testing.T, tt model.ShortLink) {
	sl, resp, err := api.ApiPost(ctx, tt)
	resp.Body.Close()
	assert.NoError(t, err, tt)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, tt)
	assert.Equal(t, tt, sl)
}

func testDuplicatedPathCreate(t *testing.T, tt model.ShortLink) {
	_, resp, err := api.ApiPost(ctx, tt)
	resp.Body.Close()
	assert.Error(t, err, tt)
	assert.Equal(t, http.StatusConflict, resp.StatusCode, tt)
}

func testBadPathCreate(t *testing.T, tt model.ShortLink) {
	_, resp, err := api.ApiPost(ctx, tt)
	resp.Body.Close()
	assert.Error(t, err, tt)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tt)
}

func testCorrectDelete(t *testing.T, tt model.ShortLink) {
	resp, err := api.ApiLinkIdDelete(ctx, tt.ShortPath)
	resp.Body.Close()
	assert.NoError(t, err, tt)
	assert.Equal(t, http.StatusOK, resp.StatusCode, tt)
}

func testNotFoundDelete(t *testing.T, tt model.ShortLink) {
	resp, err := api.ApiLinkIdDelete(ctx, tt.ShortPath)
	resp.Body.Close()
	assert.Error(t, err, tt)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, tt)
}

func testCorrectGetAll(t *testing.T) []model.ShortLink {
	sls, resp, err := api.ApiGet(ctx)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	return sls
}
