package memory

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

type User struct {
	Name string
}

var secure *securecookie.SecureCookie

func InitSecure() {
	var hashKey = []byte("very-secret-qwer")
	var blockKey = []byte("a-lot-secret-qwe")
	secure = securecookie.New(hashKey, blockKey)
}

func (s *Service) SetCookie(usr string) (*http.Cookie, error) {
	var cookie *http.Cookie
	u := &User{
		Name: usr,
	}
	if encoded, err := secure.Encode("user", u); err == nil {
		cookie = &http.Cookie{
			Name:  "user",
			Value: encoded,
		}
		return cookie, nil
	} else {
		return nil, err
	}
}

func (s *Service) ReadCookie(r *http.Request) (string, error) {
	var err error
	if cookie, err := r.Cookie("user"); err == nil {
		u := &User{}
		if err = secure.Decode("user", cookie.Value, u); err == nil {
			return u.Name, nil
		}
	}
	return "", err
}

func (s *Service) ReadCookie2(cookieValue string) (string, error) {
	var err error
	u := &User{}
	if err = secure.Decode("user", cookieValue, u); err == nil {
		return u.Name, nil
	}
	return "", err
}
