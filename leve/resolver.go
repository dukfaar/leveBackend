package leve

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	Model *Model
}

func (r *Resolver) ID(ctx context.Context) (*graphql.ID, error) {
	id := graphql.ID(r.Model.ID.Hex())
	return &id, nil
}

func (r *Resolver) XivdbID(ctx context.Context) (*string, error) {
	return &r.Model.XivdbID, nil
}

func (r *Resolver) Name(ctx context.Context) (*string, error) {
	return &r.Model.Name, nil
}

func (r *Resolver) Class(ctx context.Context) (*string, error) {
	return &r.Model.Class, nil
}

func (r *Resolver) Level(ctx context.Context) (*int32, error) {
	return &r.Model.Level, nil
}

func (r *Resolver) Gil(ctx context.Context) (*int32, error) {
	return &r.Model.Gil, nil
}

func (r *Resolver) Xp(ctx context.Context) (*int32, error) {
	return &r.Model.Xp, nil
}
