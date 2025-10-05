package domain

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
}

type Company struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Brand struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
