package usecase

import (
	"errors"
	"github.com/jackc/pgx"
	"technopark-dbms/internal/pkg/constants"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/thread"
	"technopark-dbms/internal/pkg/utilities"
	"time"
)

type threadUsecase struct {
	DB     *pgx.ConnPool
	UUCase domain.UserUsecase
}

func (t threadUsecase) CreatePosts(s utilities.SlugOrId, posts []domain.Post) ([]domain.Post, error) {
	threadInfo, err := t.GetThreadDetails(s)
	if err == thread.NotFound {
		return nil, thread.NotFound
	} else if err != nil {
		return nil, err
	} else if threadInfo == nil {
		return nil, errors.New("WTF")
	}

	now := time.Now().Format(constants.TimeLayout)
	for i, _ := range posts {
		posts[i].Created = now
		posts[i].Thread = threadInfo.ID
		posts[i].Forum = threadInfo.Forum
	}

	query := "insert into posts(parent, author, message, is_edited, thread, created, forum) values ($1, $2, $3, $4, $5, $6, $7) returning id;"
	for i, _ := range posts {
		err := t.DB.QueryRow(query, posts[i].Parent, posts[i].Author, posts[i].Message, posts[i].IsEdited, posts[i].Thread, posts[i].Created, posts[i].Forum).
			Scan(&posts[i].ID)
		if err != nil {
			return nil, err
		}
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

	//var slug *string
	var created *time.Time
	resThread := &domain.Thread{}
	err := t.DB.QueryRow(query, args...).
		Scan(&resThread.ID, &resThread.Title, &resThread.Author, &resThread.Message, &resThread.Votes, &resThread.Forum, &resThread.Slug, &created)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, thread.NotFound
		}
		return nil, err
	}
	if created != nil {
		resThread.Created = created.Format(constants.TimeLayout)
	}
	//if slug != nil {
	//	resThread.Slug = *slug
	//}

	return resThread, nil
}

func (t threadUsecase) UpdateThreadDetails(s utilities.SlugOrId, threadUpdate domain.Thread) (*domain.Thread, error) {
	panic("implement me")
}

func (t threadUsecase) GetThreadPosts(s utilities.SlugOrId, params utilities.ArrayOutParams) ([]domain.Post, error) {

}

func (t threadUsecase) VoteThread(s utilities.SlugOrId, vote domain.Vote) (*domain.Thread, error) {
	threadDetails, err := t.GetThreadDetails(s)
	if err != nil {
		return nil, err
	}

	userExists, err := t.UUCase.Exists(vote.Nickname, "")
	if err != nil {
		return nil, err
	} else if !userExists {
		return nil, thread.AuthorNotExists
	}

	var currentVoice int32
	getVoteQuery := "select voice from votes where thread = $1 and username = $2;"
	err = t.DB.QueryRow(getVoteQuery, threadDetails.ID, vote.Nickname).Scan(&currentVoice)
	if err != nil {
		if err == pgx.ErrNoRows {
			newVoteQuery := "insert into votes(thread, username, voice) values ($1, $2, $3);"
			_, err = t.DB.Exec(newVoteQuery, threadDetails.ID, vote.Nickname, vote.Voice)
			if err != nil {
				return nil, err
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
