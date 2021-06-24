package usecase

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"technopark-dbms/internal/pkg/domain"
	myerrors "technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/user"
)

const (
	createUserQuery      = "insert into users(nickname, fullname, about, email) values ($1, $2, $3, $4) returning nickname, fullname, about, email;"
	getUserDetailsQuery  = "select nickname, fullname, about, email from users where nickname = $1"
	checkUserExistsQuery = "select nickname from users where nickname = $1 or email = $2;"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type userUsecase struct {
	DB *sql.DB
}

func NewUserUsecase(db *sql.DB) domain.UserUsecase {
	return &userUsecase{
		DB: db,
	}
}

func (u *userUsecase) CreateUser(nickname string, createData domain.User) (*domain.User, error) {
	checkedUser, err := u.GetProfile(nickname)
	if err == nil {
		return checkedUser, user.AlreadyExistsError
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	query := createUserQuery
	createdUser := &domain.User{}
	err = u.DB.QueryRow(query, nickname, createData.Fullname, createData.About, createData.Email).
		Scan(&createdUser.Nickname, &createdUser.Fullname, &createdUser.About, &createdUser.Email)
	if err != nil {
		return nil, err
	}
	return createdUser, err
}

func (u *userUsecase) GetProfile(nickname string) (*domain.User, error) {
	query := getUserDetailsQuery
	foundUser := &domain.User{}
	err := u.DB.QueryRow(query, nickname).
		Scan(&foundUser.Nickname, &foundUser.Fullname, &foundUser.About, &foundUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.NotExistsError
		}
		return nil, err
	}
	return foundUser, err
}

func generateUpdateProfileRequest(nickname string, profileUpdate domain.User) (string, []interface{}, error) {
	req := psql.Update("users").Where(sq.Eq{"nickname": nickname}).Suffix("returning nickname, fullname, about, email")
	if profileUpdate.About != "" {
		req.Set("about", profileUpdate.About)
	}
	if profileUpdate.Email != "" {
		req.Set("email", profileUpdate.Email)
	}
	if profileUpdate.Fullname != "" {
		req.Set("fullname", profileUpdate.Fullname)
	}
	return req.ToSql()
}

func (u *userUsecase) UpdateProfile(nickname string, profileUpdate domain.User) (*domain.User, error) {
	userExists, err := u.Exists(nickname, "")
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, user.NotExistsError
	}

	if profileUpdate.Email != "" {
		emailConflict, err := u.Exists("", profileUpdate.Email)
		if err != nil {
			return nil, err
		}
		if emailConflict {
			return nil, user.UpdateConflict
		}
	}

	query, args, err := generateUpdateProfileRequest(nickname, profileUpdate)
	if err != nil {
		return nil, myerrors.QueryCreatingError
	}
	updatedUser := &domain.User{}
	err = u.DB.QueryRow(query, args).
		Scan(&updatedUser.Nickname, &updatedUser.Fullname, &updatedUser.About, &updatedUser.Email)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (u *userUsecase) Exists(nickname string, email string) (bool, error) {
	query := checkUserExistsQuery
	var foundNick string
	err := u.DB.QueryRow(query, nickname, email).
		Scan(&foundNick)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
