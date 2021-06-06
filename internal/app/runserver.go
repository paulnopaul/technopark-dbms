package app

import (
	"DBMSForum/internal/pkg/forum/delivery"
	"DBMSForum/internal/pkg/forum/repostory/db"
	"DBMSForum/internal/pkg/forum/usecase"
	"DBMSForum/internal/pkg/middlewares"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/fasthttp/router"

	// HERE MUST BE POSTGRES DRIVER
	_ "github.com/jackc/pgx/stdlib"
	"github.com/valyala/fasthttp"
)

// REQUIRES POSTGRES DRIVER IN IMPORT
func getPostgres() (*sql.DB, error {
	dsn := "user=postgres dbname=golang password=love host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
		return nil, err
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}

func RunServer(addr string) {
	r := router.New()

	forumRepo := db.NewForumRepository()
	forumUsecase := usecase.NewForumUsecase(forumRepo)
	delivery.NewForumHandler(r, forumUsecase)

	log.Println("Listening at: ", addr)
	err := fasthttp.ListenAndServe(addr, middlewares.Logging(r.Handler))

	if err != nil {
		log.Println(fmt.Sprint("Server error: ", err))
	}
}
