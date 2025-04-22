package midware

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/pkg/jwt"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func CheckRoleByToken(next http.Handler, strRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := GetRoleFromToken(r, strRole)
		if err != nil {
			errorDto.ShowResponseError(&w, "Ошибка авторизации", http.StatusForbidden, err)
			return
		}

		if role != strRole {
			msgErr := "Только пользователь с ролью: " + role + " может внести изменения"
			errorDto.ShowResponseError(&w, msgErr, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetRoleFromToken(r *http.Request, strRole string) (string, error) {
	bearTokenAuth := r.Header.Get("Authorization")

	if bearTokenAuth == "" || !strings.HasPrefix(bearTokenAuth, "Bearer ") {
		return "", fmt.Errorf("доступ запрещен: нет Bearer токена")
	}

	tokenString := bearTokenAuth[7:]

	TokenStr := ""
	switch strRole {
	case "client":
		TokenStr = "TOKEN_CLIENT"
	case "moderator":
		TokenStr = "TOKEN_MODERATOR"
	}

	j := jwt.NewJWT(os.Getenv(TokenStr))
	role, err := j.ParseToken(tokenString)

	if err != nil {
		return "", fmt.Errorf("невалидный токен: %w", err)
	}
	return role, nil
}
