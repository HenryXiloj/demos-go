package domain

import "context"

// UserRepo defines the contract for user-related data operations.
// Implementations (e.g., MySQL repository) handle persistence for User entities.
type UserRepo interface {
	// Create inserts a new user record into the database.
	// Returns the generated user ID or an error if the operation fails.
	Create(ctx context.Context, u *User) (int64, error)
}

// CompanyRepo defines the contract for company-related data operations.
// Implementations (e.g., PostgreSQL repository) handle persistence for Company entities.
type CompanyRepo interface {
	// Create inserts a new company record into the database.
	// Returns the generated company ID or an error if the operation fails.
	Create(ctx context.Context, c *Company) (int64, error)
}

// BrandRepo defines the contract for brand-related data operations.
// Implementations (e.g., Oracle repository) handle persistence for Brand entities.
type BrandRepo interface {
	// Create inserts a new brand record into the database.
	// Returns the generated brand ID or an error if the operation fails.
	Create(ctx context.Context, b *Brand) (int64, error)
}
