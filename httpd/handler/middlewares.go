package handler

import (
	"context"
	"net/http"
	"strconv"
	db "bibliophile-diaries/db/sqlc"
	"bibliophile-diaries/status"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
)

type key int

const (
	IDKey key = iota
	StoreKey
	JwtKey
	NameKey
	UserIDKey
	TokenKey
)

func IDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogIDstr := r.URL.Query().Get("id")
		blogID, err := strconv.Atoi(blogIDstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"please provide a correct id\"}"))
			return
		}

		ctx := r.Context()
		newCtx := context.WithValue(ctx, IDKey, blogID)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func IDMiddlewareStr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogIDstr := r.URL.Query().Get("id")

		ctx := r.Context()
		newCtx := context.WithValue(ctx, IDKey, blogIDstr)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func ProvideStore(store *db.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			newCtx := context.WithValue(ctx, StoreKey, store)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func ProvideJWTAuth(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			newCtx := context.WithValue(ctx, JwtKey, tokenAuth)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, claims, err := jwtauth.FromContext(ctx)
		if err != nil || token == nil || jwt.Validate(token) != nil {
			render.Render(w, r, status.ErrUnauthorized("Incorrect token."))
			return
		}

		name := claims["name"].(string)
		userIDstr := claims["user_id"].(string)

		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			render.Render(w, r, status.ErrUnauthorized("incorrect user id"))
			return
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, NameKey, name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticatePass(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, claims, err := jwtauth.FromContext(ctx)
		var name string
		var userIDstr string
		if err != nil || token == nil || jwt.Validate(token) != nil {
			name = "Guest"
			userIDstr = "0"
		} else {
			name = claims["name"].(string)
			userIDstr = claims["user_id"].(string)
		}

		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			render.Render(w, r, status.ErrUnauthorized("incorrect user id"))
			return
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, NameKey, name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		ctx := r.Context()
		newCtx := context.WithValue(ctx, TokenKey, token)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
