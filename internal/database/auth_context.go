package database

import "context"

type authAccountIDContextKey struct{}

func WithAuthAccountID(ctx context.Context, accountID int64) context.Context {
	return context.WithValue(ctx, authAccountIDContextKey{}, accountID)
}

func AuthAccountIDFromContext(ctx context.Context) (int64, bool) {
	accountID, ok := ctx.Value(authAccountIDContextKey{}).(int64)
	return accountID, ok && accountID > 0
}
