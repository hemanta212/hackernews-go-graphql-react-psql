package auth

import (
	"context"
	"net/http"

	"github.com/hemanta212/hackernews-go-graphql/internal/users"
	"github.com/hemanta212/hackernews-go-graphql/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			//create user and check if user exists
			user, err := users.GetUserByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)
			// and call the nextwith our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})

	}

}

// finds the user from the context, requires the middleware to have run
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
