package handler

import (
	db "bibliophile-diaries/db/sqlc"
	"bibliophile-diaries/models"
	"bibliophile-diaries/status"
	"bibliophile-diaries/utils"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	postBind := &models.PostBind{}
	// if err := postBind.FormBind(r); err != nil {
	// 	render.Render(w, r, status.ErrBadRequest(err))
	// 	return
	// }

	if err := render.Bind(r, postBind); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	createdPost, err := store.CreatePost(ctx, db.CreatePostParams{
		UserID:    userID,
		Title:     sql.NullString{String: postBind.Title, Valid: postBind.Title != ""},
		Content:   postBind.Content,
		BookTitle: postBind.BookTitle,
		Vote:      int16(postBind.Vote),
	})

	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	userName := ctx.Value(NameKey).(string)

	createdPost.Author = userName
	createdPost.AuthorID = userID

	postRow := db.ListPostsRow(createdPost)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &postRow)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	postBind := &models.PostBind{}
	if err := postBind.FormBind(r); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	ctx := r.Context()
	postID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	updatedPost, err := store.UpdatePost(ctx, db.UpdatePostParams{
		UserID:    userID,
		ID:        int64(postID),
		Title:     sql.NullString{String: postBind.Title, Valid: postBind.Title != ""},
		Content:   postBind.Content,
		BookTitle: postBind.BookTitle,
		Vote:      int16(postBind.Vote),
	})
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	username := ctx.Value(NameKey).(string)
	updatedPost.Author = username
	updatedPost.AuthorID = userID

	postRow := db.ListPostsRow(updatedPost)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &postRow)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)
	log.Println(userID)

	posts, err := store.ListPosts(ctx, userID)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	postRows := []db.ListPostsRow(posts)

	renderList := utils.Map(postRows, func(e db.ListPostsRow) render.Renderer {
		return &db.ListPostsRow{
			ID:        e.ID,
			IsLiked:   e.IsLiked,
			LikeCount: e.LikeCount,
			Author:    e.Author,
			AuthorID:  e.AuthorID,
			Title:     e.Title,
			Content:   e.Content,
			CreatedAt: e.CreatedAt,
			BookTitle: e.BookTitle,
			Vote:      e.Vote,
		}
	})

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, renderList)
}

func GetLikedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)
	//log.Println(userID)

	var posts []db.GetLikedPostsRow

	posts, err := store.GetLikedPosts(ctx, userID)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	postRows := []db.GetLikedPostsRow(posts)

	renderList := utils.Map(postRows, func(e db.GetLikedPostsRow) render.Renderer {
		return &db.ListPostsRow{
			ID:        e.ID,
			IsLiked:   e.IsLiked,
			LikeCount: e.LikeCount,
			Author:    e.Author,
			AuthorID:  e.AuthorID,
			Title:     e.Title,
			Content:   e.Content,
			CreatedAt: e.CreatedAt,
			BookTitle: e.BookTitle,
			Vote:      e.Vote,
		}
	})

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, renderList)
}

func ShowPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmplt, _ := template.ParseFiles("templates/index.html")
	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	posts, err := store.ListPosts(ctx, userID)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	err = tmplt.Execute(w, posts)

	if err != nil {
		return
	}

	render.Status(r, http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	if _, err := store.DeletePost(ctx, db.DeletePostParams{
		UserID: userID,
		ID:     int64(postID),
	}); err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	render.Render(w, r, status.DelSuccess())
}

func TogglePostLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	isLiked, err := store.TogglePostLike(ctx, db.TogglePostLikeParams{
		PostID: int64(postID),
		UserID: userID,
	})
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(isLiked)))
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)

	post, err := store.GetPost(ctx, int64(postID))
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	postRender := &db.GetPostRow{
		Title:     post.Title,
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Author:    post.Author,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		IsLiked:   post.IsLiked,
		LikeCount: post.LikeCount,
		BookTitle: post.BookTitle,
		Vote:      post.Vote,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.Render(w, r, postRender)
}
