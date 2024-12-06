package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
}

type User struct {
	Base
	Username string `gorm:"not null; uniqueIndex"`
	Email    string `gorm:"not null; uniqueIndex"`
	Password string `gorm:"not null"`
	Balance  uint64 `gorm:"not null"`
}

type Category struct {
	Base
	Name        string  `gorm:"not null"`
	Description *string `gorm:"not null"`
}

type Product struct {
	Base
	Name       string `gorm:"not null"`
	Price      uint64 `gorm:"not null"`
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
}

type Cart struct {
	Base
	Quantity  int       `gorm:"quantity"`
	UserID    uuid.UUID `gorm:"index:idx_name,unique; not null"`
	ProductID uuid.UUID `gorm:"index:idx_name,unique; not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Product   Product   `gorm:"foreignKey:ProductID"`
}

type PaymentEvent struct {
	Base
	UserID      uuid.UUID `gorm:"user_id"`
	TotalAmount uint64    `gorm:"total_amount"`
}

type PaymentOrder struct {
	Base
	Amount    float64   `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	UserID    uuid.UUID `gorm:"index:idx_name,unique; not null"`
	ProductID uuid.UUID `gorm:"index:idx_name,unique; not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Product   Product   `gorm:"foreignKey:ProductID"`
}

func (User) TableName() string {
	return "users"
}

func (Product) TableName() string {
	return "products"
}

func (Category) TableName() string {
	return "categories"
}

func (Cart) TableName() string {
	return "carts"
}

func (PaymentEvent) TableName() string {
	return "payment_events"
}

func (PaymentOrder) TableName() string {
	return "payment_orders"
}
