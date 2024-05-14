package models

import (
	"errors"
	"net/http"
	"net/mail"
	"time"
	db "bibliophile-diaries/db/sqlc"
)

type UserRegisterBind struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `name:"password"`
}

type UserLoginBind struct {
	Email    string `json:"email"`
	Password string `name:"password"`
}

type UserPayload struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

type UserProfPayload struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type BookModel struct {
	Title string `json:"title"`
	Author string `json:"author"`
	CoverURL string `json:"cover_url"`
}

type Dashboard struct {
	UserPosts    int64 `json:"user_posts"`
	UserComments int64 `json:"user_comments"`
	Bookshelf    []BookModel `json:"bookshelf"`
}

func NewUserPayload(userData db.User, token string) *UserPayload {
	return &UserPayload{
		Name:      userData.Name,
		Email:     userData.Email,
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		Token:     token,
	}
}

func (p *UserRegisterBind) Bind(r *http.Request) error {
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("your email must be a valid address")
	}

	// if _, err := verifyPassword(p.Password); err != nil {
	// 	return err
	// }

	if len(p.Password) < 6 {
		return errors.New("password must be atleast 6 characters")
	}

	if len(p.Name) == 0 {
		return errors.New("your account must have a name")
	}

	return nil
}

// func verifyPassword(pswrd string) (bool, error) {
// 	var (
// 		upp, low, num, sym bool
// 		tot                uint8
// 	)

// 	err := errors.New("your password must be valid")

// 	for _, char := range pswrd {
// 		switch {
// 		case unicode.IsUpper(char):
// 			upp = true
// 			tot++
// 		case unicode.IsLower(char):
// 			low = true
// 			tot++
// 		case unicode.IsNumber(char):
// 			num = true
// 			tot++
// 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// 			sym = true
// 			tot++
// 		default:
// 			return false, err
// 		}
// 	}

// 	if !upp || !low || !num || !sym || tot < 8 {
// 		return false, err
// 	}

// 	return true, nil
// }

func (p *UserLoginBind) Bind(r *http.Request) error {
	if len(p.Password) < 6 {
		return errors.New("password must be atleast 6 characters")
	}

	return nil
}

func UserProfilePayload(userData db.User) *UserProfPayload {
	return &UserProfPayload{
		Name:      userData.Name,
		Email:     userData.Email,
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
	}
}

func DashboardPayload(dashboard db.GetDashboardRow) *Dashboard {
	bookshelf := []BookModel{}
	
	return &Dashboard{
		UserPosts:    dashboard.UserPosts,
		UserComments: dashboard.UserComments,
		Bookshelf:    bookshelf,
	}
}

func (u *UserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UserProfPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (gd *Dashboard) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
