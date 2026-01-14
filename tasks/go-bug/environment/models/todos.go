package models

// Todo represents a todo item
type Todo struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title"`
	Completed  bool      `json:"completed"`
	CategoryID *uint     `json:"category_id,omitempty"`
	Category   *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// Category represents a todo category
type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
