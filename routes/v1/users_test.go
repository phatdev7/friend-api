package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetKaka(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hero"))
}

func TestGetAllUsers(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllUsers)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler return wrong status code: %v want %v", rr.Code, http.StatusOK)
	}
	// if rr.Body.String() != "hero" {
	// 	t.Error("not hero")
	// }
}
