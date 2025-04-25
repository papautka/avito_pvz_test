package main

import (
	"avito_pvz_test/internal/dto/payload"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {

	app := CreateRouter()
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

}
