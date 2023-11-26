package delete

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		mockError  error
		respStatus string
		respError  string
	}{
		{
			name:       "Success",
			alias:      "test_alias",
			respStatus: "OK",
			respError:  "",
		},
		{
			name:       "Fail",
			alias:      "test_alias",
			mockError:  errors.New("some error"),
			respStatus: "Error",
			respError:  "internal error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			urlDeleterMock := mocks.NewURLDeleter(t)

			urlDeleterMock.On("DeleteURL", tc.alias).
				Return(tc.mockError).Once()

			r := chi.NewRouter()
			r.Delete("/url/{alias}", New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := api.GetDelete(ts.URL + "/url/" + tc.alias)
			require.NoError(t, err)
			defer resp.Body.Close()
			var response Response
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			require.Equal(t, response.Status, tc.respStatus)
			require.Equal(t, response.Error, tc.respError)
		})
	}
}
