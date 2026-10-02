-- Master Data Locations Seed File
-- JABODETABEK 100% COVERAGE (ALL 185 DISTRICTS/KECAMATAN) + MAJOR CITIES ACROSS INDONESIA
-- Schema: (country, province, city_or_regency, district, postal_code)

INSERT INTO locations (country, province, city_or_regency, district, postal_code) VALUES

-- ==============================================================================
-- 1. DKI JAKARTA (44 KECAMATAN - 100% COMPLETE)
-- ==============================================================================
-- Jakarta Selatan (10 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Baru', '12110'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Lama', '12210'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Pesanggrahan', '12320'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Cilandak', '12430'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Pasar Minggu', '12520'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Jagakarsa', '12620'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Mampang Prapatan', '12790'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Pancoran', '12780'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Tebet', '12810'),
('Indonesia', 'DKI Jakarta', 'Jakarta Selatan', 'Setiabudi', '12910'),

-- Jakarta Pusat (8 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Gambir', '10110'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Tanah Abang', '10210'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Menteng', '10310'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Senen', '10410'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Cempaka Putih', '10510'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Johar Baru', '10560'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Kemayoran', '10610'),
('Indonesia', 'DKI Jakarta', 'Jakarta Pusat', 'Sawah Besar', '10710'),

-- Jakarta Barat (8 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Cengkareng', '11730'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Grogol Petamburan', '11470'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Kalideres', '11840'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Kebon Jeruk', '11530'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Kembangan', '11610'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Palmerah', '11480'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Taman Sari', '11110'),
('Indonesia', 'DKI Jakarta', 'Jakarta Barat', 'Tambora', '11210'),

-- Jakarta Timur (10 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Matraman', '13110'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Pulo Gadung', '13260'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Jatinegara', '13310'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Duren Sawit', '13440'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Kramat Jati', '13510'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Makasar', '13570'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Pasar Rebo', '13710'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Ciracas', '13740'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Cipayung', '13840'),
('Indonesia', 'DKI Jakarta', 'Jakarta Timur', 'Cakung', '13910'),

-- Jakarta Utara (6 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Cilincing', '14120'),
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Kelapa Gading', '14240'),
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Koja', '14210'),
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Pademangan', '14410'),
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Penjaringan', '14440'),
('Indonesia', 'DKI Jakarta', 'Jakarta Utara', 'Tanjung Priok', '14310'),

-- Kepulauan Seribu (2 Kecamatan)
('Indonesia', 'DKI Jakarta', 'Kepulauan Seribu', 'Kepulauan Seribu Selatan', '14510'),
('Indonesia', 'DKI Jakarta', 'Kepulauan Seribu', 'Kepulauan Seribu Utara', '14520'),

-- ==============================================================================
-- 2. BOGOR (KOTA BOGOR: 6 KECAMATAN, KABUPATEN BOGOR: 40 KECAMATAN - 100% COMPLETE)
-- ==============================================================================
-- Kota Bogor (6 Kecamatan)
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Bogor Barat', '16111'),
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Bogor Selatan', '16132'),
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Bogor Tengah', '16121'),
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Bogor Timur', '16143'),
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Bogor Utara', '16152'),
('Indonesia', 'Jawa Barat', 'Kota Bogor', 'Tanah Sareal', '16161'),

-- Kabupaten Bogor (40 Kecamatan)
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Babakan Madang', '16810'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Bojonggede', '16922'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Caringin', '16730'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cariu', '16840'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Ciampea', '16620'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Ciawi', '16720'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cibinong', '16911'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cibungbulang', '16630'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cigombong', '16740'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cigudeg', '16660'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cijeruk', '16740'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cileungsi', '16820'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Ciomas', '16610'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Cisarua', '16750'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Ciseeng', '16120'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Citeureup', '16810'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Dramaga', '16680'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Gunung Putri', '16961'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Gunung Sindur', '16340'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Jasinga', '16670'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Jonggol', '16830'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Kemang', '16310'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Klapanunggal', '16710'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Leuwiliang', '16640'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Leuwisadeng', '16650'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Megamendung', '16770'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Nanggung', '16650'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Pamijahan', '16630'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Parung', '16330'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Parung Panjang', '16360'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Ranca Bungur', '16310'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Rumpin', '16350'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Sukajaya', '16660'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Sukamakmur', '16830'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Sukaraja', '16710'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Tajurhalang', '16320'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Tamansari', '16610'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Tanjungsari', '16840'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Tenjo', '16370'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bogor', 'Tenjolaya', '16620'),

-- ==============================================================================
-- 3. DEPOK (KOTA DEPOK: 11 KECAMATAN - 100% COMPLETE)
-- ==============================================================================
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Beji', '16421'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Bojongsari', '16516'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Cilodong', '16413'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Cimanggis', '16451'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Cinere', '16514'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Cipayung', '16437'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Limo', '16515'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Pancoran Mas', '16431'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Sawangan', '16511'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Sukmajaya', '16412'),
('Indonesia', 'Jawa Barat', 'Kota Depok', 'Tapos', '16457'),

-- ==============================================================================
-- 4. TANGERANG (KOTA TANGERANG: 13, TANGSEL: 7, KAB TANGERANG: 29 - 100% COMPLETE)
-- ==============================================================================
-- Kota Tangerang (13 Kecamatan)
('Indonesia', 'Banten', 'Kota Tangerang', 'Batuceper', '15122'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Benda', '15121'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Cibodas', '15138'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Ciledug', '15153'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Cipondoh', '15148'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Jatiuwung', '15134'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Karangtengah', '15157'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Karawaci', '15115'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Larangan', '15154'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Neglasari', '15129'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Periuk', '15131'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Pinang', '15145'),
('Indonesia', 'Banten', 'Kota Tangerang', 'Tangerang', '15111'),

-- Kota Tangerang Selatan (7 Kecamatan)
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Ciputat', '15411'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Ciputat Timur', '15419'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Pamulang', '15417'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Pondok Aren', '15224'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Serpong', '15310'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Serpong Utara', '15326'),
('Indonesia', 'Banten', 'Kota Tangerang Selatan', 'Setu', '15314'),

-- Kabupaten Tangerang (29 Kecamatan)
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Balaraja', '15610'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Cikupa', '15710'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Cisauk', '15341'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Cisoka', '15730'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Curug', '15810'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Gunungkaler', '15620'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Jambe', '15720'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Jayanti', '15610'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Kelapa Dua', '15810'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Kemiri', '15530'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Kosambi', '15211'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Kresek', '15620'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Kronjo', '15550'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Legok', '15820'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Mauk', '15530'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Mekarbaru', '15550'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Pagedangan', '15339'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Pakuhaji', '15570'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Panongan', '15711'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Pasar Kemis', '15560'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Rajeg', '15540'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Sepatan', '15520'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Sepatan Timur', '15521'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Sindang Jaya', '15561'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Solear', '15731'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Sukadiri', '15531'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Sukamulya', '15611'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Teluknaga', '15510'),
('Indonesia', 'Banten', 'Kabupaten Tangerang', 'Tigaraksa', '15720'),

-- ==============================================================================
-- 5. BEKASI (KOTA BEKASI: 12 KECAMATAN, KABUPATEN BEKASI: 23 KECAMATAN - 100% COMPLETE)
-- ==============================================================================
-- Kota Bekasi (12 Kecamatan)
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Bantar Gebang', '17151'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Bekasi Barat', '17145'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Bekasi Selatan', '17141'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Bekasi Timur', '17111'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Bekasi Utara', '17121'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Jatiasih', '17423'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Jatisampurna', '17433'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Medan Satria', '17132'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Mustikajaya', '17158'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Pondok Gede', '17411'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Pondok Melati', '17415'),
('Indonesia', 'Jawa Barat', 'Kota Bekasi', 'Rawalumbu', '17116'),

-- Kabupaten Bekasi (23 Kecamatan)
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Babelan', '17610'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Bojongmangu', '17356'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cabangbingin', '17720'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cibarusah', '17340'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cibitung', '17520'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cikarang Barat', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cikarang Pusat', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cikarang Selatan', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cikarang Timur', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Cikarang Utara', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Karangbahagia', '17535'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Kedungwaringin', '17540'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Muara Gembong', '17730'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Pebayuran', '17710'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Serang Baru', '17330'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Setu', '17320'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Sukakarya', '17530'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Sukatani', '17630'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Sukawangi', '17620'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Tambelang', '17620'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Tambun Selatan', '17510'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Tambun Utara', '17510'),
('Indonesia', 'Jawa Barat', 'Kabupaten Bekasi', 'Tarumajaya', '17210'),

-- ==============================================================================
-- 6. PROVINSI & KOTA UTAMA LAINNYA DI INDONESIA
-- ==============================================================================
-- Jawa Barat (Bandung, Cimahi, Sukabumi, Tasikmalaya, Karawang)
('Indonesia', 'Jawa Barat', 'Kota Bandung', 'Coblong', '40132'),
('Indonesia', 'Jawa Barat', 'Kota Bandung', 'Cicendo', '40171'),
('Indonesia', 'Jawa Barat', 'Kota Bandung', 'Sumur Bandung', '40111'),
('Indonesia', 'Jawa Barat', 'Kota Bandung', 'Lengkong', '40261'),
('Indonesia', 'Jawa Barat', 'Kota Bandung', 'Sukajadi', '40161'),
('Indonesia', 'Jawa Barat', 'Kota Cimahi', 'Cimahi Tengah', '40522'),
('Indonesia', 'Jawa Barat', 'Kota Sukabumi', 'Cikole', '43111'),
('Indonesia', 'Jawa Barat', 'Kota Tasikmalaya', 'Cihideung', '46122'),
('Indonesia', 'Jawa Barat', 'Kabupaten Karawang', 'Karawang Barat', '41311'),
('Indonesia', 'Jawa Barat', 'Kabupaten Karawang', 'Karawang Timur', '41314'),

-- Banten (Serang, Cilegon, Lebak, Pandeglang)
('Indonesia', 'Banten', 'Kota Serang', 'Serang', '42111'),
('Indonesia', 'Banten', 'Kota Cilegon', 'Cibeber', '42422'),
('Indonesia', 'Banten', 'Kabupaten Lebak', 'Rangkasbitung', '42311'),
('Indonesia', 'Banten', 'Kabupaten Pandeglang', 'Pandeglang', '42211'),

-- Jawa Tengah & DIY
('Indonesia', 'Jawa Tengah', 'Kota Semarang', 'Semarang Tengah', '50131'),
('Indonesia', 'Jawa Tengah', 'Kota Surakarta (Solo)', 'Banjarsari', '57131'),
('Indonesia', 'Jawa Tengah', 'Kota Magelang', 'Magelang Tengah', '56111'),
('Indonesia', 'Jawa Tengah', 'Kabupaten Banyumas', 'Purwokerto Timur', '53111'),
('Indonesia', 'DI Yogyakarta', 'Kota Yogyakarta', 'Gondomanan', '55121'),
('Indonesia', 'DI Yogyakarta', 'Kabupaten Sleman', 'Depok', '55281'),
('Indonesia', 'DI Yogyakarta', 'Kabupaten Bantul', 'Sewon', '55188'),

-- Jawa Timur
('Indonesia', 'Jawa Timur', 'Kota Surabaya', 'Tegalsari', '60261'),
('Indonesia', 'Jawa Timur', 'Kota Surabaya', 'Gubeng', '60281'),
('Indonesia', 'Jawa Timur', 'Kota Malang', 'Lowokwaru', '65141'),
('Indonesia', 'Jawa Timur', 'Kota Batu', 'Batu', '65311'),
('Indonesia', 'Jawa Timur', 'Kabupaten Sidoarjo', 'Sidoarjo', '61211'),
('Indonesia', 'Jawa Timur', 'Kabupaten Gresik', 'Gresik', '61111'),

-- Bali
('Indonesia', 'Bali', 'Kota Denpasar', 'Denpasar Selatan', '80221'),
('Indonesia', 'Bali', 'Kabupaten Badung', 'Kuta', '80361'),
('Indonesia', 'Bali', 'Kabupaten Gianyar', 'Ubud', '80571'),

-- Sumatera
('Indonesia', 'Sumatera Utara', 'Kota Medan', 'Medan Kota', '20211'),
('Indonesia', 'Riau', 'Kota Pekanbaru', 'Tampan', '28291'),
('Indonesia', 'Kepulauan Riau', 'Kota Batam', 'Batam Kota', '29432'),
('Indonesia', 'Sumatera Barat', 'Kota Padang', 'Padang Barat', '25111'),
('Indonesia', 'Sumatera Selatan', 'Kota Palembang', 'Ilir Barat I', '30139'),
('Indonesia', 'Lampung', 'Kota Bandar Lampung', 'Tanjung Karang Pusat', '35111'),

-- Kalimantan & Sulawesi & Papua
('Indonesia', 'Kalimantan Timur', 'Kota Balikpapan', 'Balikpapan Kota', '76111'),
('Indonesia', 'Kalimantan Timur', 'Kota Samarinda', 'Samarinda Kota', '75121'),
('Indonesia', 'Kalimantan Barat', 'Kota Pontianak', 'Pontianak Selatan', '78121'),
('Indonesia', 'Sulawesi Selatan', 'Kota Makassar', 'Ujung Pandang', '90111'),
('Indonesia', 'Sulawesi Utara', 'Kota Manado', 'Wenang', '95111'),
('Indonesia', 'Nusa Tenggara Barat', 'Kota Mataram', 'Mataram', '83121'),
('Indonesia', 'Nusa Tenggara Timur', 'Kabupaten Manggarai Barat', 'Komodo', '86554'),
('Indonesia', 'Papua', 'Kota Jayapura', 'Jayapura Utara', '99111');
