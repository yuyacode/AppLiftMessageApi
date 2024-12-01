package request

import (
	"context"
)

type appKindKey struct{}
type userIDKey struct{}

func SetAppKind(ctx context.Context, appKind string) context.Context {
	return context.WithValue(ctx, appKindKey{}, appKind)
}

func GetAppKind(ctx context.Context) (string, bool) {
	appKind, ok := ctx.Value(appKindKey{}).(string)
	return appKind, ok
}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey{}).(int64)
	return userID, ok
}
