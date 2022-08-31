package main

import (
	"final-project/data"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_Pages(t *testing.T) {
	pathToTemplate = "./templates"

	tests := []struct {
		name               string
		url                string
		expectedStatusCode int
		handler            http.HandlerFunc
		sessionData        map[string]any
		expectedHTML       string
	}{
		{
			name:               "home",
			url:                "/",
			expectedStatusCode: http.StatusOK,
			handler:            testApp.HomePage,
		},
		{
			name:               "login",
			url:                "/login",
			expectedStatusCode: http.StatusOK,
			handler:            testApp.LoginPage,
			expectedHTML:       `<h1 class="mt-5">Login</h1>`,
		},
		{
			name:               "logout",
			url:                "/logout",
			expectedStatusCode: http.StatusSeeOther,
			handler:            testApp.LoginPage,
			sessionData: map[string]any{
				"userID": 1,
				"user":   data.User{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			ctx := GetCtx(req)
			req = req.WithContext(ctx)

			if len(tt.sessionData) > 0 {
				for k, v := range tt.sessionData {
					testApp.Session.Put(ctx, k, v)
				}
			}
			tt.handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected %d yet get %d", http.StatusOK, rr.Code)
			}

			if len(tt.expectedHTML) > 0 {
				html := rr.Body.String()
				if !strings.Contains(html, tt.expectedHTML) {
					t.Errorf("expected %s yet get %s", tt.expectedHTML, html)
				}
			}
		})

	}

}

func TestConfig_PostLoginPage(t *testing.T) {
	pathToTemplate = "./templates"

	postedData := url.Values{
		"email":    {"admin@example.com"},
		"password": {"abc123"},
	}
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(postedData.Encode()))
	ctx := GetCtx(req)

	req = req.WithContext(ctx)
	handler := http.HandlerFunc(testApp.PostLoginPage)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("exptected %d yet get %d", http.StatusSeeOther, rr.Code)
	}

	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("did not find in the session")
	}

}
func TestConfig_SubscribeToPlan(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/subscribe?id=1", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "user", data.User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Active:    1,
	})

	handler := http.HandlerFunc(testApp.SubscribeToPlan)
	handler.ServeHTTP(rr, req)
	testApp.Wait.Wait()
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected %d yet %d", http.StatusSeeOther, rr.Code)
	}

}
