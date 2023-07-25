package usergroup_test

import (
	"encoding/json"
	"github.com/codymj/go-service/app/services/api/handlers"
	"github.com/codymj/go-service/business/data/tests"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// UserTests holds methods for each user subtest.
type UserTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

// TestUsers is the entry point for testing user management functions.
func TestUsers(t *testing.T) {
	// init test
	test := tests.NewIntegration(
		t,
		tests.DbContainer{
			Image: "postgres:alpine",
			Port:  "5432",
			Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
		},
	)
	t.Cleanup(test.Teardown)

	// init app
	shutdown := make(chan os.Signal, 1)
	ut := UserTests{
		app: handlers.ApiMux(handlers.ApiMuxConfig{
			Shutdown: shutdown,
			Logger:   test.Logger,
			Auth:     test.Auth,
			DB:       test.Db,
		}),
		userToken:  test.Token("user@example.com", "gophers"),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	// run tests
	t.Run("getToken200", ut.getToken200)
	// todo: add more
}

// getToken200 is a subtest for testing GET /users/token, returning a 200.
func (ut *UserTests) getToken200(t *testing.T) {
	// init request
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	r.SetBasicAuth("admin@example.com", "gophers")
	w := httptest.NewRecorder()

	// do call
	ut.app.ServeHTTP(w, r)

	// run test
	t.Log("given the need to issue tokens to known users")
	{
		testId := 0
		t.Logf("\ttest %d:\twhen fetching a token with valid credentials", testId)
		{
			// assert response
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\ttest %d:\tshould receive 200 OK", tests.Failed, testId)
			}
			t.Logf("\t%s\ttest %d:\tshould recieve 200 OK", tests.Success, testId)

			// assert token
			var got struct {
				Token string `json:"token"`
			}
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to unmarshal response", tests.Failed, testId)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to unmarshal response", tests.Success, testId)
		}
	}
}
