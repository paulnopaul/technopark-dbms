package usecase

import (
	"DBMSForum/internal/pkg/domain"
	"DBMSForum/internal/pkg/utilities"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type forumUsecase struct {
	DB *sql.DB
}

func NewForumUsecase(db *sql.DB) domain.ForumUsecase {
	return &forumUsecase{
		DB: db,
	}
}

func (u *forumUsecase) Create(f domain.Forum) (*domain.Forum, error) {
	query := "insert into forums(title, username, slug) values ($1, $2, $3) returning title, username, slug, posts, threads;"
	newF := &domain.Forum{}
	err := u.DB.QueryRow(query, f.Title, f.User, f.Slug).Scan(&newF.Title, &newF.User, &newF.Slug, &newF.Posts, &newF.Threads)
	if err != nil {
		return nil, err
	}
	return newF, err
}

func (u *forumUsecase) Details(slug string) (*domain.Forum, error) {
	query := "select title, username, slug, posts, threads from forums where slug = $1"
	f := &domain.Forum{}
	err := u.DB.QueryRow(query, slug).Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (u *forumUsecase) CreateThread(slug string, t domain.Thread) (*domain.Thread, error) {
	createThreadQuery := "insert into threads(title, author, message) values ($1, $2, $3) returning id, title, author,  message, votes, created;"

	newT := &domain.Thread{}
	err := u.DB.QueryRow(createThreadQuery, t.Title, t.Author, t.Message, slug).Scan(&newT.ID, &newT.Title, &newT.Author, &newT.Message, &newT.Votes, &newT.Created)
	if err != nil {
		return nil, err
	}

	createFTQuery := "insert into f_t(f_slug, t_id) values ($1, $2)"
	_, err = u.DB.Exec(createFTQuery, slug, newT.ID)
	if err != nil {
		return nil, err
	}

	newT.Forum = slug
	return newT, nil
}

func generateUserRequest(slug string, params utilities.ArrayOutParams) (string, []interface{}, error) {
	req := psql.Select("nickname", "fullname", "about", "email").From("f_u").Join("users on f_u.u.nick = u.nickname").Where(sq.Eq{"f_u.f_slug": slug})
	if params.Since != "" {
		req = req.Where(sq.Gt{"nickname": params.Since})
	}
	if params.HasLimit {
		req = req.Limit(uint64(params.Limit))
	}
	if params.Desc {
		req = req.OrderBy("nickname desc")
	} else {
		req = req.OrderBy("nickname")
	}
	return req.ToSql()
}

func (u *forumUsecase) Users(slug string, params utilities.ArrayOutParams) ([]domain.User, error) {
	getForumUsersQuery, args, err := generateUserRequest(slug, params)
	if err != nil {
		return nil, errors.New("request creating error")
	}

	rows, err := u.DB.Query(getForumUsersQuery, args...)
	if err != nil {
		return nil, errors.New("database query error")
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

	return resUsers, nil
}

func (u *forumUsecase) Threads(slug string, params utilities.ArrayOutParams) ([]domain.Thread, error) {

}
