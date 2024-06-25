package tarantool

import (
	"context"
	"duttygram/internal"
	"duttygram/pkg/tarantool"
)

type EntityRepository struct {
	client *tarantool.Wrapper
}

func NewEntityRepository(client *tarantool.Wrapper) *EntityRepository {
	return &EntityRepository{client: client}
}

func (e *EntityRepository) Find(ctx context.Context, id string) (*internal.Entity, error) {
	return nil, nil
}

func (e *EntityRepository) Persist(ctx context.Context, entity *internal.Entity) error {
	return nil
}
