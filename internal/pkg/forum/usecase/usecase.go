package usecase

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/forum"
	"technopark-dbms/internal/pkg/thread"
	"technopark-dbms/internal/pkg/utilities"
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

func (u *forumUsecase) ForumExists(slug string) (bool, error) {
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
	authorExists, err := u.UUCase.UserExists(f.User, "")
	if err != nil {
		return nil, err
	} else if !authorExists {
		return nil, forum.AuthorNotExists
	}

	foundForum, err := u.GetForumDetails(f.Slug)
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

func (u *forumUsecase) GetForumDetails(slug string) (*domain.Forum, error) {
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
	values := make([]interface{}, 0)
	values = append(values, t.Author, forumSlug, t.Message, t.Title)
	req := "insert into threads(author, forum, message, title, created, slug) values ($1, $2, $3, $4, $5, $6)"
	if t.Created.String() != "" {
		values = append(values, t.Created)
	} else {
		values = append(values, nil)
	}
	if t.Slug != "" {
		values = append(values, t.Slug)
	} else {
		values = append(values, nil)
	}
	req += "returning id, author, (select f.slug from forums f where f.slug = $2), message, title, created, slug;"
	return req, values, nil
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
	authorExists, err := u.UUCase.UserExists(t.Author, "")
	if err != nil {
		return nil, err
	} else if !authorExists {
		return nil, forum.AuthorNotExists
	}

	forumExists, err := u.ForumExists(forumSlug)
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
	err = u.DB.QueryRow(createThreadQuery, args...).
		Scan(&newThread.ID, &newThread.Author, &newThread.Forum, &newThread.Message, &newThread.Title, &newThread.Created, &slug)
	if err != nil {
		return nil, err
	}
	if slug != nil {
		newThread.Slug = *slug
	}
	return newThread, nil
}

func generateUserRequest(slug string, params utilities.ArrayOutParams) (string, []interface{}, error) {
	var order string
	var s string
	if params.Desc {
		order, s = "desc", " < "
	} else {
		order, s = "asc", " > "
	}
	var query string
	args := make([]interface{}, 0)
	if params.Since != "" {
		query = "select u, fullname, about, email from f_u join users on f_u.u = users.nickname where f = $1 and u " + s + " $2 order by nickname " + order + " limit $3;"
		args = append(args, slug, params.Since, params.Limit)
	} else {
		query = "select u, fullname, about, email from f_u join users on f_u.u = users.nickname where f = $1 order by nickname " + order + " limit $2;"
		args = append(args, slug, params.Limit)
	}

	return query, args, nil
}

func (u *forumUsecase) GetUsers(forumSlug string, params utilities.ArrayOutParams) (domain.UserArray, error) {
	forumExists, err := u.ForumExists(forumSlug)
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

	resUsers := make(domain.UserArray, 0)
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

func (u *forumUsecase) GetThreads(forumSlug string, params utilities.ArrayOutParams) (domain.ThreadArray, error) {
	forumExists, err := u.ForumExists(forumSlug)
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
	resThreads := make(domain.ThreadArray, 0)
	for rows.Next() {
		var currentThread domain.Thread
		//var created *strfmt.DateTime
		var slug *string
		err := rows.Scan(&currentThread.ID, &currentThread.Title, &currentThread.Author, &currentThread.Forum, &currentThread.Message, &slug, &currentThread.Created, &currentThread.Votes)
		if err != nil {
			return nil, err
		}
		//if created != nil {
		//	currentThread.Created = *created
		//}
		if slug != nil {
			currentThread.Slug = *slug
		}
		resThreads = append(resThreads, currentThread)
	}

	return resThreads, nil
}
