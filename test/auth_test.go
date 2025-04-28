package test

import (
	"avito_pvz_test/cmd"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/internal/users"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app = cmd.CreateRouterTest()

func TestLoginSuccess(t *testing.T) {

	// тестовый сервер который принимает Handler
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	app.Repo.UsersRepo.CreateUser(&users.User{
		Id:       uuid.New(),
		Email:    "example1@example.com",
		Role:     "client",
		Password: "password1",
	})

	data, _ := json.Marshal(&payload.UserAuthRequest{
		Email:    "example1@example.com",
		Password: "password1",
	})

	resp, err := http.Post(ts.URL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d", 200, resp.StatusCode)
	}
	req_data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var tokenResponse payload.TokenResponse
	err = json.Unmarshal(req_data, &tokenResponse)
	if err != nil {
		t.Fatal(err)
	}
	if tokenResponse.Token == "" {
		t.Fatal("Token empty", err)
	}
	t.Cleanup(func() {
		_ = app.Repo.UsersRepo.DropUser("example1@example.com")
		log.Println("drop user example1@example.com")
	})

}

func TestLoginAndPasswordFail(t *testing.T) {
	ts := httptest.NewServer(app.Router)

	defer ts.Close()
	data, _ := json.Marshal(&payload.UserAuthRequest{
		Email:    "aboba@example.com",
		Password: "jopa",
	})
	resp, err := http.Post(ts.URL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal("Пользователя с таким именем нет в базе данных")
	}
}
