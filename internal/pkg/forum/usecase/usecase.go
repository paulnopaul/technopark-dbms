package usecase

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"technopark-dbms/internal/pkg/constants"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/forum"
	"technopark-dbms/internal/pkg/thread"
	"technopark-dbms/internal/pkg/utilities"
	"time"
)

const (
	createForumQuery     = "insert into forums(title, username, slug) values ($1, (select nickname from users u where u.nickname = $2), $3) returning title, username, slug, posts, threads;"
	forumExistsQuery     = "select slug from forums where slug = $1;"
	getForumDetailsQuery = "select title, username, slug, posts, threads from forums where slug = $1;"
	createFTQuery        = "insert into f_t(f_slug, t_id) values ($1, $2);"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type forumUsecase struct {
	DB     *pgx.ConnPool
	UUCase domain.UserUsecase
	TUCase domain.ThreadUsecase
}

func (u *forumUsecase) Exists(slug string) (bool, error) {
	var foundSlug string
	err := u.DB.QueryRow(forumExistsQuery, slug).Scan(&foundSlug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewForumUsecase(db *pgx.ConnPool, userUsecase domain.UserUsecase, threadUsecase domain.ThreadUsecase) domain.ForumUsecase {
	return &forumUsecase{
		DB:     db,
		UUCase: userUsecase,
		TUCase: threadUsecase,
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
	if err == pgx.ErrNoRows {
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
	req = req.Suffix("returning id, author, (select f.slug from forums f where f.slug = $2), message, title, created, slug")
	return req.ToSql()
}

func (u *forumUsecase) CreateThread(forumSlug string, t domain.Thread) (*domain.Thread, error) {
	if t.Slug != "" {
		foundThread, err := u.TUCase.GetThreadDetails(utilities.NewSlugOrId(t.Slug))
		if err != thread.NotFound {
			if err == nil {
				return foundThread, forum.AlreadyExists
			}
			return nil, err
		}
	}
	authorExists, err := u.UUCase.Exists(t.Author, "")
	if err != nil {
		return nil, err
	} else if !authorExists {
		return nil, forum.AuthorNotExists
	}

	forumExists, err := u.Exists(forumSlug)
	if err != nil {
		return nil, err
	} else if !forumExists {
		return nil, forum.AuthorNotExists
	}

	createThreadQuery, args, err := generateCreateThreadQuery(forumSlug, t)
	if err != nil {
		return nil, err
	}
	newThread := &domain.Thread{}
	var slug *string
	var created *time.Time
	err = u.DB.QueryRow(createThreadQuery, args...).
		Scan(&newThread.ID, &newThread.Author, &newThread.Forum, &newThread.Message, &newThread.Title, &created, &slug)
	if err != nil {
		return nil, err
	}
	if created != nil {
		newThread.Created = created.Format(constants.TimeLayout)
	}
	if slug != nil {
		newThread.Slug = *slug
	}
	return newThread, nil
}

func generateUserRequest(slug string, params utilities.ArrayOutParams) (string, []interface{}, error) {
	req := psql.Select("u.nickname", "u.fullname", "u.about", "u.email").From("f_u").Join("users u on f_u.u = u.nickname").Where(sq.Eq{"f_u.f": slug})
	if params.Desc {
		if params.Since != "" {
			req = req.Where(sq.Lt{"nickname": params.Since})
		}
		req = req.OrderBy("nickname desc")
	} else {
		if params.Since != "" {
			req = req.Where(sq.Gt{"nickname": params.Since})
		}
		req = req.OrderBy("nickname")
	}
	req = req.Limit(uint64(params.Limit))
	return req.ToSql()
}

func (u *forumUsecase) Users(forumSlug string, params utilities.ArrayOutParams) ([]domain.User, error) {
	forumExists, err := u.Exists(forumSlug)
	if err != nil {
		return nil, err
	} else if !forumExists {
		return nil, forum.NotFound
	}

	query, args, err := generateUserRequest(forumSlug, params)
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

func generateForumThreadsQuery(forum string, params utilities.ArrayOutParams) (string, []interface{}, error) {
	req := psql.Select("id, title, author, forum, message, slug, created, votes").From("threads").
		Where(sq.Eq{"forum": forum})
	if params.Desc {
		if params.Since != "" {
			req = req.Where(sq.LtOrEq{"created": params.Since})
		}
		req = req.OrderBy("created desc")
	} else {
		if params.Since != "" {
			req = req.Where(sq.GtOrEq{"created": params.Since})
		}
		req = req.OrderBy("created")
	}
	req = req.Limit(uint64(params.Limit))
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
	log.Info(getThreadsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := u.DB.Query(getThreadsQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// id, title, author, message, forum, votes, slug, created
	resThreads := make([]domain.Thread, 0)
	for rows.Next() {
		var currentThread domain.Thread
		var slug *string
		var created *time.Time
		err := rows.Scan(&currentThread.ID, &currentThread.Title, &currentThread.Author, &currentThread.Forum, &currentThread.Message, &slug, &created, &currentThread.Votes)
		if err != nil {
			return nil, err
		}
		if slug != nil {
			currentThread.Slug = *slug
		}
		if created != nil {
			currentThread.Created = created.Format(constants.TimeLayout)
		}
		resThreads = append(resThreads, currentThread)
	}

	return resThreads, nil
}
