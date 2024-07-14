package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenRequestIsCorrect(t *testing.T) {
	// Пример корректного запроса
	req := httptest.NewRequest("GET", "/cafe?count=3&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Полученный статус равен ожидаемому. Если нет, смысла проверять дальше нет
	require.Equal(t, responseRecorder.Code, http.StatusOK)

	// Тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body)
}

// Город, который передаётся в параметре `city`, не поддерживается. Сервис возвращает код ответа 400 и ошибку `wrong city value` в теле ответа
func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	// Ожидаемый ответ
	expectedAnswer := "wrong city value"

	// Пример запроса с неправильным городом
	req := httptest.NewRequest("GET", "/cafe?count=3&city=moskva", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Полученный статус равен ожидаемому. Если нет, смысла проверять дальше нет
	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)

	// Ответ соответствует ожидаемому
	body := responseRecorder.Body.String()
	assert.Equal(t, body, expectedAnswer)
}

// Если в параметре `count` указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	// Ожидаемая длина
	totalCount := 4

	// Запрос к сервису с превышением длины
	req := httptest.NewRequest("GET", "/cafe?count=25&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Полученный статус равен ожидаемому. Если нет, смысла проверять дальше нет
	require.Equal(t, responseRecorder.Code, http.StatusOK)
	// Тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body)

	// Длина ответа соответствует ожидаемому
	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")
	assert.Len(t, list, totalCount)
}
