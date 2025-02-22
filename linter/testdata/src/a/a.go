package a

import (
	"context"
	"database/sql"

	"github.com/takaaa220/bunwithoptimistic"
	"github.com/uptrace/bun"
)

type User struct {
	ID   int64
	Name string
}

func badUpdate(ctx context.Context, db *bun.DB, user *User) (sql.Result, error) {
	// want "bun Update query must be wrapped with WithOptimistic"
	return db.NewUpdate().Model(user).Exec(ctx)
}

// goodUpdate demonstrates the correct usage with WithOptimistic
func goodUpdate(ctx context.Context, db *bun.DB, user *User) (sql.Result, error) {
	return bunwithoptimistic.WithOptimistic(db.NewUpdate().Model(user)).Exec(ctx)
}

func f() {
	// The pattern can be written in regular expression.
	var gopher int // want "pattern"
	print(gopher)  // want "identifier is gopher"
}
