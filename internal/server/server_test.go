package server

import (
	//"net/http"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLink(t *testing.T) {

	th := NewTestHelper()
	t.Run("#1_PostTest", func(t *testing.T) {
		shortURL, err := th.svc.AddLink("http://blabla.ru")
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/1", shortURL)
		//fmt.Println(shortURL)
	})
	/*resp, err := th.srv.Client().Get("http://127.0.0.1:8080/1")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)*/
}
