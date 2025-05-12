// general пакет для размещения типов данных, общих для нескольких пакетов.
// Служит для ухода от перекрестных зависимостей.
package general

import (
	"errors"
	"sync/atomic"
)

// ArrLongURL структура списка длинных URL.
type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ArrRepoURL структура списка URL присланных пользователем.
type ArrRepoURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CntrAtom структура для учета открытых горутин
type CntrAtom struct {
	counter    atomic.Int64
	chCntrAtom chan bool
}

// NewCntrAtom конструктор CntrAtomVar
func NewCntrAtom() {
	CntrAtomVar = CntrAtom{}
	CntrAtomVar.chCntrAtom = make(chan bool)
}

// IncrCntr увеличиваем счетчик
func (ca *CntrAtom) IncrCntr() {
	ca.counter.Add(1)
}

// DecrCntr уменьшаем счетчик
func (ca *CntrAtom) DecrCntr() {
	ca.counter.Add(-1)
}

// GetCntr получаем счетчик
func (ca *CntrAtom) GetCntr() int64 {
	return ca.counter.Load()
}

// GetChan получаем канал
func (ca *CntrAtom) GetChan() chan bool {
	return ca.chCntrAtom
}

// SentNotif отсылаем уведомление
func (ca *CntrAtom) SentNotif() {
	ca.chCntrAtom <- true
}

// CntrAtomVar глобальная переменная атомарного счетчика запущенных горутин
var CntrAtomVar CntrAtom

// ArrGetStats структура для статистики
type ArrGetStats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// ArrShortURL структура массива коротких URL
type ArrShortURL struct {
	CorellationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ErrorURLExists структура типизированной ошибки существования длинного URL
type ErrorURLExists struct {
	Str string
}

// Error возвращает строку со значением ошибки существования длинного URL
func (e *ErrorURLExists) Error() string {
	return e.Str
}

// AsURLExists проверяет существование длинного URL
func (e *ErrorURLExists) AsURLExists(err error) bool {
	return errors.As(err, &e)
}

// User структура пользователя
type User struct {
	Name string
}
