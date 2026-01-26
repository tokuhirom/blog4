package public

import (
	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/internal"
)

func SetupPublicRoutes(r *gin.Engine, queries *publicdb.Queries, cfg *internal.Config) {
	r.GET("/", func(c *gin.Context) {
		RenderTopPage(c, queries)
	})
	r.GET("/feed", func(c *gin.Context) {
		RenderFeed(c, queries)
	})
	r.GET("/entry/*filepath", func(c *gin.Context) {
		RenderEntryPage(c, queries, cfg)
	})
	r.GET("/search", func(c *gin.Context) {
		RenderSearchPage(c, queries)
	})
	r.StaticFile("/static/main.css", "public/static/main.css")
	r.StaticFile("/build-info.json", "build-info.json")
}
