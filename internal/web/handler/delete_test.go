package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/db"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/api"
	"url-shortener/internal/web/handler/mocks"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		funcError string
		mockError error
	}{
		{
			name:  "delete success",
			alias: "abunga",
		},
		{
			name:      "delete fail",
			alias:     "Oleg",
			funcError: ErrURLNotExist.Error,
			mockError: db.ErrURLNotFound,
		},
		{
			name:      "empty alias",
			alias:     "",
			funcError: ErrEmptyAlias.Error,
		},
		{
			name:      "deletion error",
			alias:     "idk",
			funcError: ErrDeletion.Error,
			mockError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockURLDeleter := mocks.NewMockURLDeleter(t)

			if test.funcError == "" || test.mockError != nil {
				mockURLDeleter.On("DeleteURL", test.alias).Return(test.mockError).Once()
			}

			input := fmt.Sprintf(`{"alias":"%s"}`, test.alias)

			req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/url/delete", strings.NewReader(input))
			assert.NoError(t, err)

			router := gin.New()
			router.DELETE("/url/delete", Delete(logger.NewDummyLogger(), mockURLDeleter))

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if test.funcError != "" {
				var resp api.StatusResponse

				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
				assert.Equal(t, test.funcError, resp.Error)
			}
		})
	}
}
