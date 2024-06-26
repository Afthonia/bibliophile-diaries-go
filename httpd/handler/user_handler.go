package handler

import (
	db "bibliophile-diaries/db/sqlc"
	"bibliophile-diaries/models"
	"bibliophile-diaries/status"
	"bibliophile-diaries/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	deneme := &Deneme{}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("{\"error\":\"could not read body\"}"))
// 		return
// 	}

// 	err = json.Unmarshal(data, deneme)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("{\"error\":\"your json is corrupt\"}"))
// 		return
// 	}

// 	resultBytes, err := json.Marshal(deneme)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("{\"error\":\"your json is corrupt\"}"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write(resultBytes)
// }

func LoginUser(w http.ResponseWriter, r *http.Request) {
	userLoginBind := &models.UserLoginBind{}
	if err := render.Bind(r, userLoginBind); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)

	userData, err := store.GetUserByEmail(ctx, userLoginBind.Email)
	if err != nil {
		render.Render(w, r, status.ErrNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(userLoginBind.Password))
	if err != nil {
		render.Render(w, r, status.ErrUnauthorized("wrong password"))
		return
	}

	tokenAuth := ctx.Value(JwtKey).(*jwtauth.JWTAuth)
	expiration := time.Now().Add(90 * 24 * time.Hour)

	claims := map[string]any{
		"user_id": fmt.Sprintf("%d", userData.ID),
		"name":    userData.Name,
		"email":   userData.Email,
	}
	jwtauth.SetExpiry(claims, expiration)

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, models.NewUserPayload(userData, tokenString))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	userRegisterBind := &models.UserRegisterBind{}
	if err := render.Bind(r, userRegisterBind); err != nil {
		render.Render(w, r, status.ErrBadRequest(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegisterBind.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)

	userData, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:         strings.TrimSpace(userRegisterBind.Name),
		Email:        strings.TrimSpace(userRegisterBind.Email),
		PermgroupID:  1,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	///////// TOKEN CREATION
	tokenAuth := ctx.Value(JwtKey).(*jwtauth.JWTAuth)
	expiration := time.Now().Add(90 * 24 * time.Hour)

	claims := map[string]any{
		"user_id": fmt.Sprintf("%d", userData.ID),
		"name":    userData.Name,
		"email":   userData.Email,
	}
	jwtauth.SetExpiry(claims, expiration)

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}
	///////// TOKEN CREATION

	render.Status(r, http.StatusCreated)
	render.Render(w, r, models.NewUserPayload(userData, tokenString))
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	tokenStr := ctx.Value(TokenKey).(string)

	tokenAuth := ctx.Value(JwtKey).(*jwtauth.JWTAuth)

	token, err := tokenAuth.Decode(tokenStr)
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	user := token.PrivateClaims()

	userIDStr := user["user_id"].(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.Write([]byte("\"error\": \"user id is invalid\""))
		return
	}

	userProf, err := store.GetUser(ctx, int64(userID))
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, models.UserProfilePayload(userProf))

	//w.Write([]byte(fmt.Sprintf("{\"user\": \"%v\"}", userEmail)))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(IDKey).(int)

	if err := store.DeleteUser(ctx, int64(userID)); err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	render.Render(w, r, status.DelSuccess())
}

func UpdateBio(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)
	bio := ctx.Value(IDKey).(string)

	updatedBio, err := store.UpdateBio(ctx, db.UpdateBioParams{
		ID:  userID,
		Bio: sql.NullString{String: bio, Valid: bio != ""},
	})

	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strconv.FormatBool(updatedBio.Valid)))
}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(IDKey).(int)
	log.Println(userID)

	dashboard, err := store.GetDashboard(ctx, sql.NullInt64{Int64: int64(userID), Valid: userID != 0})
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, models.DashboardPayload(dashboard))
}

func ToggleBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bookID := ctx.Value(IDKey).(string)
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(UserIDKey).(int64)

	inBookshelf, err := store.ToggleBook(ctx, db.ToggleBookParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		render.Render(w, r, status.ErrNotFoundOrInternal(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(inBookshelf)))
}

func GetBookshelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store := ctx.Value(StoreKey).(*db.Store)
	userID := ctx.Value(IDKey).(int)
	log.Println(userID)

	books, err := store.ListBookshelf(ctx, int64(userID))
	if err != nil {
		render.Render(w, r, status.ErrInternal(err))
		return
	}

	bookRows := []db.ListBookshelfRow(books)

	renderList := utils.Map(bookRows, func(e db.ListBookshelfRow) render.Renderer {
		return &db.ListBookshelfRow{
			UserID:      e.UserID,
			BookID:      e.BookID,
			InBookshelf: e.InBookshelf,
			CreatedAt:   e.CreatedAt,
		}
	})

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, renderList)
}
