package main

import (
	db "bibliophile-diaries/db/sqlc"
	"bibliophile-diaries/httpd/handler"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

func createAPIRoutes(store *db.Store, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	r := chi.NewMux()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Use(handler.ProvideStore(store), handler.ProvideJWTAuth(tokenAuth), jwtauth.Verifier(tokenAuth))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/", func(r chi.Router) {

		r.With(handler.AuthenticatePass).Get("/", handler.ShowPosts)

		r.Route("/user", func(r chi.Router) {
			r.Post("/register", handler.RegisterUser)
			r.Post("/login", handler.LoginUser)

			r.Group(func(r chi.Router) {
				r.Use(handler.Authenticate)
				r.With(handler.IDMiddleware).Delete("/", handler.DeleteUser)
				r.With(handler.TokenMiddleware).Get("/profile", handler.GetUser)
				r.With(handler.IDMiddleware).Get("/dashboard", handler.GetDashboard)
			})
		})

		r.Route("/post", func(r chi.Router) {
			r.With(handler.AuthenticatePass).Get("/list", handler.GetPosts)
			r.With(handler.IDMiddleware).With(handler.AuthenticatePass).Get("/", handler.GetPost)
			r.Group(func(r chi.Router) {
				r.Use(handler.Authenticate)
				r.Post("/", handler.CreatePost)
				r.Get("/liked", handler.GetLikedPosts)
				r.With(handler.IDMiddleware).Patch("/like", handler.TogglePostLike)
				r.With(handler.IDMiddleware).Patch("/", handler.UpdatePost)
				r.With(handler.IDMiddleware).Delete("/", handler.DeletePost)
			})
		})

		r.Route("/comment", func(r chi.Router) {
			r.With(handler.IDMiddleware).Get("/post/list", handler.GetPostComments)
			r.With(handler.IDMiddleware).Get("/", handler.ShowPostComments)
			//r.With(handler.IDMiddleware).Get("/", handler.GetComment)

			r.With(handler.IDMiddleware).Get("/user", handler.GetUserComments)

			r.Group(func(r chi.Router) {
				r.Use(handler.Authenticate)
				r.Post("/", handler.CreateComment)

				r.Group(func(r chi.Router) {
					r.Use(handler.IDMiddleware)
					r.Delete("/", handler.DeleteComment)
					r.Patch("/", handler.UpdateComment)
				})
			})
		})
	})

	return r
}
