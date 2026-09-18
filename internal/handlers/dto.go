package handlers

type CreateListingDTO struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	City        string `json:"city" validate:"required,min=2,max=100"`
}
