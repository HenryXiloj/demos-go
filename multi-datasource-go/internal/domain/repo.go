package domain

import "context"

type UserRepo interface {
	Create(ctx context.Context, u *User) (int64, error)
}

type CompanyRepo interface {
	Create(ctx context.Context, c *Company) (int64, error)
}

type BrandRepo interface {
	Create(ctx context.Context, b *Brand) (int64, error)
}
