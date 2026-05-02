package gateway_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	_ "modernc.org/sqlite"

	"github.com/otakakot/sample-go-server-db-test/internal/gateway"
)

func TestGateway(t *testing.T) {
	t.Parallel()

	file := uuid.NewString()

	db, err := sql.Open("sqlite", "file:"+file+"?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	if db.Ping() != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}

		if err := os.Remove(file); err != nil {
			t.Error(err)
		}
	})

	ddl := `
	CREATE TABLE users (
		id   TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`

	if _, err := db.Exec(ddl); err != nil {
		t.Fatal(err)
	}

	gw := gateway.New(db)

	ctx := context.Background()

	name1 := uuid.NewString()

	t.Log("create user")
	created, err := gw.CreateUser(ctx, gateway.CreateUserDAI{
		Name: name1,
	})
	if err != nil {
		t.Fatal(err)
	}

	userID := created.User.ID

	t.Log("read user")
	read1, err := gw.ReadUser(ctx, gateway.ReadUserDAI{
		ID: userID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(created.User, read1.User); diff != "" {
		t.Fatalf("created: %v, read: %v", created.User, read1.User)
	}

	name2 := uuid.NewString()

	t.Log("update user")
	if _, err := gw.UpdateUser(ctx, gateway.UpdateUserDAI{
		ID:   userID,
		Name: name2,
	}); err != nil {
		t.Fatal(err)
	}

	read2, err := gw.ReadUser(ctx, gateway.ReadUserDAI{
		ID: userID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(name2, read2.User.Name); diff != "" {
		t.Fatalf("name2: %v, read2: %v", name2, read2.User.Name)
	}

	t.Log("delete user")
	if _, err := gw.DeleteUser(ctx, gateway.DeleteUserDAI{
		ID: userID,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := gw.ReadUser(ctx, gateway.ReadUserDAI{
		ID: userID,
	}); err == nil {
		t.Fatal("user should not exist")
	}
}
