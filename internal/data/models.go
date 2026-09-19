package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Animes      AnimeModel
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
}

// returns the model struct containing the initialized model
func NewModels(db *sql.DB) Models {
	return Models{
		Animes:      AnimeModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
