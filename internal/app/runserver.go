package app

import (
	forumDelivery "DBMSForum/internal/pkg/forum/delivery"
	"DBMSForum/internal/pkg/forum/usecase"
	"DBMSForum/internal/pkg/middlewares"
	postDelivery "DBMSForum/internal/pkg/post/delivery"
	usecase2 "DBMSForum/internal/pkg/post/usecase"
	serviceDelivery "DBMSForum/internal/pkg/service/delivery"
	threadDelivery "DBMSForum/internal/pkg/thread/delivery"
	userDelivery "DBMSForum/internal/pkg/user/delivery"
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	// HERE MUST BE POSTGRES DRIVER
	_ "github.com/jackc/pgx/stdlib"
)

// REQUIRES POSTGRES DRIVER IN IMPORT
func getPostgres() (*sql.DB, error) {
	dsn := "user=dbmsmaster dbname=dbmsforum password=dbms host=127.0.0.1 port=5432 sslmode=disable"
	log.Info("Connecting postgres")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
		return nil, err
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.WithError(err).Fatal("Postgres connection error")
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}

func RunServer(addr string) {
	r := router.New()

	db, err := getPostgres()
	if err != nil {
		log.WithError(err).Error("Error while getting db")
		return
	}
	forumUsecase := usecase.NewForumUsecase(db)
	postUsecase := usecase2.NewPostUsecase(db)

	forumDelivery.NewForumHandler(r, forumUsecase)
	postDelivery.NewPostHandler(r, postUsecase)
	threadDelivery.NewThreadHandler(r)
	serviceDelivery.NewServiceHandler(r)
	userDelivery.NewUserHandler(r)

	log.Println("Listening at: ", addr)
	err = fasthttp.ListenAndServe(addr, middlewares.Logging(r.Handler))
	if err != nil {
		log.Println(fmt.Sprint("Server error: ", err))
	}
}
