package usecase

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/forum"
	"technopark-dbms/internal/pkg/utilities"
)

const (
	createForumQuery     = "insert into forums(title, username, slug) values ($1, $2, $3) returning title, username, slug, posts, threads;"
	forumExistsQuery     = "select slug from forums where slug = $1;"
	getForumDetailsQuery = "select title, username, slug, posts, threads from forums where slug = $1;"
	createFTQuery        = "insert into f_t(f_slug, t_id) values ($1, $2);"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type forumUsecase struct {
	DB     *sql.DB
	UUCase domain.UserUsecase
}

func (u *forumUsecase) Exists(slug string) (bool, error) {
	var foundSlug string
	err := u.DB.QueryRow(forumExistsQuery, slug).Scan(&foundSlug)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewForumUsecase(db *sql.DB, userUsecase domain.UserUsecase) domain.ForumUsecase {
	return &forumUsecase{
		DB:     db,
		UUCase: userUsecase,
	}
}

func (u *forumUsecase) CreateForum(f domain.Forum) (*domain.Forum, error) {
	authorExists, err := u.UUCase.Exists(f.User, "")
	if err != nil {
		return nil, err
	} else if !authorExists {
		return nil, forum.AuthorNotExists
	}

	foundForum, err := u.Details(f.Slug)
	if err == nil {
		return foundForum, forum.AlreadyExists
	} else if err != forum.NotFound {
		return nil, err
	}

	createdForum := &domain.Forum{}
	err = u.DB.QueryRow(createForumQuery, f.Title, f.User, f.Slug).
		Scan(&createdForum.Title, &createdForum.User, &createdForum.Slug, &createdForum.Posts, &createdForum.Threads)
	if err != nil {
		return nil, err
	}
	return createdForum, err
}

func (u *forumUsecase) Details(slug string) (*domain.Forum, error) {
	f := &domain.Forum{}
	err := u.DB.QueryRow(getForumDetailsQuery, slug).Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
	if err == sql.ErrNoRows {
		return nil, forum.NotFound
	} else if err != nil {
		return nil, err
	}
	return f, nil
}

func generateCreateThreadQuery(forumSlug string, t domain.Thread) (string, []interface{}, error) {
	req := psql.Insert("threads").Columns("author", "forum", "message", "title")
	values := []interface{}{t.Author, forumSlug, t.Message, t.Title}
	if t.Created != "" {
		req = req.Columns("created")
		values = append(values, t.Created)
	}
	if t.Slug != "" {
		req = req.Columns("slug")
		values = append(values, t.Slug)
	}
	req = req.Values(values...)
	req = req.Suffix("returning id, author, forum, message, title, created, slug")
	return req.ToSql()
}

func (u *forumUsecase) CreateThread(forumSlug string, t domain.Thread) (*domain.Thread, error) {
	createThreadQuery, args, err := generateCreateThreadQuery(forumSlug, t)
	if err != nil {
		return nil, err
	}
	newThread := &domain.Thread{}
	created := sql.NullString{}
	slug := sql.NullString{}
	err = u.DB.QueryRow(createThreadQuery, args...).
		Scan(&newThread.ID, &newThread.Author, &newThread.Forum, &newThread.Message, &newThread.Title, &created, &slug)
	if err != nil {
		return nil, err
	}

	if created.Valid {
		newThread.Created = created.String
	}
	if slug.Valid {
		newThread.Slug = slug.String
	}
	return newThread, nil
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
	query, args, err := generateUserRequest(slug, params)
	if err != nil {
		return nil, err
	}

	rows, err := u.DB.Query(query, args...)
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
			return nil, err
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

func (u *forumUsecase) Threads(forumSlug string, params utilities.ArrayOutParams) ([]domain.Thread, error) {

	forumExists, err := u.Exists(forumSlug)
	if err != nil {
		return nil, err
	} else if !forumExists {
		return nil, forum.NotFound
	}

	getThreadsQuery, args, err := generateForumThreadsQuery(forumSlug, params)
	if err != nil {
		return nil, err
	}

	rows, err := u.DB.Query(getThreadsQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resThreads := make([]domain.Thread, 0)
	for rows.Next() {
		var currentThread domain.Thread
		if err = rows.Scan(); err != nil {
			return nil, err
		}
		resThreads = append(resThreads, currentThread)
	}

	return resThreads, nil
}
