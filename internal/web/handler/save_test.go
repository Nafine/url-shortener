package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/db"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/handler/mocks"
)

func TestSave(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		alias     string
		funcError string
		mockError error
	}{
		{
			name:  "success",
			url:   "http://test.com",
			alias: "alias",
		},
		{
			name: "empty alias",
			url:  "http://example.com/",
		},
		{
			name:      "invalid url",
			url:       "123",
			funcError: "field Url is an invalid URL",
		},
		{
			name:      "empty url",
			funcError: "field Url is a required field",
		},
		{
			name:      "save fail",
			url:       "http://test.com",
			alias:     "alias",
			funcError: ErrFailSaveURL.Error,
			mockError: errors.New("unexpected error"),
		},
		{
			name:      "existing url",
			url:       "http://existing.com",
			funcError: ErrURLExists.Error,
			mockError: db.ErrURLAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mockURLSaver := mocks.NewMockURLSaver(t)

			if test.funcError == "" || test.mockError != nil {
				mockURLSaver.On(
					"SaveURL",
					test.url,
					mock.AnythingOfType("string"),
				).Return(test.mockError).Once()
			}

			input := fmt.Sprintf(
				`{"url": "%s", "alias": "%s"}`,
				test.url,
				test.alias,
			)

			req, err := http.NewRequest(http.MethodPost, "/url", strings.NewReader(input))
			require.NoError(t, err)

			router := gin.New()
			router.POST("/url", Save(logger.NewDummyLogger(), mockURLSaver))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var resp Response

			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, test.funcError, resp.Error)
		})
	}
}
