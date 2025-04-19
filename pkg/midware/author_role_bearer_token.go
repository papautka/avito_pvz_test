package midware

import (
	"avito_pvz_test/pkg/jwt"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func CheckRoleByToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearTokenAuth := r.Header.Get("Authorization")

		// проверка наличия Bearer в заголовке
		if bearTokenAuth == "" || !strings.HasPrefix(bearTokenAuth, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		// Извлекаем сам токен (удаляем префикс "Bearer ")
		tokenString := bearTokenAuth[7:]

		// Выводим токен для примера
		fmt.Println("Extracted token:", tokenString)

		// Здесь можно добавить логику для валидации токена
		newJwt := jwt.NewJWT(os.Getenv("TOKEN_MODERATOR"))
		role, err := newJwt.ParseToken(tokenString)
		fmt.Println("ROLE:", role)
		if err != nil || role != "moderator" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		// Переход к следующему обработчику, если токен валиден
		next.ServeHTTP(w, r)
	})
}
