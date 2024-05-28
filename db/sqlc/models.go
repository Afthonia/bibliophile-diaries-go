// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type PermType string

const (
	PermTypeManagePost     PermType = "ManagePost"
	PermTypeCreatePost     PermType = "CreatePost"
	PermTypeManageComments PermType = "ManageComments"
	PermTypeCreateComment  PermType = "CreateComment"
	PermTypeManageUsers    PermType = "ManageUsers"
)

func (e *PermType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PermType(s)
	case string:
		*e = PermType(s)
	default:
		return fmt.Errorf("unsupported scan type for PermType: %T", src)
	}
	return nil
}

type NullPermType struct {
	PermType PermType `json:"PermType"`
	Valid    bool     `json:"valid"` // Valid is true if PermType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPermType) Scan(value interface{}) error {
	if value == nil {
		ns.PermType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PermType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPermType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PermType), nil
}

type Bookshelf struct {
	BookID      string    `json:"book_id"`
	UserID      int64     `json:"user_id"`
	InBookshelf bool      `json:"in_bookshelf"`
	CreatedAt   time.Time `json:"created_at"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Permgroup struct {
	ID          int32      `json:"id"`
	Name        string     `json:"name"`
	Permissions []PermType `json:"permissions"`
}

type Post struct {
	ID        int64          `json:"id"`
	BookTitle string         `json:"book_title"`
	Vote      int16          `json:"vote"`
	UserID    int64          `json:"user_id"`
	Title     sql.NullString `json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type PostLike struct {
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           int64          `json:"id"`
	PermgroupID  int32          `json:"permgroup_id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"password_hash"`
	Bio          sql.NullString `json:"bio"`
	CreatedAt    time.Time      `json:"created_at"`
}
