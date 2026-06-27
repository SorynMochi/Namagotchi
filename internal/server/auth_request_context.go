package server

import (
	"context"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type authAccountRequestContextKey struct{}

func withAuthAccountRequestContext(ctx context.Context, account database.AuthAccount) context.Context {
	if account.ID < 1 {
		return ctx
	}

	return context.WithValue(ctx, authAccountRequestContextKey{}, account)
}

func authAccountFromRequestContext(ctx context.Context) (database.AuthAccount, bool) {
	account, ok := ctx.Value(authAccountRequestContextKey{}).(database.AuthAccount)
	return account, ok && account.ID > 0
}
