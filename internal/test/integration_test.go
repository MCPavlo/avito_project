package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"tz/internal/handler"
)

func TestDummyLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dummyLogin" {
			handler.NewDummyLoginHandler().Login(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	employeeToken := getToken(t, ts, "employee")
	if employeeToken == "" {
		t.Fatalf("Failed to get employee token")
	}
	moderatorToken := getToken(t, ts, "moderator")
	if moderatorToken == "" {
		t.Fatalf("Failed to get moderator token")
	}

	t.Log("✅ DummyLogin test passed successfully!")
}

// --------------------  Утилиты --------------------

// getToken отправляет запрос на /dummyLogin и возвращает полученный токен
func getToken(t *testing.T, ts *httptest.Server, role string) string {
	body := map[string]string{"role": role}
	reqBody, _ := json.Marshal(body)

	resp, err := http.Post(fmt.Sprintf("%s/dummyLogin", ts.URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to get %s token: %v", role, err)
	}
	defer resp.Body.Close()

	raw, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Raw response from /dummyLogin:", string(raw))

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
	}
	var tokenResp map[string]string
	if err := json.Unmarshal(raw, &tokenResp); err != nil {
		t.Fatalf("Failed to parse %s token response: %v", role, err)
	}
	return tokenResp["token"]
}
