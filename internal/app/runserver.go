package app

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	forumDelivery "technopark-dbms/internal/pkg/forum/delivery"
	forumDBUsecase "technopark-dbms/internal/pkg/forum/usecase"
	"technopark-dbms/internal/pkg/middlewares"
	postDelivery "technopark-dbms/internal/pkg/post/delivery"
	postDBUsecase "technopark-dbms/internal/pkg/post/usecase"
	serviceDelivery "technopark-dbms/internal/pkg/service/delivery"
	serviceDBUsecase "technopark-dbms/internal/pkg/service/usecase"
	threadDelivery "technopark-dbms/internal/pkg/thread/delivery"
	threadDBUsecase "technopark-dbms/internal/pkg/thread/usecase"
	userDelivery "technopark-dbms/internal/pkg/user/delivery"
	userDBUsecase "technopark-dbms/internal/pkg/user/usecase"
)

// REQUIRES POSTGRES DRIVER IN IMPORT
func getPostgres() *pgx.ConnPool {
	conf := pgx.ConnConfig{
		User:                 "dbmsmaster",
		Database:             "dbmsforum",
		Password:             "dbms",
		PreferSimpleProtocol: false,
	}
	poolConf := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 11,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	conn, err := pgx.NewConnPool(poolConf)
	if err != nil {
		log.WithError(err).Fatal("postgres connection error")
	}
	return conn
}

func RunServer(addr string) {
	r := router.New()

	db := getPostgres()

	postUsecase := postDBUsecase.NewPostUsecase(db)
	serviceUsecase := serviceDBUsecase.NewServiceUsecase(db)
	userUsecase := userDBUsecase.NewUserUsecase(db)
	threadUsecase := threadDBUsecase.NewThreadUsecase(db, userUsecase)
	forumUsecase := forumDBUsecase.NewForumUsecase(db, userUsecase, threadUsecase)

	forumDelivery.NewForumHandler(r, forumUsecase)
	postDelivery.NewPostHandler(r, postUsecase)
	threadDelivery.NewThreadHandler(r, threadUsecase)
	serviceDelivery.NewServiceHandler(r, serviceUsecase)
	userDelivery.NewUserHandler(r, userUsecase)

	log.Println("Listening at: ", addr)
	err := fasthttp.ListenAndServe(addr, middlewares.Logging(r.Handler))
	if err != nil {
		log.Println(fmt.Sprint("Server error: ", err))
	}
}
