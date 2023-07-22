package user_test

import (
	"context"
	"errors"
	"github.com/codymj/go-service/business/data/store/user"
	"github.com/codymj/go-service/business/data/tests"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

var dbc = tests.DbContainer{
	Image: "postgres:alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {
	// create db
	logger, db, teardown := tests.NewUnit(t, dbc)
	t.Cleanup(teardown)

	// create user store
	store := user.NewStore(logger, db)

	// run test
	t.Log("given the need to work with user records")
	{
		testId := 0
		t.Logf("\ttest %d:\twhen handling a single user", testId)
		{
			// mocks for create
			ctx := context.Background()
			now := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
			nu := user.NewUser{
				Name:            "John Smith",
				Email:           "jsmith@example.com",
				Roles:           []string{auth.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

			// invoke create
			usr, err := store.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to create user: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to create user", tests.Success, testId)

			// mocks for query
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:  "service project",
					Subject: usr.Id,
					ExpiresAt: &jwt.NumericDate{
						Time: time.Now().Add(time.Hour).UTC(),
					},
					IssuedAt: &jwt.NumericDate{
						Time: time.Now().UTC(),
					},
				},
				Roles: []string{auth.RoleUser},
			}

			// invoke query (by ID)
			saved, err := store.QueryById(ctx, claims, usr.Id)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to query user by ID: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to query user by ID", tests.Success, testId)

			// assert retrieved user equals the one we saved
			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\ttest %d:\tshould get back same user: %s", tests.Failed, testId, diff)
			}
			t.Logf("\t%s\ttest %d:\tshould get back same user", tests.Success, testId)

			// mocks for update
			uu := user.UpdateUser{
				Name:  tests.StrPtr("Cody Johnson"),
				Email: tests.StrPtr("cjohnson@example.com"),
			}
			claims = auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: "service project",
					ExpiresAt: &jwt.NumericDate{
						Time: time.Now().Add(time.Hour).UTC(),
					},
					IssuedAt: &jwt.NumericDate{
						Time: time.Now().UTC(),
					},
				},
				Roles: []string{auth.RoleAdmin},
			}

			// invoke update
			if err = store.Update(ctx, claims, usr.Id, uu, now); err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to update user: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to update user", tests.Success, testId)

			// invoke query (by email)
			saved, err = store.QueryByEmail(ctx, claims, *uu.Email)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to query user by email: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to query user by email", tests.Success, testId)

			// assert updated fields have been saved
			if *uu.Name != saved.Name {
				t.Fatalf("\t%s\ttest %d:\tshould have updated user.name: %s", tests.Failed, testId, saved.Name)
			}
			if *uu.Email != saved.Email {
				t.Fatalf("\t%s\ttest %d:\tshould have updated user.email: %s", tests.Failed, testId, saved.Email)
			}
			t.Logf("\t%s\ttest %d:\tshould get back same user", tests.Success, testId)

			// invoke delete
			if err = store.Delete(ctx, claims, usr.Id); err != nil {
				t.Fatalf("\t%s\ttest %d:\tshould be able to delete user by ID: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould be able to delete user by ID", tests.Success, testId)

			// invoke query (by ID) which should return an error
			if _, err = store.QueryById(ctx, claims, usr.Id); !errors.Is(err, database.ErrNotFound) {
				t.Fatalf("\t%s\ttest %d:\tshould not be able to query user by ID: %s", tests.Failed, testId, err)
			}
			t.Logf("\t%s\ttest %d:\tshould not be able to query user by ID", tests.Success, testId)
		}
	}
}
