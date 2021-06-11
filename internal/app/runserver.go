package app

import (
	forumDelivery "DBMSForum/internal/pkg/forum/delivery"
	forumDBUsecase "DBMSForum/internal/pkg/forum/usecase"
	"DBMSForum/internal/pkg/middlewares"
	postDelivery "DBMSForum/internal/pkg/post/delivery"
	postDBUsecase "DBMSForum/internal/pkg/post/usecase"
	serviceDelivery "DBMSForum/internal/pkg/service/delivery"
	serviceDBUsecase "DBMSForum/internal/pkg/service/usecase"
	threadDelivery "DBMSForum/internal/pkg/thread/delivery"
	threadDBUsecase "DBMSForum/internal/pkg/thread/usecase"
	userDelivery "DBMSForum/internal/pkg/user/delivery"
	userDBUsecase "DBMSForum/internal/pkg/user/usecase"
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

	forumUsecase := forumDBUsecase.NewForumUsecase(db)
	postUsecase := postDBUsecase.NewPostUsecase(db)
	threadUsecase := threadDBUsecase.NewThreadUsecase(db)
	serviceUsecase := serviceDBUsecase.NewServiceUsecase(db)
	userUsecase := userDBUsecase.NewUserUsecase(db)

	forumDelivery.NewForumHandler(r, forumUsecase)
	postDelivery.NewPostHandler(r, postUsecase)
	threadDelivery.NewThreadHandler(r, threadUsecase)
	serviceDelivery.NewServiceHandler(r, serviceUsecase)
	userDelivery.NewUserHandler(r, userUsecase)

	log.Println("Listening at: ", addr)
	err = fasthttp.ListenAndServe(addr, middlewares.Logging(r.Handler))
	if err != nil {
		log.Println(fmt.Sprint("Server error: ", err))
	}
}
