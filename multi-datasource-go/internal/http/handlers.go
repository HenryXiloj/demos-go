package http

import (
	"context"
	"net/http"
	"time"

	"multi-datasource-go/internal/domain"

	"github.com/gin-gonic/gin"
)

// Handlers groups all HTTP handler dependencies, including
// repositories for users, companies, and brands.
// It also holds a configurable request timeout used for
// per-request context control.
type Handlers struct {
	Users     domain.UserRepo    // Repository for MySQL-backed user operations
	Companies domain.CompanyRepo // Repository for PostgreSQL-backed company operations
	Brands    domain.BrandRepo   // Repository for Oracle-backed brand operations
	Timeout   time.Duration      // Timeout duration applied to each incoming request
}

// Register registers all versioned HTTP routes handled by this service.
// It organizes endpoints under /api/v1, /api/v2, and /api/v3 prefixes
// to reflect the data source each route interacts with.
func (h *Handlers) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.POST("/users", h.createUser)

	v2 := r.Group("/api/v2")
	v2.POST("/companies", h.createCompany)

	v3 := r.Group("/api/v3")
	v3.POST("/brands", h.createBrand)
}

// ctx creates a derived context with the configured timeout.
// Each request will automatically cancel operations after the
// timeout expires or if the client disconnects.
func (h *Handlers) ctx(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), h.Timeout)
}

// createUser handles POST /api/v1/users requests.
// It binds the request body to a domain.User, validates it,
// and calls the MySQL repository to persist the record.
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

// createCompany handles POST /api/v2/companies requests.
// It binds incoming JSON to a domain.Company object and uses
// the PostgreSQL repository to insert a new record.
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

// createBrand handles POST /api/v3/brands requests.
// It binds the JSON payload to a domain.Brand and
// calls the Oracle repository to persist the data.
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
