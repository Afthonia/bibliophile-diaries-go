package handler

import (
	db "bibliophile-diaries/db/sqlc"
	"bibliophile-diaries/models"
	"bibliophile-diaries/status"
	"bibliophile-diaries/utils"
	"html/template"
	"net/http"

	"github.com/go-chi/render"
)

type Data struct {
	CommentList []db.ListPostCommentsRow
	Post        db.GetPostRow
}

func CreateComment(w http.ResponseWriter, r *http.Request) {

	commentBind := &models.CommentBind{}

	if err := render.Bind(r, commentBind); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	createdComment, err := store.CreateComment(ctx, db.CreateCommentParams{
		PostID:  commentBind.PostID,
		UserID:  userID,
		Content: commentBind.Content,
	})

	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	commenter := ctx.Value(NameKey).(string)

	createdComment.Commenter = commenter
	createdComment.CommenterID = userID

	commentRow := db.ListPostCommentsRow(createdComment)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &commentRow)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {

	commentBind := &models.CommentBind{}
	if err := render.Bind(r, commentBind); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	ctx := r.Context()
	commentID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	updatedComment, err := store.UpdateComment(ctx, db.UpdateCommentParams{
		ID:      int64(commentID),
		Content: commentBind.Content,
		UserID:  userID,
	})

	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	commenter := ctx.Value(NameKey).(string)
	updatedComment.Commenter = commenter
	updatedComment.CommenterID = userID

	commentRow := db.ListPostCommentsRow(updatedComment)

	render.Render(w, r, &commentRow)
}

func GetPostComments(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	postID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)

	postComments, err := store.ListPostComments(ctx, int64(postID))
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	renderList := utils.Map(postComments, func(e db.ListPostCommentsRow) render.Renderer {

		return &db.ListPostCommentsRow{
			ID:          e.ID,
			CommenterID: e.CommenterID,
			Commenter:   e.Commenter,
			Content:     e.Content,
			CreatedAt:   e.CreatedAt,
		}
	})

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, renderList)
}

func ShowPostComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	postID := ctx.Value(IDKey).(int)

	tmplt, err := template.ParseFiles("templates/post_comments.html")
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	postComments, err := store.ListPostComments(ctx, int64(postID))
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	post, err := store.GetPost(ctx, int64(postID))
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	data := Data{
		CommentList: postComments,
		Post:        post,
	}

	err = tmplt.Execute(w, data)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
	}

	render.Status(r, http.StatusOK)
}

func GetUserComments(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userID := ctx.Value(UserIDKey).(int64)
	store := ctx.Value(StoreKey).(*db.Store)

	userComments, err := store.ListUserComments(ctx, userID)
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	renderList := utils.Map(userComments, func(e db.ListUserCommentsRow) render.Renderer {
		return &db.ListPostCommentsRow{
			ID:          e.ID,
			Commenter:   e.Commenter,
			CommenterID: e.CommenterID,
			Content:     e.Content,
			CreatedAt:   e.CreatedAt,
		}
	})

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, renderList)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	commentID := ctx.Value(IDKey).(int)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	if _, err := store.DeleteComment(ctx, db.DeleteCommentParams{
		UserID: userID,
		ID:     int64(commentID),
	}); err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	render.Render(w, r, status.DelSuccess())
}

// func GetComment(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	commentID := ctx.Value(IDKey).(int)
// 	store := ctx.Value(StoreKey).(*db.Store)

// 	comment, err := store.GetComment(ctx, int64(commentID))

// 	if err != nil {
// 		render.Render(w, r, status.ErrNotFoundOrInternal(err))
// 		return
// 	}

// 	commentRender := db.ListPostCommentsRow(comment)

// 	render.Status(r, http.StatusOK)
// 	render.Render(w, r, &commentRender)
// }
