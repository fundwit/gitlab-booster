package servehttp

import (
	"gitlab-booster/persistence"

	"github.com/fundwit/go-commons/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/sony/sonyflake"
)

type manifestHandler struct {
	idWorker  *sonyflake.Sonyflake
	validator *validator.Validate
	ds        *persistence.DatasourceManager
}

// RegisterManifestHandlers 注册路由
func RegisterManifestHandlers(root *gin.Engine, dsManager *persistence.DatasourceManager, middlewares ...gin.HandlerFunc) {
	r := root.Group("/v1/manifests")
	r.Use(middlewares...)

	handler := &manifestHandler{
		idWorker:  sonyflake.NewSonyflake(sonyflake.Settings{}),
		validator: validator.New(),
		ds:        dsManager,
	}

	r.POST("", handler.newManifest)
	r.GET("", handler.queryManifests)
	r.DELETE(":id", handler.deleteManifest)
	r.GET(":id", handler.manifestDetail)
	r.PUT(":id", handler.updateManifest)

	r.POST(":id/items", handler.appendManifestItems)
	r.DELETE(":id/items", handler.removeManifestItems)
}

// Manifest 配置清单
type Manifest struct {
	ID   types.ID `json:"id"`
	Name string   `json:"name" validate:"required"`
	Note string   `json:"note"`
}

type ManifestUpdate struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

type ManifestWithRepos struct {
	Manifest
	Repos []RepositoryRef `json:"repos"`
}

type pathParams struct {
	ID types.ID `uri:"id" validate:"required"`
}
type RepositoryRef struct {
	GroupID types.ID `uri:"id" json:"groupId" validate:"required" gorm:"primary_key" sql:"type:BIGINT UNSIGNED NOT NULL"`

	ID        types.ID `json:"id" validate:"required" gorm:"primary_key" sql:"type:BIGINT UNSIGNED NOT NULL"`
	Name      string   `json:"name" validate:"required"`
	Namespace string   `json:"namespace" validate:"required"`
	URL       string   `json:"url" validate:"required"`
}

func (m *manifestHandler) newManifest(c *gin.Context) {
	body := Manifest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(body); err != nil {
		panic(err)
	}
	id, err := m.idWorker.NextID()
	if err != nil {
		panic(err)
	}
	body.ID = types.ID(id)
	if err := m.ds.GromDB().Save(&body).Error; err != nil {
		panic(err)
	}
	c.JSON(201, &body)
}

func (m *manifestHandler) queryManifests(c *gin.Context) {
	list := []Manifest{}
	if err := m.ds.GromDB().Model(&Manifest{}).Scan(&list).Error; err != nil {
		panic(err)
	}
	c.JSON(200, &list)
}

func (m *manifestHandler) updateManifest(c *gin.Context) {
	p := pathParams{}
	if err := c.ShouldBindUri(&p); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(p); err != nil {
		panic(err)
	}

	body := ManifestUpdate{}
	if err := c.ShouldBindJSON(&body); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(body); err != nil {
		panic(err)
	}

	err := m.ds.GromDB().Transaction(func(tx *gorm.DB) error {
		if err := m.ds.GromDB().Model(&Manifest{}).Where(&Manifest{ID: p.ID}).Update(body).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	c.AbortWithStatus(200)
}

func (m *manifestHandler) deleteManifest(c *gin.Context) {
	p := pathParams{}
	if err := c.ShouldBindUri(&p); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(p); err != nil {
		panic(err)
	}

	err := m.ds.GromDB().Transaction(func(tx *gorm.DB) error {
		if err := m.ds.GromDB().Model(&Manifest{}).Delete(&Manifest{ID: p.ID}).Error; err != nil {
			return err
		}
		if err := m.ds.GromDB().Where(&RepositoryRef{GroupID: p.ID}).Delete(&RepositoryRef{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	c.AbortWithStatus(204)
}

func (m *manifestHandler) manifestDetail(c *gin.Context) {
	p := pathParams{}
	if err := c.ShouldBindUri(&p); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(p); err != nil {
		panic(err)
	}
	man := Manifest{}
	if err := m.ds.GromDB().Model(&Manifest{}).First(&man, &Manifest{ID: p.ID}).Error; err != nil {
		panic(err)
	}

	out := ManifestWithRepos{
		Manifest: man,
	}

	if err := m.ds.GromDB().Model(&RepositoryRef{}).Where(&RepositoryRef{GroupID: p.ID}).Scan(&out.Repos).Error; err != nil {
		panic(err)
	}

	c.JSON(200, &out)
}

func (m *manifestHandler) appendManifestItems(c *gin.Context) {
	items := []RepositoryRef{}
	if err := c.ShouldBindJSON(&items); err != nil {
		panic(err)
	}
	num := len(items)
	if num == 0 {
		c.AbortWithStatus(200)
		return
	}
	p := pathParams{}
	if err := c.ShouldBindUri(&p); err != nil {
		panic(err)
	}

	// range 是拷贝
	for i := 0; i < num; i++ {
		items[i].GroupID = p.ID
		if err := m.validator.Struct(&items[i]); err != nil {
			panic(err)
		}
	}

	err := m.ds.GromDB().Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	c.AbortWithStatus(200)
}

func (m *manifestHandler) removeManifestItems(c *gin.Context) {
	repoIds := []types.ID{}
	if err := c.ShouldBindJSON(&repoIds); err != nil {
		panic(err)
	}
	if len(repoIds) == 0 {
		c.AbortWithStatus(204)
		return
	}
	p := pathParams{}
	if err := c.ShouldBindUri(&p); err != nil {
		panic(err)
	}

	err := m.ds.GromDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&RepositoryRef{GroupID: p.ID}).Where("id IN (?)", repoIds).Delete(&RepositoryRef{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	c.AbortWithStatus(204)
}
