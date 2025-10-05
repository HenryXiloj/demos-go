package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

// ---- Service interfaces (what handlers depend on) ----

type UserService interface {
	CreateUser(ctx context.Context, name, lastName string) (int64, error)
}

type CompanyService interface {
	CreateCompany(ctx context.Context, name string) (int64, error)
}

type BrandService interface {
	CreateBrand(ctx context.Context, name string) (int64, error)
}

// ---- Concrete service implementations ----

// common small helper to enforce timeouts per call
func withTimeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		// no timeout requested; still allow cancellation propagation
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, d)
}

// ---------- User ----------

type userService struct {
	repo    UserRepo
	timeout time.Duration
}

func NewUserService(repo UserRepo, timeout time.Duration) UserService {
	return &userService{repo: repo, timeout: timeout}
}

func (s *userService) CreateUser(ctx context.Context, name, lastName string) (int64, error) {
	name = strings.TrimSpace(name)
	lastName = strings.TrimSpace(lastName)
	if name == "" || lastName == "" {
		return 0, errors.New("name and lastName are required")
	}

	u := &User{Name: name, LastName: lastName}

	cctx, cancel := withTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Create(cctx, u)
}

// ---------- Company ----------

type companyService struct {
	repo    CompanyRepo
	timeout time.Duration
}

func NewCompanyService(repo CompanyRepo, timeout time.Duration) CompanyService {
	return &companyService{repo: repo, timeout: timeout}
}

func (s *companyService) CreateCompany(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("company name is required")
	}

	c := &Company{Name: name}

	cctx, cancel := withTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Create(cctx, c)
}

// ---------- Brand ----------

type brandService struct {
	repo    BrandRepo
	timeout time.Duration
}

func NewBrandService(repo BrandRepo, timeout time.Duration) BrandService {
	return &brandService{repo: repo, timeout: timeout}
}

func (s *brandService) CreateBrand(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("brand name is required")
	}

	b := &Brand{Name: name}

	cctx, cancel := withTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Create(cctx, b)
}
