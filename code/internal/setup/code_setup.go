package setup

import (
	"context"
	"embed"
	"flag"
	"net/http"

	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/usecases"
	httphandlers "github.com/SOAT1StackGoLang/tech-challenge/internal/handlers/http"
	pgxrepo "github.com/SOAT1StackGoLang/tech-challenge/internal/repositories/postgres"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

// swagger embed files in binary
// dist folder from https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.1.0.zip
//
//go:embed apidocs/*
var swaggerUI embed.FS

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

	// Configure Swagger and Redirect / to /apidocs/
	configureSwagger()
	ws := new(restful.WebService)
	httphandlers.NewUserHandler(ctx, userUseCase, ws)
	httphandlers.NewCategoriesHttpHandler(ctx, catUseCase, ws)

	restful.Add(ws)

	// Configure Swagger and Redirect / to /apidocs/
	configureSwagger()

	ws := new(restful.WebService)
	httphandlers.NewUserHandler(ctx, userUseCase, ws)
	httphandlers.NewCategoriesHttpHandler(ctx, catUseCase, ws)

	restful.Add(ws)

	log.Info("listening...")
	log.Panic(http.ListenAndServe(binding, nil))
}

func configureSwagger() {

	// Serve Swagger UI files
	fs := http.FileServer(http.FS(swaggerUI))
	http.Handle("/apidocs/", http.StripPrefix("/", fs))

	// Set up Swagger API
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       "/apidocs/openapi.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}
	restful.Add(restfulspec.NewOpenAPIService(config))

	// Redirect root to /apidocs/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/apidocs/", http.StatusMovedPermanently)
	})
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}
