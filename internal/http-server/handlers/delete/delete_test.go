package delete_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"URLShortener/internal/http-server/handlers/delete"
	"URLShortener/internal/http-server/handlers/delete/mocks"
	"URLShortener/internal/lib/logger/handlers/slogdiscard"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		respError    string
		expectedJSON string
		mockError    error
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			expectedJSON: `{"status":"OK","alias":"test_alias"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Delete("/url/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req := httptest.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Check the response status code
			assert.Equal(t, http.StatusOK, w.Code)

			//// Check the response JSON
			assert.JSONEq(t, tc.expectedJSON, w.Body.String())

		})
	}
}
