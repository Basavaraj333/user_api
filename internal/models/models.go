package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1"`
	Dob  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1"`
	Dob  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type UserWithAge struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}

type PaginatedResponse struct {
	Data       []UserWithAge `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
