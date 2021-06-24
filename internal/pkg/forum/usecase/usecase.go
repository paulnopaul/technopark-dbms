package usecase

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"technopark-dbms/internal/pkg/domain"
	myerrors "technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/utilities"
)

const (
	createForumQuery     = "insert into forums(title, username, slug) values ($1, $2, $3) returning title, username, slug, posts, threads;"
	getForumDetailsQuery = "select title, username, slug, posts, threads from forums where slug = $1"
	createThreadQuery    = "insert into threads(title, author, message) values ($1, $2, $3) returning id, title, author,  message, votes, created;"
	createFTQuery        = "insert into f_t(f_slug, t_id) values ($1, $2)"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type forumUsecase struct {
	DB     *sql.DB
	UUCase domain.UserUsecase
}

func NewForumUsecase(db *sql.DB) domain.ForumUsecase {
	return &forumUsecase{
		DB: db,
	}
}

func (u *forumUsecase) CreateForum(f domain.Forum) (*domain.Forum, error) {
	query := createForumQuery
	createdForum := &domain.Forum{}
	err := u.DB.QueryRow(query, f.Title, f.User, f.Slug).
		Scan(&createdForum.Title, &createdForum.User, &createdForum.Slug, &createdForum.Posts, &createdForum.Threads)
	if err != nil {
		return nil, err
	}
	return createdForum, err
}

func (u *forumUsecase) Details(slug string) (*domain.Forum, error) {
	query := getForumDetailsQuery
	f := &domain.Forum{}
	err := u.DB.QueryRow(query, slug).Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (u *forumUsecase) CreateThread(slug string, t domain.Thread) (*domain.Thread, error) {
	query := createThreadQuery

	newT := &domain.Thread{}
	err := u.DB.QueryRow(query, t.Title, t.Author, t.Message, slug).Scan(&newT.ID, &newT.Title, &newT.Author, &newT.Message, &newT.Votes, &newT.Created)
	if err != nil {
		return nil, err
	}

	query = createFTQuery
	_, err = u.DB.Exec(query, slug, newT.ID)
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
		return nil, myerrors.QueryCreatingError
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

func generateForumThreadsQuery(slug string, params utilities.ArrayOutParams) (string, []interface{}, error) {
	req := psql.Select("id, title, author, message, forum, votes, slug, created").From("threads").
		Where(sq.Eq{"slug": slug})
	if params.Since != "" {
		req = req.Where(sq.Gt{"created": params.Since})
	}
	if params.HasLimit {
		req = req.Limit(uint64(params.Limit))
	}
	if params.Desc {
		req = req.OrderBy("created desc")
	} else {
		req = req.OrderBy("created")
	}
	return req.ToSql()
}

func (u *forumUsecase) Threads(slug string, params utilities.ArrayOutParams) ([]domain.Thread, error) {
	getThreadsQuery, args, err := generateForumThreadsQuery(slug, params)
	if err != nil {
		return nil, myerrors.QueryCreatingError
	}

	rows, err := u.DB.Query(getThreadsQuery, args...)
	if err != nil {
		return nil, errors.New("database query error")
	}
	defer rows.Close()

	resThreads := make([]domain.Thread, 0)
	for rows.Next() {
		var currentThread domain.Thread
		if err = rows.Scan(); err != nil {
			return nil, errors.New("row scan error")
		}
		resThreads = append(resThreads, currentThread)
	}

	return resThreads, nil
}
