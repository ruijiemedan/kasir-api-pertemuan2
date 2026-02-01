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

// Config adalah tempat menyimpan setting aplikasi
// Seperti buku pengaturan restoran
type Config struct {
	Port   string `mapstructure:"PORT"`    // Port berapa server jalan
	DBConn string `mapstructure:"DB_CONN"` // Alamat database
}

func main() {
	// ===== BAGIAN 1: BACA SETTING (CONFIG) =====
	// Viper adalah alat untuk baca file .env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Cek apakah ada file .env
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// Ambil setting dari .env
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Println("📋 Config loaded:")
	fmt.Println("   Port:", config.Port)
	fmt.Println("   DB Connected")

	// ===== BAGIAN 2: HUBUNGKAN KE DATABASE =====
	// Buka pintu ke gudang data (database)
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("❌ Failed to initialize database:", err)
	}
	defer db.Close() // Tutup koneksi saat program selesai

	// ===== BAGIAN 3: DEPENDENCY INJECTION =====
	// Ini seperti Manager yang kenalkan semua staff

	// 1. Buat Anak Gudang (Repository)
	productRepo := repositories.NewProductRepository(db)
	fmt.Println("✅ Repository created (Anak Gudang siap)")

	// 2. Buat Koki (Service) dan kenalkan dengan Anak Gudang
	productService := services.NewProductService(productRepo)
	fmt.Println("✅ Service created (Koki siap)")

	// 3. Buat Pelayan (Handler) dan kenalkan dengan Koki
	productHandler := handlers.NewProductHandler(productService)
	fmt.Println("✅ Handler created (Pelayan siap)")

	// ===== BAGIAN 4: SETUP ROUTES (JALUR PESANAN) =====
	// Ini seperti papan menu di restoran
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	fmt.Println("✅ Routes configured")

	// ===== BAGIAN 5: JALANKAN SERVER =====
	// Buka restoran untuk pelanggan!
	addr := "0.0.0.0:" + config.Port
	fmt.Println("\n🚀 Server running di", addr)
	fmt.Println("📝 Endpoints:")
	fmt.Println("   GET    /api/produk      → Lihat semua produk")
	fmt.Println("   POST   /api/produk      → Tambah produk baru")
	fmt.Println("   GET    /api/produk/{id} → Lihat satu produk")
	fmt.Println("   PUT    /api/produk/{id} → Update produk")
	fmt.Println("   DELETE /api/produk/{id} → Hapus produk")
	fmt.Println()

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("❌ Gagal running server:", err)
	}
}
