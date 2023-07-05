package setup

import (
	"context"
	"flag"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/usecases"
	httphandlers "github.com/SOAT1StackGoLang/tech-challenge/internal/handlers/http"
	pgxrepo "github.com/SOAT1StackGoLang/tech-challenge/internal/repositories/postgres"
	"github.com/emicklei/go-restful/v3"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

var (
	binding    string
	log        = helpers.NewLogger()
	connString string
)

func init() {
	flag.StringVar(&binding, "httpbind", ":8000", "address/port to bind listen socket")
	flag.Parse()
	godotenv.Load()
	helpers.ReadPgxConnEnvs()
	connString = helpers.ToDsnWithDbName()
}

func SetupCode() {
	ctx := context.Background()

	gormDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	userRepo := pgxrepo.NewPgxUsersRepository(gormDB, log)
	userUseCase := usecases.NewUsersUseCase(userRepo, log)

	catRepo := pgxrepo.NewPgxCategoriesRepository(gormDB, log)
	catUseCase := usecases.NewCategoriesUseCase(log, catRepo, userUseCase)

	ws := new(restful.WebService)
	httphandlers.NewUserHandler(ctx, userUseCase, ws)
	httphandlers.NewCategoriesHttpHandler(ctx, catUseCase, ws)

	restful.Add(ws)

	log.Info("listening...")
	log.Panic(http.ListenAndServe(binding, nil))
}
