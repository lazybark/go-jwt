package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-jwt/config"
	"github.com/lazybark/go-jwt/storage"

	"github.com/lazybark/lazyevent/v2/events"
	"github.com/lazybark/lazyevent/v2/lproc"
)

type Api struct {
	ver    semver.Ver
	db     storage.Storage
	conf   *config.Config
	logger *lproc.LogProcessor
}

/*
type ILogger interface {
	Log(e events.Event)
	FatalInCaseErr(any)
}*/

var ApiVer = semver.Ver{
	Major:         1,
	Minor:         0,
	Patch:         0,
	ReleaseNote:   "RED-backend service api",
	BuildMetadata: "",
	Stable:        true,
}

func New(db storage.Storage, conf config.Config, l *lproc.LogProcessor) *Api {

	return &Api{ver: ApiVer, db: db, conf: &conf, logger: l}
}

func (a *Api) Start() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		//Check exact api version
		r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf(ApiStringResult, a.ver.String())))
		})
		//Create new user
		r.With(parseFormMiddleware).Post("/users/add", func(w http.ResponseWriter, r *http.Request) {
			a.ResponseUserAdd(r, w)
		})
		//Authenticate user
		r.With(parseFormMiddleware).Post("/users/login", func(w http.ResponseWriter, r *http.Request) {
			a.ResponseUserLogin(r, w)
		})
	})

	a.logger.Log(events.Info(fmt.Sprintf("Listening on %s", a.conf.Host)).Green())
	http.ListenAndServe(a.conf.Host, r)
}

func (a *Api) StorageFlush() error {
	return a.db.Flush()
}

func (a *Api) StorgageInit() error {
	return a.db.Init()
}
