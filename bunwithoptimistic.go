package bunwithoptimistic

import "github.com/uptrace/bun"

func WithOptimistic(query *bun.UpdateQuery) *WithOptimisticUpdateQuery {
	model, ok := query.GetModel().(WithOptimisticModel)
	if !ok {
		return &WithOptimisticUpdateQuery{UpdateQuery: query}
	}

	query.Where("? = ?", query.FQN(model.VersionColumn()), model.CurrentVersion())
	model.IncrementVersion()

	return &WithOptimisticUpdateQuery{UpdateQuery: query}
}

type WithOptimisticUpdateQuery struct {
	*bun.UpdateQuery
}

type WithOptimisticModel interface {
	CurrentVersion() int
	IncrementVersion() int
	VersionColumn() string
}
