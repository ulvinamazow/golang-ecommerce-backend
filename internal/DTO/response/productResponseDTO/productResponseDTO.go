package productresponsedto

type ProductResponseDTO struct {
	Name        string  `json:"name"`
	CreatedBy   int64   `json:"created_by"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"category_id"`
}
