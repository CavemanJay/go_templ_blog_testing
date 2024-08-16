package main

import (
	"blog/article"
	"blog/sqlite"
	"blog/views"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"embed"
)

//go:embed static
var staticFiles embed.FS

const UseEmbededStaticFiles = false

func HomeHandler(ctx echo.Context, db *sqlite.Queries) error {
	articles, err := db.QueryArticles(ctx.Request().Context())
	if err != nil {
		return err
	}

	return views.Home(articles).Render(ctx.Request().Context(), ctx.Response())
}

func ArticleHandler(ctx echo.Context, db *sqlite.Queries, parser article.Parser) error {
	slug := ctx.Param("slug")
	article, err := db.QueryArticleBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	articleContent, err := parser.Parse(article.Filename)
	if err != nil {
		return err
	}

	return views.Article(article.Title, articleContent).Render(ctx.Request().Context(), ctx.Response())
}

// func ErrorHandler() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) (err error) {
// 			if err = next(c); err != nil {
// 				c.Error(err)
// 			}
// 			return
// 		}
// 	}
// }

func main() {
	var (
		dbCnn  = mustSetupDatabaseConnection()
		db     = sqlite.New(dbCnn)
		parser = article.NewParser()
		e      = echo.New()
		fs     = echo.MustSubFS(staticFiles, "static")
	)

	if UseEmbededStaticFiles {
		e.StaticFS("/static", fs)
	} else {
		e.Static("/static", "static")
	}

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return HomeHandler(c, db)
	})
	e.GET("/articles/:slug", func(c echo.Context) error {
		return ArticleHandler(c, db, parser)
	})

	e.HTTPErrorHandler = func(err error, c echo.Context) {

		if he, ok := err.(*echo.HTTPError); ok && he.Code == 404 {
			e.DefaultHTTPErrorHandler(he, c)
			return
		}

		if e := views.Error(err).Render(c.Request().Context(), c.Response()); e != nil {
			panic(e)
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}

func mustSetupDatabaseConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		panic(err)
	}
	return db
}
