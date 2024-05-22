package errs

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//nolint:stylecheck
var (
	IncorrectJSON  = errors.New("incorrect json")
	Duplicate      = errors.New("some of input data already exists")
	NilID          = errors.New("nil ID")
	WrongID        = errors.New("wrong id")
	IncorrectToken = errors.New("incorrect token")
	Auth           = errors.New("auth error")
	InvalidToken   = errors.New("invalid token")
	NoReference    = errors.New("no such tag/feature")
	Unknown        = errors.New("unknown error")
	Validation     = errors.New("incorrect input data")

	NotAdminErr    = errors.New("you are not admin")
	NoUserSuchUser = errors.New("no such user")
	NotLogged      = errors.New("you are not logged")
	WrongRole      = errors.New("wrong role")
	NoRowsInResult = errors.New("no such data")
)

func Handle(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return NoRowsInResult

	case errors.Is(err, jwt.ErrSignatureInvalid), errors.Is(err, jwt.ErrTokenMalformed):
		return Auth

	}
	if pgErr, ok := err.(*pgconn.PgError); ok { //nolint:errorlint
		switch pgErr.Code {
		case "23505":
			return Duplicate
		case "23503":
			return NoReference
		default:
			return Unknown
		}
	}

	if _, ok := err.(validator.ValidationErrors); ok { //nolint:errorlint
		return Validation
	}
	var jsErr *json.InvalidUnmarshalError
	if errors.As(err, &jsErr) { //nolint: errorlint
		return Validation
	}
	if _, ok := err.(schema.EmptyFieldError); ok { //nolint:errorlint
		return Validation
	}
	if _, ok := err.(schema.MultiError); ok { //nolint:errorlint
		return Validation
	}
	if _, ok := err.(redis.Error); ok { //nolint:errorlint
		return Unknown
	}

	return err
}
