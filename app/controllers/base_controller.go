package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"simpleWebCart/app/database/seeders"
	"simpleWebCart/app/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	Router    *mux.Router
	AppConfig *AppConfig
}

type AppConfig struct {
	AppName string
	AppURL  string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

type PageLink struct {
	Page          int32
	Url           string
	IsCurrentPage bool
}

type PaginationLinks struct {
	CurrentPage string
	NextPage    string
	PrevPage    string
	TotalRows   int32
	TotalPages  int32
	Links       []PageLink
}

type PaginationParams struct {
	Path        string
	TotalRows   int32
	PerPage     int32
	CurrentPage int32
}

var (
	store               = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	sessionShoppingCart = "shopping-cart-session"
	sessionUser         = "user-session"
	sessionFlash        = "flash-session"
)

func (server *Server) Initialize(dbConf DBConfig) {

	server.InitializeDB(dbConf)
	server.InitializeRoutes()
}

func (Server *Server) Run(addr string) {
	fmt.Printf("Listening on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, Server.Router))
}

func (server *Server) InitializeDB(dbConf DBConfig) {

	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect databse")
	}

}

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated Success")
}

func (server *Server) InitCommands(dbConf DBConfig) {
	server.InitializeDB(dbConf)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func GetPaginationLinks(conf *AppConfig, params PaginationParams) (PaginationLinks, error) {
	var (
		links    []PageLink
		nextPage int32
		prevPage int32
	)

	totalPages := int32(math.Ceil(float64(params.TotalRows) / float64(params.PerPage)))

	for i := 1; int32(i) <= totalPages; i++ {
		links = append(links, PageLink{
			Page:          int32(i),
			Url:           fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprint(i)),
			IsCurrentPage: int32(i) == params.CurrentPage,
		})
	}

	prevPage = 1
	nextPage = totalPages

	if params.CurrentPage > 2 {
		prevPage = params.CurrentPage - 1
	}

	if params.CurrentPage < totalPages {
		nextPage = params.CurrentPage + 1
	}

	return PaginationLinks{
		CurrentPage: fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprint(params.CurrentPage)),
		NextPage:    fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprint(nextPage)),
		PrevPage:    fmt.Sprintf("%s/%s?page=%s", os.Getenv("APP_URL"), params.Path, fmt.Sprint(prevPage)),
		TotalRows:   params.TotalRows,
		TotalPages:  totalPages,
		Links:       links,
	}, nil
}

func SetFlash(w http.ResponseWriter, r *http.Request, name, value string) {
	sessions, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions.AddFlash(value, name)
	sessions.Save(r, w)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) []string {
	sessions, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	fm := sessions.Flashes(name)
	if len(fm) < 0 {
		return nil
	}

	sessions.Save(r, w)
	var flashes []string
	for _, fl := range fm {
		flashes = append(flashes, fl.(string))
	}

	return flashes
}

func IsLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, sessionUser)
	if session.Values["id"] == nil {
		return false
	}

	return true
}

func ComparePassword(password string, hashedPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password)) == nil
}

func (server *Server) CurrentUser(w http.ResponseWriter, r *http.Request) *models.User {
	if !IsLoggedIn(r) {
		return nil
	}

	session, _ := store.Get(r, sessionUser)

	userModel := models.User{}
	user, err := userModel.FindByID(server.DB, session.Values["id"].(string))
	if err != nil {
		session.Values["id"] = nil
		session.Save(r, w)
		return nil
	}

	return user
}
