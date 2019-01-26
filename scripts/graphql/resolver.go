package graphql

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) Refresh(ctx context.Context, token string) (string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (User, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id string) (User, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context) ([]User, error) {
	panic("not implemented")
}
