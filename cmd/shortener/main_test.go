package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/developerc/reductorUrl/internal/app"
	"github.com/developerc/reductorUrl/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//var settings *config.ServerSettings
//var shortURLAttr *app.ShortURLAttr

func TestPostHandler(t *testing.T) {
	//t.Error("this is error!")
	settings := config.NewServerSettings()
	shortURLAttr := app.NewShortURLAttr(*settings)

	t.Run("#1_PostTest", func(t *testing.T) {
		longURL := strings.NewReader("http://blabla.ru")
		request := httptest.NewRequest(http.MethodPost, "/", longURL)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(app.PostHandler(*shortURLAttr))
		h(w, request)

		result := w.Result()

		assert.Equal(t, 201, result.StatusCode)
		defer result.Body.Close()
		resBody, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		assert.Equal(t, "text/plain", result.Header.Get("Content-Type"))
		assert.Equal(t, `http://localhost:8080/1`, string(resBody))
	})

	t.Run("#2_BadTest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/", nil)
		w := httptest.NewRecorder()
		app.BadHandler(w, request)

		result := w.Result()

		assert.Equal(t, 400, result.StatusCode)
		defer result.Body.Close()
		//assert.Equal(t, "", result.Header.Get("Location"))
	})
}

/*func TestGetHandler(t *testing.T) {

	//settings := config.GetSrvSetGlob()
	shortURLAttr := app.GetShortURLAttr()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/1", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(app.GetHandler(*shortURLAttr))
	h(w, request)

	result := w.Result()

	assert.Equal(t, 307, result.StatusCode)
	defer result.Body.Close()
	assert.Equal(t, "", result.Header.Get("Location"))
}*/

/*func TestHandlerPostGet(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "post_test_#1",
			want: want{
				code:        201,
				response:    `/1`,
				contentType: "text/plain",
			},
		},
		{
			name: "get_test_#2",
			want: want{
				code:        307,
				response:    `http://example.com/1`,
				contentType: "",
			},
		},
		{
			name: "bad_req_test_#3",
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "post_test_#1":
				longURL := strings.NewReader("http://blabla.ru")
				request := httptest.NewRequest(http.MethodPost, "/", longURL)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				PostHandler(w, request)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.want.code, res.StatusCode)
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
				assert.Equal(t, test.want.response, string(resBody))
			case "get_test_#2":
				request := httptest.NewRequest(http.MethodGet, "/1", nil)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				GetHandler(w, request)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.want.code, res.StatusCode)
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				assert.Equal(t, test.want.contentType, res.Header.Get("Location"))
			case "bad_req_test_#3":
				request := httptest.NewRequest(http.MethodDelete, "/", nil)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				BadHandler(w, request)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
			}

		})
	}
}*/
