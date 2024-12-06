package main

import (
	"log"

	"github.com/Roixys/e-fast-store-api/config"
	"github.com/Roixys/e-fast-store-api/model"
	"gorm.io/gorm"
)

func seedCategories(db *gorm.DB) {
	categories := []model.Category{
		{
			Name:        "Electronics",
			Description: stringPtr("Gadgets and electronic devices"),
		},
		{
			Name:        "Clothing",
			Description: stringPtr("Apparel and fashion items"),
		},
		{
			Name:        "Home & Kitchen",
			Description: stringPtr("Furniture, appliances, and home decor"),
		},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Fatalf("Could not seed category: %v", err)
		}
	}
}

func seedProducts(db *gorm.DB) {
	var categories []model.Category
	if err := db.Find(&categories).Error; err != nil {
		log.Fatalf("Could not retrieve categories for seeding products: %v", err)
	}

	products := []model.Product{
		{
			Name:       "Smartphone",
			Price:      500000,
			CategoryID: categories[0].ID,
		},
		{
			Name:       "Jeans",
			Price:      150000,
			CategoryID: categories[1].ID,
		},
		{
			Name:       "Blender",
			Price:      300000,
			CategoryID: categories[2].ID,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			log.Fatalf("Could not seed product: %v", err)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	configuration := config.LoadConfig(".")

	db := config.NewPostgres(configuration.DBSource, configuration.Environment)

	// Seed the database
	seedCategories(db)
	seedProducts(db)

	log.Println("Seeding completed successfully")
}
