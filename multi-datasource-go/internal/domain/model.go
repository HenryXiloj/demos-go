package domain

// User represents an application user entity.
// It maps to the "users" table in the MySQL database.
type User struct {
	ID       int64  `json:"id"`       // Unique identifier for the user (auto-incremented primary key)
	Name     string `json:"name"`     // First name of the user
	LastName string `json:"lastName"` // Last name of the user
}

// Company represents a company entity.
// It maps to the "companies" table in the PostgreSQL database.
type Company struct {
	ID   int64  `json:"id"`   // Unique identifier for the company (auto-incremented primary key)
	Name string `json:"name"` // Name of the company
}

// Brand represents a brand entity.
// It maps to the "brands" table in the Oracle database.
type Brand struct {
	ID   int64  `json:"id"`   // Unique identifier for the brand (auto-incremented primary key)
	Name string `json:"name"` // Name of the brand
}
