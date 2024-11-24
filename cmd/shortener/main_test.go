package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerPostGet(t *testing.T) {
	type want struct {
		code           int
		response       string
		contentType    string
		headerLocation string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "post_test_#1",
			want: want{
				code:        201,
				response:    `http://example.com/0`,
				contentType: "text/plain",
			},
		},
		{
			name: "get_test_#2",
			want: want{
				code:        307,
				response:    `http://example.com/0`,
				contentType: "http://blabla.ru",
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
				HandlerPostGet(w, request)
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
				request := httptest.NewRequest(http.MethodGet, "/0", nil)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				HandlerPostGet(w, request)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.want.code, res.StatusCode)
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				//resBody, err := io.ReadAll(res.Body)
				//require.NoError(t, err)
				assert.Equal(t, test.want.contentType, res.Header.Get("Location"))
			case "bad_req_test_#3":
				request := httptest.NewRequest(http.MethodDelete, "/", nil)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				HandlerPostGet(w, request)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
			}

		})
	}
}
