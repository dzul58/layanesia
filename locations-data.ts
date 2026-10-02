export interface LocationItem {
  country: string;
  province: string;
  city: string;
  district: string;
  postalCode?: string;
}

export const LOCATIONS: LocationItem[] = [
  // DKI Jakarta - Jakarta Selatan
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Kebayoran Baru', postalCode: '12110' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Kebayoran Lama', postalCode: '12210' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Pesanggrahan', postalCode: '12320' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Cilandak', postalCode: '12430' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Pasar Minggu', postalCode: '12520' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Jagakarsa', postalCode: '12620' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Mampang Prapatan', postalCode: '12790' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Pancoran', postalCode: '12780' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Tebet', postalCode: '12810' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Selatan', district: 'Setiabudi', postalCode: '12910' },

  // DKI Jakarta - Jakarta Pusat
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Gambir', postalCode: '10110' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Tanah Abang', postalCode: '10210' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Menteng', postalCode: '10310' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Senen', postalCode: '10410' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Cempaka Putih', postalCode: '10510' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Johar Baru', postalCode: '10560' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Kemayoran', postalCode: '10610' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Pusat', district: 'Sawah Besar', postalCode: '10710' },

  // DKI Jakarta - Jakarta Barat
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Cengkareng', postalCode: '11730' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Grogol Petamburan', postalCode: '11470' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Kalideres', postalCode: '11840' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Kebon Jeruk', postalCode: '11530' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Kembangan', postalCode: '11610' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Palmerah', postalCode: '11480' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Taman Sari', postalCode: '11110' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Barat', district: 'Tambora', postalCode: '11210' },

  // DKI Jakarta - Jakarta Timur
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Matraman', postalCode: '13110' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Pulo Gadung', postalCode: '13260' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Jatinegara', postalCode: '13310' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Duren Sawit', postalCode: '13440' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Kramat Jati', postalCode: '13510' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Makasar', postalCode: '13570' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Pasar Rebo', postalCode: '13710' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Ciracas', postalCode: '13740' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Cipayung', postalCode: '13840' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Timur', district: 'Cakung', postalCode: '13910' },

  // DKI Jakarta - Jakarta Utara
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Cilincing', postalCode: '14120' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Kelapa Gading', postalCode: '14240' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Koja', postalCode: '14210' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Pademangan', postalCode: '14410' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Penjaringan', postalCode: '14440' },
  { country: 'Indonesia', province: 'DKI Jakarta', city: 'Jakarta Utara', district: 'Tanjung Priok', postalCode: '14310' },

  // Jawa Barat - Kota Bogor
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Bogor Barat', postalCode: '16111' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Bogor Selatan', postalCode: '16132' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Bogor Tengah', postalCode: '16121' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Bogor Timur', postalCode: '16143' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Bogor Utara', postalCode: '16152' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bogor', district: 'Tanah Sareal', postalCode: '16161' },

  // Jawa Barat - Kota Depok
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Beji', postalCode: '16421' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Bojongsari', postalCode: '16516' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Cilodong', postalCode: '16413' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Cimanggis', postalCode: '16451' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Cinere', postalCode: '16514' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Cipayung', postalCode: '16437' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Limo', postalCode: '16515' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Pancoran Mas', postalCode: '16431' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Sawangan', postalCode: '16511' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Sukmajaya', postalCode: '16412' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Depok', district: 'Tapos', postalCode: '16457' },

  // Banten - Kota Tangerang
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang', district: 'Batuceper', postalCode: '15122' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang', district: 'Ciledug', postalCode: '15153' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang', district: 'Cipondoh', postalCode: '15148' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang', district: 'Karawaci', postalCode: '15115' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang', district: 'Tangerang', postalCode: '15111' },

  // Banten - Tangerang Selatan
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang Selatan', district: 'Ciputat', postalCode: '15411' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang Selatan', district: 'Pamulang', postalCode: '15417' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang Selatan', district: 'Pondok Aren', postalCode: '15224' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang Selatan', district: 'Serpong', postalCode: '15310' },
  { country: 'Indonesia', province: 'Banten', city: 'Kota Tangerang Selatan', district: 'Serpong Utara', postalCode: '15326' },

  // Jawa Barat - Kota Bekasi
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bekasi', district: 'Bekasi Barat', postalCode: '17145' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bekasi', district: 'Bekasi Selatan', postalCode: '17141' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bekasi', district: 'Bekasi Timur', postalCode: '17111' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bekasi', district: 'Bekasi Utara', postalCode: '17121' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bekasi', district: 'Pondok Gede', postalCode: '17411' },

  // Jawa Barat - Kota Bandung
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bandung', district: 'Coblong', postalCode: '40132' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bandung', district: 'Cicendo', postalCode: '40171' },
  { country: 'Indonesia', province: 'Jawa Barat', city: 'Kota Bandung', district: 'Sukajadi', postalCode: '40161' },

  // Jawa Timur - Kota Surabaya
  { country: 'Indonesia', province: 'Jawa Timur', city: 'Kota Surabaya', district: 'Tegalsari', postalCode: '60261' },
  { country: 'Indonesia', province: 'Jawa Timur', city: 'Kota Surabaya', district: 'Gubeng', postalCode: '60281' },

  // Bali - Denpasar & Badung
  { country: 'Indonesia', province: 'Bali', city: 'Kota Denpasar', district: 'Denpasar Selatan', postalCode: '80221' },
  { country: 'Indonesia', province: 'Bali', city: 'Kabupaten Badung', district: 'Kuta', postalCode: '80361' },

  // Sumatera Utara - Kota Medan
  { country: 'Indonesia', province: 'Sumatera Utara', city: 'Kota Medan', district: 'Medan Kota', postalCode: '20211' }
];
