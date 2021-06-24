package app

import (
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
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

	postUsecase := postDBUsecase.NewPostUsecase(db)
	threadUsecase := threadDBUsecase.NewThreadUsecase(db)
	serviceUsecase := serviceDBUsecase.NewServiceUsecase(db)
	userUsecase := userDBUsecase.NewUserUsecase(db)
	forumUsecase := forumDBUsecase.NewForumUsecase(db, userUsecase)

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
