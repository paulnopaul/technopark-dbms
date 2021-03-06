package usecase

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"strconv"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/post"
	"technopark-dbms/internal/pkg/thread"
	"technopark-dbms/internal/pkg/utilities"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type threadUsecase struct {
	DB     *pgx.ConnPool
	UUCase domain.UserUsecase
}

func (t threadUsecase) GetThreadIdAndForum(s utilities.SlugOrId) (*domain.Thread, error) {
	query := "select id, forum from threads where"
	args := make([]interface{}, 0)
	if s.IsSlug {
		query += " slug = $1;"
		args = append(args, s.Slug)
	} else {
		query += " id = $1;"
		args = append(args, s.ID)
	}

	resThread := &domain.Thread{}
	err := t.DB.QueryRow(query, args...).
		Scan(&resThread.ID, &resThread.Forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, thread.NotFound
		}
		return nil, err
	}
	return resThread, nil
}

func (t threadUsecase) CreatePosts(s utilities.SlugOrId, posts domain.PostArray) (domain.PostArray, error) {
	threadInfo, err := t.GetThreadIdAndForum(s)
	if err == thread.NotFound {
		return nil, thread.NotFound
	} else if err != nil {
		return nil, err
	} else if threadInfo == nil {
		return nil, errors.New("WTF")
	}
	if len(posts) == 0 {
		return posts, nil
	}

	now := strfmt.DateTime(time.Now())
	//for i, _ := range posts {
	//	posts[i].Created = now
	//	posts[i].Thread = threadInfo.ID
	//	posts[i].Forum = threadInfo.Forum
	//}

	//query := "insert into posts(parent, author, message, is_edited, thread, created, forum) values ($1, $2, $3, $4, $5, $6, $7) returning id;"
	req := psql.Insert("posts(parent, author, message, is_edited, thread, created, forum)")
	userReq := psql.Insert("f_u(f, u)")
	for i, _ := range posts {
		posts[i].Created = now
		posts[i].Thread = threadInfo.ID
		posts[i].Forum = threadInfo.Forum
		req = req.Values(posts[i].Parent, posts[i].Author, posts[i].Message, posts[i].IsEdited, posts[i].Thread, posts[i].Created, posts[i].Forum)
		userReq = userReq.Values(threadInfo.Forum, posts[i].Author)
	}
	req = req.Suffix("returning id")
	userReq = userReq.Suffix("on conflict do nothing")
	query, args, _ := req.ToSql()

	tx, err := t.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&posts[i].ID)
		if err != nil {
			return nil, err
		}
	}
	if rows.Err() != nil {
		if pgErr, ok := rows.Err().(pgx.PgError); ok {
			if pgErr.Code == "66666" {
				return nil, post.InvalidParentError
			}
			if pgErr.Code == "23503" {
				return nil, thread.AuthorNotExists
			}
			return nil, rows.Err()
		}
	}
	query = "update forums set posts = posts + $1 where slug = $2;"
	_, err = tx.Exec(query, len(posts), threadInfo.Forum)
	if err != nil {
		return nil, err
	}

	query, args, _ = userReq.ToSql()
	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (t threadUsecase) GetThreadDetails(s utilities.SlugOrId) (*domain.Thread, error) {
	query := "select id, title, author, message, votes, forum, slug, created from threads where "
	args := make([]interface{}, 0)
	if s.IsSlug {
		query += "slug = $1;"
		args = append(args, s.Slug)
	} else {
		query += "id = $1;"
		args = append(args, s.ID)
	}

	var slug *string
	resThread := &domain.Thread{}
	err := t.DB.QueryRow(query, args...).
		Scan(&resThread.ID, &resThread.Title, &resThread.Author, &resThread.Message, &resThread.Votes, &resThread.Forum, &slug, &resThread.Created)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, thread.NotFound
		}
		return nil, err
	}
	if slug != nil {
		resThread.Slug = *slug
	}
	return resThread, nil
}

func (t threadUsecase) UpdateThreadDetails(s utilities.SlugOrId, threadUpdate domain.Thread) (*domain.Thread, error) {
	threadDetails, err := t.GetThreadDetails(s)
	if err != nil {
		return nil, err
	}
	if threadUpdate.Message == "" && threadUpdate.Title == "" {
		return threadDetails, nil
	}

	updateThreadQuery := "update threads set title=coalesce(nullif($1, ''), title), message=coalesce(nullif($2, ''), message) where id = $3;"
	_, err = t.DB.Exec(updateThreadQuery, threadUpdate.Title, threadUpdate.Message, threadDetails.ID)
	if err != nil {
		return nil, err
	}
	if threadUpdate.Title != "" {
		threadDetails.Title = threadUpdate.Title
	}
	if threadUpdate.Message != "" {
		threadDetails.Message = threadUpdate.Message
	}
	return threadDetails, nil
}

func parentPostsQuery(id int32, limit int, since int64, desc bool) (string, []interface{}) {
	order, s := "asc", " > "
	if desc {
		order = "desc"
	} else {
		order = "asc"
	}
	if desc && since != 0 {
		s = " < "
	} else {
		s = " > "
	}
	args := []interface{}{id, limit}
	query := `select p.id, p.parent, p.author, p.message, p.is_edited, p.forum, p.thread, p.created
			from posts p where p.way[2] in (select id from posts where thread = $1 and way[3] is null`
	if since != 0 {
		query += " and way[2] " + s + "(select way[2] from posts where id = $3)"
		args = append(args, since)
	}
	query += " order by id " + order + " limit $2) order by p.way[2] " + order + ",p.way asc, p.id asc"
	return query, args
}

func flatPostsQuery(id int32, limit int, since int64, desc bool) (string, []interface{}) {
	order, s := "asc", " > "
	if desc {
		order, s = "desc", " < "
	}
	args := []interface{}{id, limit}
	query := `select p.id, p.parent, p.author, p.message, p.is_edited, p.forum, p.thread, p.created
				from posts p where p.thread = $1`
	if since != 0 {
		query += " and p.id " + s + " $3"
		args = append(args, since)
	}
	query += " order by p.id " + order + " limit $2 "
	return query, args
}

func treePostsQuery(id int32, limit int, since int64, desc bool) (string, []interface{}) {
	order, s := "asc", " > "
	if desc {
		order = "desc"
	}
	if desc && since != 0 {
		s = " < "
	}
	args := []interface{}{id, limit}
	query := `select p.id, p.parent, p.author, p.message, p.is_edited, p.forum, p.thread, p.created
		from posts p where p.thread = $1`
	if since != 0 {
		query += "and way " + s + "(select way from posts where id = $3)"
		args = append(args, since)
	}
	query += "order by p.way " + order + ", p.created " + order + ", p.id asc limit $2 "
	return query, args
}

func generateGetPostsQuery(threadId int32, params utilities.ArrayOutParams) (string, []interface{}, error) {
	since := int64(0)
	if params.Since != "" {
		parsedSince, err := strconv.ParseInt(params.Since, 10, 64)
		if err != nil {
			return "", nil, err
		}
		since = parsedSince
	}
	var query string
	var args []interface{}
	switch params.Sort {
	case "flat":
		query, args = flatPostsQuery(threadId, int(params.Limit), since, params.Desc)
	case "tree":
		query, args = treePostsQuery(threadId, int(params.Limit), since, params.Desc)
	case "parent_tree":
		query, args = parentPostsQuery(threadId, int(params.Limit), since, params.Desc)
	default:
		query, args = flatPostsQuery(threadId, int(params.Limit), since, params.Desc)
	}
	return query, args, nil
}

func (t threadUsecase) GetThreadPosts(s utilities.SlugOrId, params utilities.ArrayOutParams) (domain.PostArray, error) {
	threadDetails, err := t.GetThreadIdAndForum(s)
	if err != nil {
		return nil, err
	}

	getPostsQuery, args, err := generateGetPostsQuery(threadDetails.ID, params)
	if err != nil {
		return nil, err
	}

	rows, err := t.DB.Query(getPostsQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// p.id, p.parent, p.author, p.message, p.is_edited, p.forum, p.thread, p.created
	resPosts := make(domain.PostArray, 0)
	for rows.Next() {
		var p domain.Post
		//var created *strfmt.DateTime
		err := rows.Scan(&p.ID, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
		if err != nil {
			return nil, err
		}
		//if created != nil {
		//	p.Created = *created
		//}
		resPosts = append(resPosts, p)
	}

	return resPosts, nil
}

func (t threadUsecase) CreateThreadVote(s utilities.SlugOrId, vote domain.Vote) (*domain.Thread, error) {
	threadDetails, err := t.GetThreadDetails(s)
	if err != nil {
		return nil, err
	}

	var currentVoice int32
	getVoteQuery := "select voice from votes where thread = $1 and username = $2;"
	err = t.DB.QueryRow(getVoteQuery, threadDetails.ID, vote.Nickname).Scan(&currentVoice)
	if err != nil {
		if err == pgx.ErrNoRows {
			newVoteQuery := "insert into votes(thread, username, voice) values ($1, $2, $3);"
			_, err = t.DB.Exec(newVoteQuery, threadDetails.ID, vote.Nickname, vote.Voice)
			if err != nil {
				return nil, thread.AuthorNotExists
			}
			threadDetails.Votes += vote.Voice
			return threadDetails, nil
		}
		return nil, err
	}
	updateVoteQuery := "update votes set voice=$3 where thread=$1 and username=$2"
	_, err = t.DB.Exec(updateVoteQuery, threadDetails.ID, vote.Nickname, vote.Voice)
	if err != nil {
		return nil, err
	}
	if vote.Voice != currentVoice {
		threadDetails.Votes = threadDetails.Votes - currentVoice + vote.Voice
	}
	return threadDetails, nil
}

func NewThreadUsecase(db *pgx.ConnPool, userUsecase domain.UserUsecase) domain.ThreadUsecase {
	return &threadUsecase{
		DB:     db,
		UUCase: userUsecase,
	}
}
