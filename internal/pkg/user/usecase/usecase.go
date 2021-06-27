package usecase

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
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
	DB *pgx.ConnPool
}

func (u *userUsecase) GetProfiles(nickname, email string) (domain.UserArray, error) {
	query := getUsersDetailsQuery

	rows, err := u.DB.Query(query, nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resUsers := make(domain.UserArray, 0)
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

func NewUserUsecase(db *pgx.ConnPool) domain.UserUsecase {
	return &userUsecase{
		DB: db,
	}
}

func (u *userUsecase) CreateUser(nickname string, createData domain.User) (*domain.User, error, domain.UserArray) {
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
		if err == pgx.ErrNoRows {
			return nil, user.NotExistsError
		}
		return nil, err
	}
	return foundUser, err
}
func (u *userUsecase) UpdateUser(nickname string, userUpdate domain.User) (*domain.User, error) {
	if userUpdate.Email == "" && userUpdate.About == "" && userUpdate.Fullname == "" {
		return u.GetProfile(nickname)
	}

	foundUser, err := u.GetProfile(nickname)
	if err != nil {
		return nil, err
	}

	if userUpdate.Email != "" {
		emailConflict, err := u.UserExists("", userUpdate.Email)
		if err != nil {
			return nil, err
		}
		if emailConflict {
			return nil, user.UpdateConflict
		}
	}

	if userUpdate.Email == "" {
		userUpdate.Email = foundUser.Email
	} else {
		foundUser.Email = userUpdate.Email
	}
	if userUpdate.About == "" {
		userUpdate.About = foundUser.About
	} else {
		foundUser.About = userUpdate.About
	}
	if userUpdate.Fullname == "" {
		userUpdate.Fullname = foundUser.Fullname
	} else {
		foundUser.Fullname = userUpdate.Fullname
	}

	query := "update users set fullname = $1, about = $2, email = $3 where nickname = $4"
	_, err = u.DB.Exec(query, userUpdate.Fullname, userUpdate.About, userUpdate.Email, nickname)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (u *userUsecase) UserExists(nickname string, email string) (bool, error) {
	var foundNick string
	err := u.DB.QueryRow(checkUserExistsQuery, nickname, email).
		Scan(&foundNick)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
