package midware

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/pkg/jwt"
	"net/http"
	"os"
	"strings"
)

func CheckRoleByToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearTokenAuth := r.Header.Get("Authorization")

		// проверка наличия Bearer в заголовке
		if bearTokenAuth == "" || !strings.HasPrefix(bearTokenAuth, "Bearer ") {
			strError := "Доступ запрещен!"
			errorDto.ShowResponseError(&w, strError, http.StatusForbidden)
			return
		}
		// Извлекаем сам токен (удаляем префикс "Bearer ")
		tokenString := bearTokenAuth[7:]

		// Здесь можно добавить логику для валидации токена
		newJwt := jwt.NewJWT(os.Getenv("TOKEN_MODERATOR"))
		role, err := newJwt.ParseToken(tokenString)
		if err != nil {
			msgErr := "Инвалидный Bearer Token"
			errorDto.ShowResponseError(&w, msgErr, http.StatusForbidden, err)
			return
		}
		if role != "moderator" {
			msgErr := "Только пользователь с role moderator может создать PVZ"
			errorDto.ShowResponseError(&w, msgErr, http.StatusForbidden)
			return
		}
		// Переход к следующему обработчику, если токен валиден
		next.ServeHTTP(w, r)
	})
}
