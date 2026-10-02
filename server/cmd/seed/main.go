package main

import (
	"log"
	"time"

	"layanesia-server/config"
	dbmigrate "layanesia-server/db"
	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// Seeder data contoh untuk pengembangan/staging. JANGAN dijalankan di production.
func main() {
	log.Println("🌱 Starting Layanesia Database Seeder (GORM)...")

	cfg := config.MustLoadConfig()
	if cfg.IsProduction() {
		log.Fatal("❌ Seeder tidak boleh dijalankan dengan ENV=production.")
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	if err := dbmigrate.MigrateUp(cfg.DSN()); err != nil {
		log.Fatalf("❌ migrasi gagal: %v", err)
	}

	// 2. Seed Master Data Locations
	log.Println("📍 Seeding Master Data Locations...")
	locations := []models.Location{
		// ==========================================
		// 1. DKI JAKARTA (44 KECAMATAN - 100%)
		// ==========================================
		// Jakarta Selatan (10 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Kebayoran Baru", PostalCode: stringPtr("12110")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Kebayoran Lama", PostalCode: stringPtr("12210")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Pesanggrahan", PostalCode: stringPtr("12320")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Cilandak", PostalCode: stringPtr("12430")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Pasar Minggu", PostalCode: stringPtr("12520")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Jagakarsa", PostalCode: stringPtr("12620")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Mampang Prapatan", PostalCode: stringPtr("12790")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Pancoran", PostalCode: stringPtr("12780")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Tebet", PostalCode: stringPtr("12810")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Selatan", District: "Setiabudi", PostalCode: stringPtr("12910")},

		// Jakarta Pusat (8 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Gambir", PostalCode: stringPtr("10110")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Tanah Abang", PostalCode: stringPtr("10210")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Menteng", PostalCode: stringPtr("10310")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Senen", PostalCode: stringPtr("10410")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Cempaka Putih", PostalCode: stringPtr("10510")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Johar Baru", PostalCode: stringPtr("10560")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Kemayoran", PostalCode: stringPtr("10610")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Pusat", District: "Sawah Besar", PostalCode: stringPtr("10710")},

		// Jakarta Barat (8 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Cengkareng", PostalCode: stringPtr("11730")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Grogol Petamburan", PostalCode: stringPtr("11470")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Kalideres", PostalCode: stringPtr("11840")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Kebon Jeruk", PostalCode: stringPtr("11530")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Kembangan", PostalCode: stringPtr("11610")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Palmerah", PostalCode: stringPtr("11480")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Taman Sari", PostalCode: stringPtr("11110")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Barat", District: "Tambora", PostalCode: stringPtr("11210")},

		// Jakarta Timur (10 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Matraman", PostalCode: stringPtr("13110")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Pulo Gadung", PostalCode: stringPtr("13260")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Jatinegara", PostalCode: stringPtr("13310")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Duren Sawit", PostalCode: stringPtr("13440")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Kramat Jati", PostalCode: stringPtr("13510")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Makasar", PostalCode: stringPtr("13570")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Pasar Rebo", PostalCode: stringPtr("13710")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Ciracas", PostalCode: stringPtr("13740")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Cipayung", PostalCode: stringPtr("13840")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Timur", District: "Cakung", PostalCode: stringPtr("13910")},

		// Jakarta Utara (6 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Cilincing", PostalCode: stringPtr("14120")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Kelapa Gading", PostalCode: stringPtr("14240")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Koja", PostalCode: stringPtr("14210")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Pademangan", PostalCode: stringPtr("14410")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Penjaringan", PostalCode: stringPtr("14440")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Jakarta Utara", District: "Tanjung Priok", PostalCode: stringPtr("14310")},

		// Kepulauan Seribu (2 Kecamatan)
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Kepulauan Seribu", District: "Kepulauan Seribu Selatan", PostalCode: stringPtr("14510")},
		{Country: "Indonesia", Province: "DKI Jakarta", CityOrRegency: "Kepulauan Seribu", District: "Kepulauan Seribu Utara", PostalCode: stringPtr("14520")},

		// ==========================================
		// 2. BOGOR (KOTA: 6, KABUPATEN: 40 - 100%)
		// ==========================================
		// Kota Bogor (6 Kecamatan)
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Bogor Barat", PostalCode: stringPtr("16111")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Bogor Selatan", PostalCode: stringPtr("16132")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Bogor Tengah", PostalCode: stringPtr("16121")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Bogor Timur", PostalCode: stringPtr("16143")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Bogor Utara", PostalCode: stringPtr("16152")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bogor", District: "Tanah Sareal", PostalCode: stringPtr("16161")},

		// Kabupaten Bogor (40 Kecamatan)
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Babakan Madang", PostalCode: stringPtr("16810")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Bojonggede", PostalCode: stringPtr("16922")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Caringin", PostalCode: stringPtr("16730")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cariu", PostalCode: stringPtr("16840")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Ciampea", PostalCode: stringPtr("16620")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Ciawi", PostalCode: stringPtr("16720")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cibinong", PostalCode: stringPtr("16911")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cibungbulang", PostalCode: stringPtr("16630")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cigombong", PostalCode: stringPtr("16740")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cigudeg", PostalCode: stringPtr("16660")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cijeruk", PostalCode: stringPtr("16740")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cileungsi", PostalCode: stringPtr("16820")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Ciomas", PostalCode: stringPtr("16610")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Cisarua", PostalCode: stringPtr("16750")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Ciseeng", PostalCode: stringPtr("16120")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Citeureup", PostalCode: stringPtr("16810")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Dramaga", PostalCode: stringPtr("16680")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Gunung Putri", PostalCode: stringPtr("16961")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Gunung Sindur", PostalCode: stringPtr("16340")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Jasinga", PostalCode: stringPtr("16670")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Jonggol", PostalCode: stringPtr("16830")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Kemang", PostalCode: stringPtr("16310")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Klapanunggal", PostalCode: stringPtr("16710")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Leuwiliang", PostalCode: stringPtr("16640")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Leuwisadeng", PostalCode: stringPtr("16650")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Megamendung", PostalCode: stringPtr("16770")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Nanggung", PostalCode: stringPtr("16650")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Pamijahan", PostalCode: stringPtr("16630")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Parung", PostalCode: stringPtr("16330")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Parung Panjang", PostalCode: stringPtr("16360")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Ranca Bungur", PostalCode: stringPtr("16310")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Rumpin", PostalCode: stringPtr("16350")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Sukajaya", PostalCode: stringPtr("16660")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Sukamakmur", PostalCode: stringPtr("16830")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Sukaraja", PostalCode: stringPtr("16710")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Tajurhalang", PostalCode: stringPtr("16320")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Tamansari", PostalCode: stringPtr("16610")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Tanjungsari", PostalCode: stringPtr("16840")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Tenjo", PostalCode: stringPtr("16370")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bogor", District: "Tenjolaya", PostalCode: stringPtr("16620")},

		// ==========================================
		// 3. DEPOK (11 KECAMATAN - 100%)
		// ==========================================
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Beji", PostalCode: stringPtr("16421")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Bojongsari", PostalCode: stringPtr("16516")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Cilodong", PostalCode: stringPtr("16413")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Cimanggis", PostalCode: stringPtr("16451")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Cinere", PostalCode: stringPtr("16514")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Cipayung", PostalCode: stringPtr("16437")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Limo", PostalCode: stringPtr("16515")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Pancoran Mas", PostalCode: stringPtr("16431")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Sawangan", PostalCode: stringPtr("16511")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Sukmajaya", PostalCode: stringPtr("16412")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Depok", District: "Tapos", PostalCode: stringPtr("16457")},

		// ==========================================
		// 4. TANGERANG (KOTA: 13, TANGSEL: 7, KAB: 29 - 100%)
		// ==========================================
		// Kota Tangerang (13 Kecamatan)
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Batuceper", PostalCode: stringPtr("15122")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Benda", PostalCode: stringPtr("15121")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Cibodas", PostalCode: stringPtr("15138")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Ciledug", PostalCode: stringPtr("15153")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Cipondoh", PostalCode: stringPtr("15148")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Jatiuwung", PostalCode: stringPtr("15134")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Karangtengah", PostalCode: stringPtr("15157")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Karawaci", PostalCode: stringPtr("15115")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Larangan", PostalCode: stringPtr("15154")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Neglasari", PostalCode: stringPtr("15129")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Periuk", PostalCode: stringPtr("15131")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Pinang", PostalCode: stringPtr("15145")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang", District: "Tangerang", PostalCode: stringPtr("15111")},

		// Kota Tangerang Selatan (7 Kecamatan)
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Ciputat", PostalCode: stringPtr("15411")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Ciputat Timur", PostalCode: stringPtr("15419")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Pamulang", PostalCode: stringPtr("15417")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Pondok Aren", PostalCode: stringPtr("15224")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Serpong", PostalCode: stringPtr("15310")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Serpong Utara", PostalCode: stringPtr("15326")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Tangerang Selatan", District: "Setu", PostalCode: stringPtr("15314")},

		// Kabupaten Tangerang (29 Kecamatan)
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Balaraja", PostalCode: stringPtr("15610")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Cikupa", PostalCode: stringPtr("15710")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Cisauk", PostalCode: stringPtr("15341")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Cisoka", PostalCode: stringPtr("15730")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Curug", PostalCode: stringPtr("15810")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Gunungkaler", PostalCode: stringPtr("15620")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Jambe", PostalCode: stringPtr("15720")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Jayanti", PostalCode: stringPtr("15610")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Kelapa Dua", PostalCode: stringPtr("15810")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Kemiri", PostalCode: stringPtr("15530")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Kosambi", PostalCode: stringPtr("15211")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Kresek", PostalCode: stringPtr("15620")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Kronjo", PostalCode: stringPtr("15550")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Legok", PostalCode: stringPtr("15820")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Mauk", PostalCode: stringPtr("15530")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Mekarbaru", PostalCode: stringPtr("15550")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Pagedangan", PostalCode: stringPtr("15339")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Pakuhaji", PostalCode: stringPtr("15570")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Panongan", PostalCode: stringPtr("15711")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Pasar Kemis", PostalCode: stringPtr("15560")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Rajeg", PostalCode: stringPtr("15540")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Sepatan", PostalCode: stringPtr("15520")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Sepatan Timur", PostalCode: stringPtr("15521")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Sindang Jaya", PostalCode: stringPtr("15561")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Solear", PostalCode: stringPtr("15731")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Sukadiri", PostalCode: stringPtr("15531")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Sukamulya", PostalCode: stringPtr("15611")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Teluknaga", PostalCode: stringPtr("15510")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Tangerang", District: "Tigaraksa", PostalCode: stringPtr("15720")},

		// ==========================================
		// 5. BEKASI (KOTA: 12, KABUPATEN: 23 - 100%)
		// ==========================================
		// Kota Bekasi (12 Kecamatan)
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Bantar Gebang", PostalCode: stringPtr("17151")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Bekasi Barat", PostalCode: stringPtr("17145")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Bekasi Selatan", PostalCode: stringPtr("17141")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Bekasi Timur", PostalCode: stringPtr("17111")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Bekasi Utara", PostalCode: stringPtr("17121")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Jatiasih", PostalCode: stringPtr("17423")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Jatisampurna", PostalCode: stringPtr("17433")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Medan Satria", PostalCode: stringPtr("17132")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Mustikajaya", PostalCode: stringPtr("17158")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Pondok Gede", PostalCode: stringPtr("17411")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Pondok Melati", PostalCode: stringPtr("17415")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bekasi", District: "Rawalumbu", PostalCode: stringPtr("17116")},

		// Kabupaten Bekasi (23 Kecamatan)
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Babelan", PostalCode: stringPtr("17610")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Bojongmangu", PostalCode: stringPtr("17356")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cabangbingin", PostalCode: stringPtr("17720")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cibarusah", PostalCode: stringPtr("17340")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cibitung", PostalCode: stringPtr("17520")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cikarang Barat", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cikarang Pusat", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cikarang Selatan", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cikarang Timur", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Cikarang Utara", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Karangbahagia", PostalCode: stringPtr("17535")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Kedungwaringin", PostalCode: stringPtr("17540")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Muara Gembong", PostalCode: stringPtr("17730")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Pebayuran", PostalCode: stringPtr("17710")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Serang Baru", PostalCode: stringPtr("17330")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Setu", PostalCode: stringPtr("17320")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Sukakarya", PostalCode: stringPtr("17530")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Sukatani", PostalCode: stringPtr("17630")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Sukawangi", PostalCode: stringPtr("17620")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Tambelang", PostalCode: stringPtr("17620")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Tambun Selatan", PostalCode: stringPtr("17510")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Tambun Utara", PostalCode: stringPtr("17510")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Bekasi", District: "Tarumajaya", PostalCode: stringPtr("17210")},

		// ==========================================
		// 6. PROVINSI & KOTA UTAMA LAINNYA
		// ==========================================
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bandung", District: "Coblong", PostalCode: stringPtr("40132")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bandung", District: "Cicendo", PostalCode: stringPtr("40171")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bandung", District: "Sumur Bandung", PostalCode: stringPtr("40111")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bandung", District: "Lengkong", PostalCode: stringPtr("40261")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Bandung", District: "Sukajadi", PostalCode: stringPtr("40161")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Cimahi", District: "Cimahi Tengah", PostalCode: stringPtr("40522")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Sukabumi", District: "Cikole", PostalCode: stringPtr("43111")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kota Tasikmalaya", District: "Cihideung", PostalCode: stringPtr("46122")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Karawang", District: "Karawang Barat", PostalCode: stringPtr("41311")},
		{Country: "Indonesia", Province: "Jawa Barat", CityOrRegency: "Kabupaten Karawang", District: "Karawang Timur", PostalCode: stringPtr("41314")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Serang", District: "Serang", PostalCode: stringPtr("42111")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kota Cilegon", District: "Cibeber", PostalCode: stringPtr("42422")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Lebak", District: "Rangkasbitung", PostalCode: stringPtr("42311")},
		{Country: "Indonesia", Province: "Banten", CityOrRegency: "Kabupaten Pandeglang", District: "Pandeglang", PostalCode: stringPtr("42211")},
		{Country: "Indonesia", Province: "Jawa Tengah", CityOrRegency: "Kota Semarang", District: "Semarang Tengah", PostalCode: stringPtr("50131")},
		{Country: "Indonesia", Province: "Jawa Tengah", CityOrRegency: "Kota Surakarta (Solo)", District: "Banjarsari", PostalCode: stringPtr("57131")},
		{Country: "Indonesia", Province: "Jawa Tengah", CityOrRegency: "Kota Magelang", District: "Magelang Tengah", PostalCode: stringPtr("56111")},
		{Country: "Indonesia", Province: "Jawa Tengah", CityOrRegency: "Kabupaten Banyumas", District: "Purwokerto Timur", PostalCode: stringPtr("53111")},
		{Country: "Indonesia", Province: "DI Yogyakarta", CityOrRegency: "Kota Yogyakarta", District: "Gondomanan", PostalCode: stringPtr("55121")},
		{Country: "Indonesia", Province: "DI Yogyakarta", CityOrRegency: "Kabupaten Sleman", District: "Depok", PostalCode: stringPtr("55281")},
		{Country: "Indonesia", Province: "DI Yogyakarta", CityOrRegency: "Kabupaten Bantul", District: "Sewon", PostalCode: stringPtr("55188")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kota Surabaya", District: "Tegalsari", PostalCode: stringPtr("60261")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kota Surabaya", District: "Gubeng", PostalCode: stringPtr("60281")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kota Malang", District: "Lowokwaru", PostalCode: stringPtr("65141")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kota Batu", District: "Batu", PostalCode: stringPtr("65311")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kabupaten Sidoarjo", District: "Sidoarjo", PostalCode: stringPtr("61211")},
		{Country: "Indonesia", Province: "Jawa Timur", CityOrRegency: "Kabupaten Gresik", District: "Gresik", PostalCode: stringPtr("61111")},
		{Country: "Indonesia", Province: "Bali", CityOrRegency: "Kota Denpasar", District: "Denpasar Selatan", PostalCode: stringPtr("80221")},
		{Country: "Indonesia", Province: "Bali", CityOrRegency: "Kabupaten Badung", District: "Kuta", PostalCode: stringPtr("80361")},
		{Country: "Indonesia", Province: "Bali", CityOrRegency: "Kabupaten Gianyar", District: "Ubud", PostalCode: stringPtr("80571")},
		{Country: "Indonesia", Province: "Sumatera Utara", CityOrRegency: "Kota Medan", District: "Medan Kota", PostalCode: stringPtr("20211")},
		{Country: "Indonesia", Province: "Riau", CityOrRegency: "Kota Pekanbaru", District: "Tampan", PostalCode: stringPtr("28291")},
		{Country: "Indonesia", Province: "Kepulauan Riau", CityOrRegency: "Kota Batam", District: "Batam Kota", PostalCode: stringPtr("29432")},
		{Country: "Indonesia", Province: "Sumatera Barat", CityOrRegency: "Kota Padang", District: "Padang Barat", PostalCode: stringPtr("25111")},
		{Country: "Indonesia", Province: "Sumatera Selatan", CityOrRegency: "Kota Palembang", District: "Ilir Barat I", PostalCode: stringPtr("30139")},
		{Country: "Indonesia", Province: "Lampung", CityOrRegency: "Kota Bandar Lampung", District: "Tanjung Karang Pusat", PostalCode: stringPtr("35111")},
		{Country: "Indonesia", Province: "Kalimantan Timur", CityOrRegency: "Kota Balikpapan", District: "Balikpapan Kota", PostalCode: stringPtr("76111")},
		{Country: "Indonesia", Province: "Kalimantan Timur", CityOrRegency: "Kota Samarinda", District: "Samarinda Kota", PostalCode: stringPtr("75121")},
		{Country: "Indonesia", Province: "Kalimantan Barat", CityOrRegency: "Kota Pontianak", District: "Pontianak Selatan", PostalCode: stringPtr("78121")},
		{Country: "Indonesia", Province: "Sulawesi Selatan", CityOrRegency: "Kota Makassar", District: "Ujung Pandang", PostalCode: stringPtr("90111")},
		{Country: "Indonesia", Province: "Sulawesi Utara", CityOrRegency: "Kota Manado", District: "Wenang", PostalCode: stringPtr("95111")},
		{Country: "Indonesia", Province: "Nusa Tenggara Barat", CityOrRegency: "Kota Mataram", District: "Mataram", PostalCode: stringPtr("83121")},
		{Country: "Indonesia", Province: "Nusa Tenggara Timur", CityOrRegency: "Kabupaten Manggarai Barat", District: "Komodo", PostalCode: stringPtr("86554")},
		{Country: "Indonesia", Province: "Papua", CityOrRegency: "Kota Jayapura", District: "Jayapura Utara", PostalCode: stringPtr("99111")},
	}

	for _, loc := range locations {
		db.Where("province = ? AND city_or_regency = ? AND district = ?", loc.Province, loc.CityOrRegency, loc.District).FirstOrCreate(&loc)
	}
	log.Println("✅ Locations seeded successfully.")

	// 3. Seed Users with PLAIN TEXT passwords (GORM BeforeCreate hook will hash automatically)
	log.Println("👥 Seeding Sample Users (Plaintext password auto-hashed by GORM)...")

	userWorker1 := models.User{
		ID:            uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		Email:         "budi.worker@gmail.com",
		Phone:         "081234567890",
		Name:          "Budi Santoso",
		Password:      "Password123!", // Plain text password -> auto-hashed by BeforeCreate hook
		ActiveMode:    models.ModeSeeker,
		IsKTPVerified: true,
		KTPNumber:     stringPtr("3174012304900001"),
		KTPImageURL:   stringPtr("https://res.cloudinary.com/layanesia/image/upload/v123/ktp_budi.jpg"),
		ResumeURL:     stringPtr("https://res.cloudinary.com/layanesia/image/upload/v123/cv_budi.pdf"),
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Kebayoran Baru",
		AddressDetail: stringPtr("Jl. Wijaya I No. 42, Kebayoran Baru"),
		Latitude:      floatPtr(-6.2435),
		Longitude:     floatPtr(106.8021),
		OSMPlaceID:    stringPtr("osm_rel_12345"),
	}

	userWorker2 := models.User{
		ID:            uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"),
		Email:         "siti.art@gmail.com",
		Phone:         "081987654321",
		Name:          "Siti Rahmawati",
		Password:      "Password123!", // Plain text password -> auto-hashed by BeforeCreate hook
		ActiveMode:    models.ModeSeeker,
		IsKTPVerified: true,
		KTPNumber:     stringPtr("3174025508920003"),
		KTPImageURL:   stringPtr("https://res.cloudinary.com/layanesia/image/upload/v123/ktp_siti.jpg"),
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Cilandak",
		AddressDetail: stringPtr("Jl. Fatmawati Raya No. 18, Cilandak"),
		Latitude:      floatPtr(-6.2912),
		Longitude:     floatPtr(106.7978),
	}

	userEmployer1 := models.User{
		ID:            uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33"),
		Email:         "hrd@logisticsjaya.co.id",
		Phone:         "081122334455",
		Name:          "PT Logistics Jaya Mandiri (Hendra)",
		Password:      "Password123!", // Plain text password -> auto-hashed by BeforeCreate hook
		ActiveMode:    models.ModeEmployer,
		IsKTPVerified: true,
		KTPNumber:     stringPtr("3174091211850005"),
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Kebayoran Baru",
		AddressDetail: stringPtr("Gedung Graha Logistics Lt. 4, Kebayoran Baru"),
		Latitude:      floatPtr(-6.2410),
		Longitude:     floatPtr(106.8045),
	}

	userEmployer2 := models.User{
		ID:            uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44"),
		Email:         "anita.wibowo@gmail.com",
		Phone:         "081555666777",
		Name:          "Ibu Anita Wibowo",
		Password:      "Password123!", // Plain text password -> auto-hashed by BeforeCreate hook
		ActiveMode:    models.ModeEmployer,
		IsKTPVerified: true,
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Pondok Indah",
		AddressDetail: stringPtr("Jl. Metro Pondok Indah Blok TA No. 12"),
		Latitude:      floatPtr(-6.2650),
		Longitude:     floatPtr(106.7840),
	}

	users := []models.User{userWorker1, userWorker2, userEmployer1, userEmployer2}
	now := time.Now()
	for i := range users {
		u := &users[i]
		// Semua user contoh sudah terverifikasi KTP agar alur lamaran/tawaran bisa diuji.
		u.KTPStatus = models.KTPVerified
		u.IsKTPVerified = true
		u.KTPSubmittedAt = &now
		u.KTPVerifiedAt = &now
		// Password plaintext di-hash oleh hook BeforeCreate.
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"password", "name", "phone", "active_mode",
				"is_ktp_verified", "ktp_status", "ktp_submitted_at", "ktp_verified_at", "ktp_number",
				"province", "city", "district", "address_detail", "updated_at",
			}),
		}).Create(u).Error; err != nil {
			log.Fatalf("❌ seeding user %s: %v", u.Email, err)
		}
	}
	log.Println("✅ Sample Users seeded and auto-hashed successfully.")

	// 4. Seed Subscriptions (sudah PAID/ACTIVE; order_id LEGACY-SEED-*)
	log.Println("💳 Seeding Subscriptions...")
	subs := []models.Subscription{
		{
			ID:              uuid.MustParse("1bafdaee-cfea-4894-82d1-8c163a1f6acd"),
			UserID:          userWorker1.ID,
			PlanType:        models.PlanPostSkill10K,
			Amount:          10000,
			Status:          models.SubStatusActive,
			PaymentStatus:   models.PaymentPaid,
			PaymentProvider: "SEED",
			OrderID:         "LEGACY-SEED-1bafdaee",
			PaidAt:          &now,
			StartsAt:        now,
			ExpiresAt:       now.Add(models.SubscriptionDuration),
		},
		{
			ID:              uuid.MustParse("097c880d-1d28-424e-aeab-aba9ff9d4550"),
			UserID:          userWorker2.ID,
			PlanType:        models.PlanApplyJob5K,
			Amount:          5000,
			Status:          models.SubStatusActive,
			PaymentStatus:   models.PaymentPaid,
			PaymentProvider: "SEED",
			OrderID:         "LEGACY-SEED-097c880d",
			PaidAt:          &now,
			StartsAt:        now,
			ExpiresAt:       now.Add(models.SubscriptionDuration),
		},
	}
	for _, s := range subs {
		var existing models.Subscription
		if err := db.Where("id = ?", s.ID).First(&existing).Error; err != nil {
			if err := db.Create(&s).Error; err != nil {
				log.Printf("⚠️  seeding subscription %s dilewati: %v", s.ID, err)
			}
		}
	}
	log.Println("✅ Subscriptions seeded successfully.")

	// 5. Seed Model 1: Skill Postings
	log.Println("🛠️ Seeding Skill Postings (Model 1)...")
	skill1 := models.SkillPosting{
		ID:           uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b11"),
		UserID:       userWorker1.ID,
		Category:     models.CategoryProfesional,
		Title:        "Supir Operasional & Pengemudi Pribadi Pengganti",
		Description:  "Pengalaman 6 tahun supir eksekutif & operasional boks/kantor. SIM A & B1 aktif. Siap kerja insidental harian.",
		Country:      "Indonesia",
		Province:     "DKI Jakarta",
		City:         "Jakarta Selatan",
		District:     "Kebayoran Baru",
		RateType:     models.RatePerDay,
		RateAmount:   250000,
		Availability: models.AvailAvailable,
	}

	skill2 := models.SkillPosting{
		ID:           uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b22"),
		UserID:       userWorker2.ID,
		Category:     models.CategorySerabutan,
		Title:        "Bantuan Asisten Rumah Tangga & Cuci Setrika Harian",
		Description:  "Bisa membantu cuci setrika, bersihkan rumah harian, dan bantu persiapan konsumsi acara keluarga. Teliti & jujur.",
		Country:      "Indonesia",
		Province:     "DKI Jakarta",
		City:         "Jakarta Selatan",
		District:     "Cilandak",
		RateType:     models.RatePerHour,
		RateAmount:   35000,
		Availability: models.AvailAvailable,
	}

	skills := []models.SkillPosting{skill1, skill2}
	for _, sk := range skills {
		var existing models.SkillPosting
		if err := db.Where("id = ?", sk.ID).First(&existing).Error; err != nil {
			db.Create(&sk)
		}
	}
	log.Println("✅ Skill Postings seeded successfully.")

	// 6. Seed Model 2: Job Postings
	log.Println("📋 Seeding Job Postings (Model 2)...")
	job1 := models.JobPosting{
		ID:            uuid.MustParse("c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c11"),
		EmployerID:    userEmployer1.ID,
		Category:      models.CategoryProfesional,
		Title:         "Dibutuhkan Supir Operasional Pengganti Shift 2 Hari",
		Description:   "Dibutuhkan pengemudi mobil boks operasional kantor untuk pengiriman area Jabodetabek selama 2 hari pengganti staf sakit.",
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Kebayoran Baru",
		WorkDate:      time.Now().AddDate(0, 0, 2),
		DurationType:  models.DurationDays,
		DurationValue: 2,
		Budget:        600000,
		Status:        models.JobStatusOpen,
	}

	job2 := models.JobPosting{
		ID:            uuid.MustParse("c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c22"),
		EmployerID:    userEmployer2.ID,
		Category:      models.CategorySerabutan,
		Title:         "Bantuan Insidental Bersihkan Kebun & Halaman Rumah",
		Description:   "Dibutuhkan tenaga bantuan 1 hari (sekitar 4 jam) untuk potong rumput & pembersihan halaman belakang rumah.",
		Country:       "Indonesia",
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Pondok Indah",
		WorkDate:      time.Now().AddDate(0, 0, 1),
		DurationType:  models.DurationHours,
		DurationValue: 4,
		Budget:        150000,
		Status:        models.JobStatusOpen,
	}

	jobs := []models.JobPosting{job1, job2}
	for _, j := range jobs {
		var existing models.JobPosting
		if err := db.Where("id = ?", j.ID).First(&existing).Error; err != nil {
			db.Create(&j)
		}
	}
	log.Println("✅ Job Postings seeded successfully.")

	log.Println("🎉 Database Seeding Completed Cleanly!")
}

func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
