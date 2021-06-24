package usecase

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/user"
)

const (
	createUserQuery      = "insert into users(nickname, fullname, about, email) values ($1, $2, $3, $4) returning nickname, fullname, about, email;"
	getUserDetailsQuery  = "select nickname, fullname, about, email from users where nickname = $1;"
	getUsersDetailsQuery = "select nickname, fullname, about, email from users where nickname = $1 or email = $2;"
	checkUserExistsQuery = "select nickname from users where nickname = $1 or email = $2;"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type userUsecase struct {
	DB *sql.DB
}

func (u *userUsecase) GetProfiles(nickname, email string) ([]domain.User, error) {
	query := getUsersDetailsQuery

	rows, err := u.DB.Query(query, nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resUsers := make([]domain.User, 0)
	for rows.Next() {
		var currentUser domain.User
		if err = rows.Scan(&currentUser.Nickname,
			&currentUser.Fullname,
			&currentUser.About,
			&currentUser.Email); err != nil {
			return nil, errors.New("row scan error")
		}
		resUsers = append(resUsers, currentUser)
	}
	if len(resUsers) == 0 {
		return nil, user.NotExistsError
	}
	return resUsers, nil
}

func NewUserUsecase(db *sql.DB) domain.UserUsecase {
	return &userUsecase{
		DB: db,
	}
}

func (u *userUsecase) CreateUser(nickname string, createData domain.User) (*domain.User, error, []domain.User) {
	checkedProfiles, err := u.GetProfiles(nickname, createData.Email)
	if err == nil {
		return nil, user.AlreadyExistsError, checkedProfiles
	} else if err != user.NotExistsError {
		return nil, err, nil
	}

	query := createUserQuery
	createdUser := &domain.User{}
	err = u.DB.QueryRow(query, nickname, createData.Fullname, createData.About, createData.Email).
		Scan(&createdUser.Nickname, &createdUser.Fullname, &createdUser.About, &createdUser.Email)
	if err != nil {
		return nil, err, nil
	}
	return createdUser, nil, nil
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
		req = req.Set("about", profileUpdate.About)
	}
	if profileUpdate.Email != "" {
		req = req.Set("email", profileUpdate.Email)
	}
	if profileUpdate.Fullname != "" {
		req = req.Set("fullname", profileUpdate.Fullname)
	}
	return req.ToSql()
}

func (u *userUsecase) UpdateProfile(nickname string, profileUpdate domain.User) (*domain.User, error) {
	if profileUpdate.Email == "" && profileUpdate.About == "" && profileUpdate.Fullname == "" {
		return u.GetProfile(nickname)
	}

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
		return nil, err
	}
	updatedUser := &domain.User{}
	err = u.DB.QueryRow(query, args...).
		Scan(&updatedUser.Nickname, &updatedUser.Fullname, &updatedUser.About, &updatedUser.Email)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (u *userUsecase) Exists(nickname string, email string) (bool, error) {
	var foundNick string
	err := u.DB.QueryRow(checkUserExistsQuery, nickname, email).
		Scan(&foundNick)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
