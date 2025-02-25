package costco

import (
	"costco/internal/metadata"
	"costco/internal/metadata/postgres"
	"costco/internal/storage"
	"costco/internal/storage/minio"

	"github.com/gin-gonic/gin"
)

type registry struct {
	g *gin.Engine
	m metadata.MetadataStorage
	s storage.Storage
}

func Init() registry {
	g := gin.Default()

	m := postgres.Init()
	s := minio.Init()

	r := registry{
		g: g,
		m: &m,
		s: &s,
	}

	return r
}

func (r *registry) RegisterEndpoints() {
	r.g.GET("/get", r.Get)
	r.g.GET("/", r.Metadata)
	r.g.POST("/post", r.Put)
}

func (r *registry) Start() {
	r.g.Run()
}
