package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

// =====================================================
// Service Interfaces
// =====================================================

// UserService defines business operations related to users.
// Handlers depend on this interface rather than concrete repositories.
type UserService interface {
	// CreateUser validates and creates a new user record.
	// Returns the created user ID or an error.
	CreateUser(ctx context.Context, name, lastName string) (int64, error)
}

// CompanyService defines business operations related to companies.
type CompanyService interface {
	// CreateCompany validates and creates a new company record.
	CreateCompany(ctx context.Context, name string) (int64, error)
}

// BrandService defines business operations related to brands.
type BrandService interface {
	// CreateBrand validates and creates a new brand record.
	CreateBrand(ctx context.Context, name string) (int64, error)
}

// =====================================================
// Helper: Context Timeout
// =====================================================

// withTimeout creates a child context with the specified timeout duration.
// If `d <= 0`, it still returns a cancelable context that inherits from parent.
func withTimeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		// No timeout configured, but still allow manual cancellation.
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, d)
}

// =====================================================
// User Service Implementation
// =====================================================

// userService provides user-related business logic and enforces validation and timeouts.
type userService struct {
	repo    UserRepo      // Underlying data repository for users
	timeout time.Duration // Operation timeout duration
}

// NewUserService creates a new instance of UserService with the given repository and timeout.
func NewUserService(repo UserRepo, timeout time.Duration) UserService {
	return &userService{repo: repo, timeout: timeout}
}

// CreateUser validates input and delegates user creation to the repository layer.
func (s *userService) CreateUser(ctx context.Context, name, lastName string) (int64, error) {
	// Clean and validate input
	name = strings.TrimSpace(name)
	lastName = strings.TrimSpace(lastName)
	if name == "" || lastName == "" {
		return 0, errors.New("name and lastName are required")
	}

	// Build user entity
	u := &User{Name: name, LastName: lastName}

	// Apply timeout for database operation
	cctx, cancel := withTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Create(cctx, u)
}

// =====================================================
// Company Service Implementation
// =====================================================

// companyService provides company-related business logic.
type companyService struct {
	repo    CompanyRepo
	timeout time.Duration
}

// NewCompanyService creates a new instance of CompanyService with timeout.
func NewCompanyService(repo CompanyRepo, timeout time.Duration) CompanyService {
	return &companyService{repo: repo, timeout: timeout}
}

// CreateCompany validates the company name and creates a record via the repository.
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

// =====================================================
// Brand Service Implementation
// =====================================================

// brandService provides brand-related business logic.
type brandService struct {
	repo    BrandRepo
	timeout time.Duration
}

// NewBrandService creates a new instance of BrandService with timeout.
func NewBrandService(repo BrandRepo, timeout time.Duration) BrandService {
	return &brandService{repo: repo, timeout: timeout}
}

// CreateBrand validates the brand name and delegates creation to the repository.
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
