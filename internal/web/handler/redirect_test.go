package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/api"
	"url-shortener/internal/web/handler/mocks"
)

func TestRedirect(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		funcError string
		mockError error
	}{
		{
			name:  "success",
			alias: "123",
		},
		{
			name:      "error",
			alias:     "1 23",
			funcError: ErrInvalidAlias.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mockURLGetter := mocks.NewMockURLGetter(t)

			if test.mockError != nil || test.funcError == "" {
				mockURLGetter.On("GetURL", test.alias).
					Return(test.alias, test.mockError).
					Once()
			}

			router := gin.New()
			router.GET("/url/:alias", Redirect(logger.NewDummyLogger(), mockURLGetter))

			req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/url/"+test.alias, nil)
			require.NoError(t, err)

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
