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
	return db.NewUpdate().Model(user).Exec(ctx) // want "bun Update query must be wrapped with WithOptimistic"
}

// goodUpdate demonstrates the correct usage with WithOptimistic
func goodUpdate(ctx context.Context, db *bun.DB, user *User) (sql.Result, error) {
	return bunwithoptimistic.WithOptimistic(db.NewUpdate().Model(user)).Exec(ctx)
}
