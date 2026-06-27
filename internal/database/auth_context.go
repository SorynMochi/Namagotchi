package database

import "context"

type authAccountIDContextKey struct{}
type authPlayerIDContextKey struct{}

func WithAuthAccountID(ctx context.Context, accountID int64) context.Context {
	return context.WithValue(ctx, authAccountIDContextKey{}, accountID)
}

func AuthAccountIDFromContext(ctx context.Context) (int64, bool) {
	accountID, ok := ctx.Value(authAccountIDContextKey{}).(int64)
	return accountID, ok && accountID > 0
}

func WithAuthPlayerID(ctx context.Context, playerID int64) context.Context {
	return context.WithValue(ctx, authPlayerIDContextKey{}, playerID)
}

func AuthPlayerIDFromContext(ctx context.Context) (int64, bool) {
	playerID, ok := ctx.Value(authPlayerIDContextKey{}).(int64)
	return playerID, ok && playerID > 0
}
