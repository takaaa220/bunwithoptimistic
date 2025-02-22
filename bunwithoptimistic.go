package bunwithoptimistic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
)

func WithOptimistic(query *bun.UpdateQuery) *WithOptimisticUpdateQuery {
	model, ok := query.GetModel().Value().(WithOptimisticModel)
	if !ok {
		return &WithOptimisticUpdateQuery{UpdateQuery: query}
	}

	query.Where("? = ?", query.FQN(model.VersionColumn()), model.CurrentVersion())
	model.IncrementVersion()

	return &WithOptimisticUpdateQuery{UpdateQuery: query}
}

func (q *WithOptimisticUpdateQuery) Exec(ctx context.Context) (sql.Result, error) {
	res, err := q.UpdateQuery.Exec(ctx)
	if err != nil {
		return res, err
	}

	if model, ok := q.GetModel().(WithOptimisticModel); ok {
		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, &OptimisticConflictError{
				Version:   model.CurrentVersion() - 1,
				TableName: q.GetTableName(),
			}
		}
	}

	// TODO: support multiple models

	return res, nil
}

type OptimisticConflictError struct {
	Version   int
	TableName string
}

func (e *OptimisticConflictError) Error() string {
	return fmt.Sprintf("bun: version %d is out of version on table %s", e.Version, e.TableName)
}

type WithOptimisticUpdateQuery struct {
	*bun.UpdateQuery
}

type WithOptimisticModel interface {
	CurrentVersion() int
	IncrementVersion() int
	VersionColumn() string
}
