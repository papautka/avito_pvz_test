package test

import (
	"avito_pvz_test/cmd"
	"avito_pvz_test/internal/dto/payload"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {

	app := cmd.CreateRouter()
	// тестовый сервер который принимает Handler
	ts := httptest.NewServer(app)
	defer ts.Close()

	data, _ := json.Marshal(&payload.UserAuthRequest{
		Email:    "user1422483@example.com",
		Password: "string",
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
}

func TestLoginAndPasswordFail(t *testing.T) {
	app := cmd.CreateRouter()
	ts := httptest.NewServer(app)

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
