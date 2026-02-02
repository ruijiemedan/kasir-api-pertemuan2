package main

import (
	"fmt"
	"kasir-api-pertemuan2/database"
	"kasir-api-pertemuan2/handlers"
	"kasir-api-pertemuan2/repositories"
	"kasir-api-pertemuan2/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Load config dari .env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Println("📋 Config loaded:")
	fmt.Println("   Port:", config.Port)
	fmt.Println("   DB Connected")

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Dependency Injection - Product
	productRepo := repositories.NewProductRepository(db)
	fmt.Println("✅ Product Repository created")

	productService := services.NewProductService(productRepo)
	fmt.Println("✅ Product Service created")

	productHandler := handlers.NewProductHandler(productService)
	fmt.Println("✅ Product Handler created")

	// Dependency Injection - Category
	categoryRepo := repositories.NewCategoryRepository(db)
	fmt.Println("✅ Category Repository created")

	categoryService := services.NewCategoryService(categoryRepo)
	fmt.Println("✅ Category Service created")

	categoryHandler := handlers.NewCategoryHandler(categoryService)
	fmt.Println("✅ Category Handler created")

	// Setup routes - Product
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Setup routes - Category
	http.HandleFunc("/api/kategori", categoryHandler.HandleCategories)
	http.HandleFunc("/api/kategori/", categoryHandler.HandleCategoryByID)

	fmt.Println("✅ Routes configured")

	// Start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("🚀 Server running di", addr)
	fmt.Println("📝 Endpoints:")
	fmt.Println("   Product Endpoints:")
	fmt.Println("   GET    /api/produk      → Lihat semua produk")
	fmt.Println("   POST   /api/produk      → Tambah produk baru")
	fmt.Println("   GET    /api/produk/{id} → Lihat satu produk")
	fmt.Println("   PUT    /api/produk/{id} → Update produk")
	fmt.Println("   DELETE /api/produk/{id} → Hapus produk")
	fmt.Println()
	fmt.Println("   Category Endpoints:")
	fmt.Println("   GET    /api/kategori      → Lihat semua kategori")
	fmt.Println("   POST   /api/kategori      → Tambah kategori baru")
	fmt.Println("   GET    /api/kategori/{id} → Lihat satu kategori")
	fmt.Println("   PUT    /api/kategori/{id} → Update kategori")
	fmt.Println("   DELETE /api/kategori/{id} → Hapus kategori")
	fmt.Println()

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
