package memory

import (
	"fmt"
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
	fmt.Println("from ReadCookieHandler")
	var err error
	if cookie, err := r.Cookie("user"); err == nil {
		u := &User{}
		if err = secure.Decode("user", cookie.Value, u); err == nil {
			return u.Name, nil
		}
	}
	return "", err
}

/*func (s *Service) GetCripto() (string, error) {
	src := []byte("Ключ от сердца") // данные, которые хотим зашифровать
	fmt.Printf("original: %s\n", src)

	// будем использовать AES-256, создав ключ длиной 32 байта

	key := sha256.Sum256([]byte("super_secret****"))
	// NewCipher создает и возвращает новый cipher.Block.
	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
	// для выбора AES-128, AES-192 или AES-256.
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
	fmt.Printf("encrypted: %x\n", dst)
	return hex.EncodeToString(dst), nil
}

func generateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}*/
