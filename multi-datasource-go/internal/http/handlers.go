package http

import (
	"context"
	"net/http"
	"time"

	"multi-datasource-go/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Users     domain.UserRepo
	Companies domain.CompanyRepo
	Brands    domain.BrandRepo
	Timeout   time.Duration
}

func (h *Handlers) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.POST("/users", h.createUser)

	v2 := r.Group("/api/v2")
	v2.POST("/companies", h.createCompany)

	v3 := r.Group("/api/v3")
	v3.POST("/brands", h.createBrand)
}

func (h *Handlers) ctx(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), h.Timeout)
}

func (h *Handlers) createUser(c *gin.Context) {
	var u domain.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := h.ctx(c)
	defer cancel()
	id, err := h.Users.Create(ctx, &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handlers) createCompany(c *gin.Context) {
	var m domain.Company
	if err := c.BindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := h.ctx(c)
	defer cancel()
	id, err := h.Companies.Create(ctx, &m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handlers) createBrand(c *gin.Context) {
	var b domain.Brand
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := h.ctx(c)
	defer cancel()
	id, err := h.Brands.Create(ctx, &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
