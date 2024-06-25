package internal

import (
	"context"
	"fmt"
)

type Messenger interface {
	SendText(ctx context.Context, content string) error
}

type EntityRepository interface {
	Find(ctx context.Context, id string) (*Entity, error)
	Persist(ctx context.Context, entity *Entity) error
}

type Entity struct {
	Id    string
	Stage Stage
}

func NewEntity(id string, stage Stage) *Entity {
	return &Entity{Id: id, Stage: stage}
}

type StateFactory struct {
	repo     EntityRepository
	scenario *Scenario
}

func (f *StateFactory) Create(ctx context.Context, id string, messenger Messenger) (*State, error) {
	entity, err := f.repo.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't execute finding state process: %w", err)
	}
	if entity == nil {
		entity = NewEntity(id, f.scenario.GetStartStage())
		err = f.repo.Persist(ctx, entity)
		if err != nil {
			return nil, fmt.Errorf("can't execute persisting state process: %w", err)
		}
	}

	return NewState(entity, f.scenario, messenger, f.repo), nil
}

func NewStateFactory(repo EntityRepository, scenario *Scenario) *StateFactory {
	return &StateFactory{repo: repo, scenario: scenario}
}

type State struct {
	entity    *Entity
	scenario  *Scenario
	messenger Messenger
	repo      EntityRepository
}

func NewState(entity *Entity, scenario *Scenario, messenger Messenger, repo EntityRepository) *State {
	return &State{entity: entity, scenario: scenario, messenger: messenger, repo: repo}
}
