import React, { useState } from 'react';
import {
  ShieldCheck, MapPin, Zap, DollarSign,
  Users, Briefcase, User, Sparkles, ChevronDown,
  Globe, ArrowRight, Calculator, Award,
  MessageSquare, Star, Search
} from 'lucide-react';

export default function LandingSections({
  onOpenPostSkill,
  onOpenPostJob,
  onOpenMap,
  onScrollToCatalog,
  activeMode
}) {
  // Workflow toggle state
  const [activeWorkflowRole, setActiveWorkflowRole] = useState(
    activeMode === 'EMPLOYER' ? 'EMPLOYER' : 'SEEKER'
  );

  // Calculator State
  const [selectedSkillType, setSelectedSkillType] = useState('supir');
  const [workDaysPerMonth, setWorkDaysPerMonth] = useState(20);
  const [customDailyRate, setCustomDailyRate] = useState(250000);

  // FAQ Accordion state (index or null)
  const [openFaqIndex, setOpenFaqIndex] = useState(0);

  // Skill presets for calculator
  const skillPresets = {
    supir: { label: 'Supir Pribadi / Operasional Boks', defaultRate: 250000, rateType: 'hari' },
    art: { label: 'Bantuan Rumah Tangga & Cuci Setrika', defaultRate: 35000 * 8, rateType: 'hari' },
    tukang: { label: 'Tukang Servis & Perbaikan Rumah', defaultRate: 200000, rateType: 'hari' },
    gudang: { label: 'Staf Gudang & Bongkar Muat', defaultRate: 220000, rateType: 'hari' },
    acara: { label: 'Kru Lapangan / Bantuan Katering & Acara', defaultRate: 180000, rateType: 'hari' },
  };

  const handlePresetChange = (presetKey) => {
    setSelectedSkillType(presetKey);
    setCustomDailyRate(skillPresets[presetKey].defaultRate);
  };

  // Calculations
  const monthlyGrossIncome = workDaysPerMonth * customDailyRate;
  const layanesiaSubscription = 10000; // Rp 10.000 / month flat
  const platformCommissionComparison = monthlyGrossIncome * 0.20; // 20% on typical gig apps
  const monthlyNetIncome = monthlyGrossIncome - layanesiaSubscription;
  const savedSavings = platformCommissionComparison - layanesiaSubscription;

  const faqs = [
    {
      q: 'Apa itu Layanesia dan bagaimana sistem "Single Account Dual-Role"?',
      a: 'Layanesia adalah platform on-demand lokal untuk tenaga kerja harian, supir pengganti, asisten rumah tangga insidental, hingga staf operasional mendesak di Indonesia. Dengan sistem Single Account Dual-Role, Anda hanya butuh 1 akun untuk beralih secara instan antara mencari pekerjaan harian (Pencari Kerja) atau merekrut tenaga kerja darurat (Pemberi Kerja) tanpa perlu mendaftar ulang.'
    },
    {
      q: 'Mengapa Layanesia mewajibkan Verifikasi KTP resmi?',
      a: 'Keamanan dan keandalan adalah prioritas utama kami. Verifikasi KTP (NIK 16 digit dan foto identitas asli) wajib untuk kedua belah pihak sebelum bertransaksi. Ini mencegah akun bodong, penipuan identitas, maupun tindakan mangkir (no-show) saat pekerjaan berlangsung.'
    },
    {
      q: 'Berapa biaya menggunakan Layanesia? Apakah ada potongan komisi dari gaji?',
      a: 'Pemberi Kerja bebas biaya selamanya: pasang lowongan 100% GRATIS. Untuk Pencari Kerja, Layanesia menggunakan sistem langganan mikro flat: Rp 5.000 / 30 hari untuk melamar semua lowongan (Model 2), atau Rp 10.000 / 30 hari untuk memajang profil keahlian di etalase (Model 1) via pembayaran Midtrans resmi. Hasil jerih payah gaji pekerja 100% utuh tanpa potongan komisi sepeser pun.'
    },
    {
      q: 'Bagaimana alur komunikasi setelah tawaran atau lamaran disetujui?',
      a: 'Begitu tawaran pekerjaan diterima atau pelamar disetujui, sistem Kontak Terhubung langsung membuka detail nomor telepon dan tautan WhatsApp resmi partner. Anda bisa langsung berkoordinasi mengenai waktu, patokan alamat, dan detail teknis pekerjaan secara cepat tanpa birokrasi berbelit.'
    },
    {
      q: 'Bagaimana pembayaran upah pekerja dilakukan?',
      a: 'Pembayaran upah honor disepakati langsung antara Pemberi Kerja dan Pekerja sesuai nominal yang tertera pada tawaran/lowongan. Pembayaran dapat dilakukan secara tunai (cash on delivery) maupun transfer bank/e-wallet langsung begitu pekerjaan selesai.'
    },
    {
      q: 'Bagaimana jika pekerjaan telah selesai? Bagaimana sistem ratingnya?',
      a: 'Setelah pekerjaan rampung, salah satu pihak menandai "Selesai" di menu Kontak Terhubung. Kemudian kedua pihak dapat saling memberikan rating bintang (1-5) dan ulasan pengalaman kerja. Sistem rating dua arah ini menjamin reputasi digital pekerja dan pemberi kerja tetap transparan dan terpercaya.'
    }
  ];

  return (
    <div className="landing-sections-wrapper">
      {/* 1. HERO SECTION (Always Dark Luxury Banner in both Dark & Light themes) */}
      <section id="beranda" className="hero-section-card" style={{
        padding: '48px 36px',
        marginBottom: '40px',
        position: 'relative',
        overflow: 'hidden',
        border: '1px solid rgba(255, 255, 255, 0.12)',
        borderRadius: '20px',
        background: 'linear-gradient(145deg, #070b14 0%, #0f172a 55%, #0a0f1d 100%)',
        boxShadow: '0 20px 50px -10px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05)',
      }}>
        {/* Glow ambient background elements */}
        <div style={{
          position: 'absolute',
          top: '-60px',
          right: '-40px',
          width: '320px',
          height: '320px',
          background: 'radial-gradient(circle, rgba(16, 185, 129, 0.22) 0%, transparent 70%)',
          borderRadius: '50%',
          pointerEvents: 'none',
        }} />
        <div style={{
          position: 'absolute',
          bottom: '-80px',
          left: '20%',
          width: '280px',
          height: '280px',
          background: 'radial-gradient(circle, rgba(99, 102, 241, 0.18) 0%, transparent 70%)',
          borderRadius: '50%',
          pointerEvents: 'none',
        }} />

        <div style={{ position: 'relative', zIndex: 2 }}>
          {/* Tagline kicker */}
          <div className="hero-tagline-kicker" style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: '8px',
            padding: '6px 14px',
            borderRadius: '20px',
            background: 'rgba(16, 185, 129, 0.15)',
            border: '1px solid rgba(16, 185, 129, 0.4)',
            marginBottom: '20px',
            fontSize: '0.85rem',
            color: '#34d399',
            fontWeight: 600,
          }}>
            <Sparkles size={16} color="#34d399" /> Platform On-Demand Tenaga Kerja Pengganti & Insidental No. 1 di Indonesia
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1.4fr) minmax(0, 1fr)', gap: '36px', alignItems: 'center' }}>
            <div>
              <h1 className="hero-heading" style={{
                fontSize: 'clamp(2rem, 3.8vw, 3rem)',
                lineHeight: 1.18,
                marginBottom: '18px',
                color: '#ffffff',
                letterSpacing: '-0.03em',
                textWrap: 'balance',
              }}>
                Solusi Cepat Tenaga Kerja Pengganti & Pekerjaan Insidental{' '}
                <span style={{
                  background: 'linear-gradient(135deg, #34d399 0%, #2dd4bf 50%, #60a5fa 100%)',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}>
                  Dalam Hitungan Jam
                </span>
              </h1>

              <p className="hero-subtext" style={{
                color: '#cbd5e1',
                fontSize: '1.05rem',
                lineHeight: 1.65,
                marginBottom: '28px',
                maxWidth: '620px',
              }}>
                Hubungkan kebutuhan supir pengganti mendesak, asisten rumah tangga harian, staf gudang, hingga tukang perbaikan berlisensi KTP di sekitar kecamatan Anda secara aman, cepat, dan <strong style={{ color: '#f8fafc', fontWeight: 700 }}>tanpa potongan komisi gaji</strong>.
              </p>

              {/* Action Buttons */}
              <div style={{ display: 'flex', gap: '14px', flexWrap: 'wrap', marginBottom: '32px' }}>
                <button
                  className="btn btn-primary"
                  onClick={onScrollToCatalog}
                  style={{
                    padding: '12px 24px',
                    fontSize: '1rem',
                    fontWeight: 700,
                    background: 'linear-gradient(135deg, #10b981 0%, #0d9488 100%)',
                    color: '#ffffff',
                    boxShadow: '0 4px 16px rgba(16, 185, 129, 0.4)',
                  }}
                >
                  <Search size={18} /> Eksplorasi Jasa & Loker
                </button>

                <button
                  className="btn hero-btn-post-job"
                  onClick={onOpenPostJob}
                  style={{
                    border: '1px solid rgba(129, 140, 248, 0.45)',
                    color: '#c7d2fe',
                    background: 'rgba(99, 102, 241, 0.14)',
                    padding: '12px 22px',
                    fontSize: '0.95rem',
                    fontWeight: 600,
                  }}
                >
                  <Briefcase size={18} /> Pasang Lowongan (Gratis)
                </button>

                <button
                  className="btn hero-btn-post-skill"
                  onClick={onOpenPostSkill}
                  style={{
                    border: '1px solid rgba(52, 211, 153, 0.45)',
                    color: '#a7f3d0',
                    background: 'rgba(16, 185, 129, 0.14)',
                    padding: '12px 22px',
                    fontSize: '0.95rem',
                    fontWeight: 600,
                  }}
                >
                  <Zap size={18} /> Pasang Keahlian Saya
                </button>
              </div>

              {/* Unboxed Trust Signals */}
              <div className="hero-trust-bar" style={{
                display: 'flex',
                alignItems: 'center',
                gap: '12px',
                flexWrap: 'wrap',
                fontSize: '0.85rem',
                color: '#cbd5e1',
                paddingTop: '16px',
                borderTop: '1px solid rgba(255, 255, 255, 0.12)',
              }}>
                <span className="hero-trust-item-verified" style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', color: '#34d399', fontWeight: 600 }}>
                  <ShieldCheck size={16} color="#34d399" /> 100% Verifikasi KTP
                </span>
                <span className="hero-trust-divider" aria-hidden="true" style={{ color: 'rgba(255, 255, 255, 0.35)' }}>·</span>
                <span className="hero-trust-item-location" style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', color: '#cbd5e1' }}>
                  <MapPin size={16} color="#818cf8" /> 185+ Kecamatan Jabodetabek
                </span>
                <span className="hero-trust-divider" aria-hidden="true" style={{ color: 'rgba(255, 255, 255, 0.35)' }}>·</span>
                <span className="hero-trust-item-salary" style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', color: '#fbbf24', fontWeight: 600 }}>
                  <Award size={16} color="#fbbf24" /> Gaji 100% Milik Pekerja
                </span>
              </div>
            </div>

            {/* Hero Quick Preview Showcase Box */}
            <div className="hero-preview-box" style={{
              background: 'rgba(255, 255, 255, 0.03)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: '16px',
              padding: '24px',
              backdropFilter: 'blur(12px)',
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                <span style={{ fontSize: '0.85rem', fontWeight: 700, color: '#e2e8f0', letterSpacing: '0.02em' }}>
                  SISTEM DUA MODEL KERJA
                </span>
                <span style={{ fontSize: '0.78rem', color: '#34d399', fontWeight: 600 }}>
                  Live Marketplace
                </span>
              </div>

              {/* Model 1 Mini Showcase Card */}
              <div className="hero-mini-card" style={{
                background: 'rgba(15, 23, 42, 0.85)',
                border: '1px solid rgba(16, 185, 129, 0.35)',
                borderRadius: '12px',
                padding: '14px',
                marginBottom: '12px',
                boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
              }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '6px' }}>
                  <span style={{ fontSize: '0.72rem', color: '#34d399', fontWeight: 700, textTransform: 'uppercase' }}>
                    Model 1 · Talent Showcase
                  </span>
                  <span style={{ fontSize: '0.78rem', color: '#34d399', fontWeight: 700 }}>
                    Rp 250.000 / hari
                  </span>
                </div>
                <div className="hero-mini-title" style={{ fontSize: '0.92rem', fontWeight: 600, color: '#ffffff', marginBottom: '4px' }}>
                  Supir Pengganti SIM A & B1 Aktif
                </div>
                <div className="hero-mini-sub" style={{ fontSize: '0.78rem', color: '#94a3b8' }}>
                  Budi Santoso · Kebayoran Baru, Jakarta Selatan
                </div>
              </div>

              {/* Model 2 Mini Showcase Card */}
              <div className="hero-mini-card" style={{
                background: 'rgba(15, 23, 42, 0.85)',
                border: '1px solid rgba(99, 102, 241, 0.35)',
                borderRadius: '12px',
                padding: '14px',
                marginBottom: '16px',
                boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
              }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '6px' }}>
                  <span style={{ fontSize: '0.72rem', color: '#a5b4fc', fontWeight: 700, textTransform: 'uppercase' }}>
                    Model 2 · Demand Job Board
                  </span>
                  <span style={{ fontSize: '0.78rem', color: '#a5b4fc', fontWeight: 700 }}>
                    Rp 600.000 (2 Hari)
                  </span>
                </div>
                <div className="hero-mini-title" style={{ fontSize: '0.92rem', fontWeight: 600, color: '#ffffff', marginBottom: '4px' }}>
                  Dibutuhkan Supir Operasional Kantor
                </div>
                <div className="hero-mini-sub" style={{ fontSize: '0.78rem', color: '#94a3b8' }}>
                  PT Logistics Jaya Mandiri · Post Gratis 100%
                </div>
              </div>

              {/* Interactive map trigger mini banner */}
              <button
                onClick={onOpenMap}
                className="btn btn-outline"
                style={{
                  width: '100%',
                  padding: '10px',
                  fontSize: '0.85rem',
                  borderColor: 'rgba(16, 185, 129, 0.45)',
                  color: '#34d399',
                  background: 'rgba(16, 185, 129, 0.12)',
                }}
              >
                <Globe size={16} /> Buka Peta Interaktif Lokasi Kerja
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* 2. ADVANTAGES SECTION (KEUNGGULAN UTAMA) */}
      <section id="keunggulan" style={{ marginBottom: '56px' }}>
        <div style={{ textAlign: 'center', maxWidth: '720px', margin: '0 auto 36px auto' }}>
          <div style={{ fontSize: '0.85rem', fontWeight: 700, color: '#10b981', textTransform: 'uppercase', letterSpacing: '0.08em', marginBottom: '8px' }}>
            Mengapa Memilih Layanesia
          </div>
          <h2 style={{ fontSize: '2.1rem', color: 'var(--text-main)', letterSpacing: '-0.02em', marginBottom: '12px' }}>
            Standar Baru Keamanan & Transparansi Pekerja Harian
          </h2>
          <p style={{ color: 'var(--text-muted)', fontSize: '1rem', lineHeight: 1.6 }}>
            Kami merancang ekosistem gig terpercaya yang memprioritaskan keamanan identitas, kecepatan penugasan, dan keadilan finansial bagi semua pihak.
          </p>
        </div>

        {/* Bento Grid */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
          gap: '20px',
        }}>
          {/* Advantage 1 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(16, 185, 129, 0.15)',
              color: '#10b981',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <ShieldCheck size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              01. Verifikasi KTP Resmi Anti-Akun Bodong
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Pekerja dan pemberi kerja wajib memvalidasi NIK 16 digit dan foto KTP resmi. Menghapus risiko penipuan identitas, profil fiktif, maupun kekhawatiran barang hilang di lokasi kerja.
            </p>
          </div>

          {/* Advantage 2 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(99, 102, 241, 0.15)',
              color: '#818cf8',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <MapPin size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              02. Presisi Hiperlokal Tingkat Kecamatan
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Didukung 100% master data 185 kecamatan Jabodetabek dan kota besar Indonesia. Dilengkapi Leaflet OpenStreetMap untuk menemukan kandidat terdekat tanpa buang waktu macet di jalan.
            </p>
          </div>

          {/* Advantage 3 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(245, 158, 11, 0.15)',
              color: '#f59e0b',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <DollarSign size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              03. Tanpa Potongan Komisi Hasil Keringat
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Bebas potongan 20%-30% dari upah. Pekerja hanya membayar langganan flat terjangkau (Rp 5.000 atau Rp 10.000/30 hari via Midtrans). Gaji diterima 100% utuh langsung ke kantong pekerja.
            </p>
          </div>

          {/* Advantage 4 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(14, 165, 233, 0.15)',
              color: '#38bdf8',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <Users size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              04. Single Account Dual-Role Fleksibel
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Satu akun untuk dua fungsi. Pagi hari Anda bisa bertindak sebagai employer mencari tukang ledeng darurat, lalu di akhir pekan beralih menawarkan keahlian fotografi atau les privat.
            </p>
          </div>

          {/* Advantage 5 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(34, 197, 94, 0.15)',
              color: '#4ade80',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <MessageSquare size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              05. Kontak WhatsApp Instan Begitu Matched
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Saat tawaran disetujui atau pelamar diterima, tombol direct chat WhatsApp partner langsung terbuka otomatis. Koordinasi alamat patokan dan jadwal langsung lancar tanpa jeda.
            </p>
          </div>

          {/* Advantage 6 */}
          <div className="glass-panel" style={{ padding: '28px', borderRadius: '16px' }}>
            <div style={{
              width: '46px',
              height: '46px',
              borderRadius: '12px',
              background: 'rgba(236, 72, 153, 0.15)',
              color: '#f472b6',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '18px',
            }}>
              <Star size={26} />
            </div>
            <h3 style={{ fontSize: '1.18rem', color: 'var(--text-main)', marginBottom: '10px' }}>
              06. Evaluasi Rating & Ulasan Dua Arah
            </h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Keduanya saling mengulas setelah pekerjaan diselesaikan. Pemberi kerja mendapatkan kepastian etos kerja, pekerja mendapatkan bukti portofolio reputasi digital yang kredibel.
            </p>
          </div>
        </div>
      </section>

      {/* 3. HOW IT WORKS (CARA KERJA DUA ARAH) */}
      <section id="cara-kerja" style={{ marginBottom: '56px' }}>
        <div style={{ textAlign: 'center', maxWidth: '720px', margin: '0 auto 32px auto' }}>
          <div style={{ fontSize: '0.85rem', fontWeight: 700, color: '#6366f1', textTransform: 'uppercase', letterSpacing: '0.08em', marginBottom: '8px' }}>
            Alur Penggunaan
          </div>
          <h2 style={{ fontSize: '2.1rem', color: 'var(--text-main)', letterSpacing: '-0.02em', marginBottom: '12px' }}>
            Bagaimana Cara Kerja Layanesia?
          </h2>
          <p style={{ color: 'var(--text-muted)', fontSize: '1rem', lineHeight: 1.6 }}>
            Pilih alur peran Anda di bawah ini untuk melihat langkah mudah memulai di Layanesia:
          </p>

          {/* Interactive Role Segmented Switcher */}
          <div style={{
            display: 'inline-flex',
            padding: '4px',
            borderRadius: '14px',
            background: 'var(--bg-secondary)',
            border: '1px solid var(--border-color)',
            marginTop: '16px',
          }}>
            <button
              onClick={() => setActiveWorkflowRole('SEEKER')}
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: '8px',
                padding: '10px 24px',
                borderRadius: '10px',
                fontSize: '0.92rem',
                fontWeight: 600,
                border: 'none',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                background: activeWorkflowRole === 'SEEKER' ? 'var(--gradient-primary)' : 'transparent',
                color: activeWorkflowRole === 'SEEKER' ? '#ffffff' : 'var(--text-muted)',
              }}
            >
              <User size={16} /> Saya Pencari Kerja (Worker)
            </button>
            <button
              onClick={() => setActiveWorkflowRole('EMPLOYER')}
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: '8px',
                padding: '10px 24px',
                borderRadius: '10px',
                fontSize: '0.92rem',
                fontWeight: 600,
                border: 'none',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                background: activeWorkflowRole === 'EMPLOYER' ? 'var(--gradient-secondary)' : 'transparent',
                color: activeWorkflowRole === 'EMPLOYER' ? '#ffffff' : 'var(--text-muted)',
              }}
            >
              <Briefcase size={16} /> Saya Pemberi Kerja (Employer)
            </button>
          </div>
        </div>

        {/* WORKFLOW STEPS */}
        {activeWorkflowRole === 'SEEKER' ? (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))',
            gap: '20px',
          }}>
            {/* Step 1 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px', position: 'relative' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#10b981', marginBottom: '8px' }}>
                LANGKAH 01
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Daftar & Verifikasi KTP
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Daftarkan akun gratis, isi data domisili, unggah foto KTP resmi dan file CV/Resume PDF Anda untuk mendapatkan lencana terverifikasi.
              </p>
            </div>

            {/* Step 2 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#10b981', marginBottom: '8px' }}>
                LANGKAH 02
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Pasang Jasa / Lamar Loker
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Pilih model Anda: Pasang profil keahlian di Etalase (Model 1, Rp 10k/30 hari) atau lamar lowongan insidental terbuka (Model 2, Rp 5k/30 hari).
              </p>
            </div>

            {/* Step 3 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#10b981', marginBottom: '8px' }}>
                LANGKAH 03
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Terima Tawaran & Chat WA
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Terima penawaran langsung atau notifikasi lamaran diterima. Buka menu Kontak Terhubung untuk chat langsung via WhatsApp resmi.
              </p>
            </div>

            {/* Step 4 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#10b981', marginBottom: '8px' }}>
                LANGKAH 04
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Kerja & Terima Gaji Utuh
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Selesaikan pekerjaan di lokasi. Terima bayaran honor penuh tanpa potongan sepersen pun, lalu tukar ulasan rating untuk portofolio Anda.
              </p>
            </div>
          </div>
        ) : (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))',
            gap: '20px',
          }}>
            {/* Step 1 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#6366f1', marginBottom: '8px' }}>
                LANGKAH 01
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Buka Akun Pemberi Kerja
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Masuk ke akun Layanesia Anda. Anda bisa mengaktifkan mode Pemberi Kerja kapan saja tanpa biaya pendaftaran.
              </p>
            </div>

            {/* Step 2 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#6366f1', marginBottom: '8px' }}>
                LANGKAH 02
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Post Loker Gratis / Cari Talent
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Publikasikan lowongan kerja insidental secara 100% GRATIS atau telusuri profil pekerja terdekat di etalase hiperlokal.
              </p>
            </div>

            {/* Step 3 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#6366f1', marginBottom: '8px' }}>
                LANGKAH 03
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Kirim Tawaran / Pilih Pelamar
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Kirim Direct Offer dengan budget Anda ke pekerja, atau reviu daftar CV pelamar yang masuk lalu terima kandidat terbaik.
              </p>
            </div>

            {/* Step 4 */}
            <div className="glass-panel" style={{ padding: '24px', borderRadius: '16px' }}>
              <div style={{ fontSize: '0.8rem', fontWeight: 800, color: '#6366f1', marginBottom: '8px' }}>
                LANGKAH 04
              </div>
              <h4 style={{ fontSize: '1.1rem', color: 'var(--text-main)', marginBottom: '8px' }}>
                Terhubung Langsung & Evaluasi
              </h4>
              <p style={{ fontSize: '0.86rem', color: 'var(--text-muted)', lineHeight: 1.55 }}>
                Hubungi via WhatsApp, koordinasikan pekerjaan, dan setelah selesai berikan bintang ulasan untuk memastikan standar kualitas.
              </p>
            </div>
          </div>
        )}
      </section>

      {/* 4. INTERACTIVE CALCULATOR (SIMULASI PENGHASILAN & ANGGARAN) */}
      <section id="kalkulator" className="glass-panel" style={{
        padding: '36px',
        borderRadius: '20px',
        marginBottom: '56px',
        border: '1px solid var(--border-color)',
        background: 'var(--calc-bg)',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', color: '#10b981', marginBottom: '8px' }}>
          <Calculator size={22} />
          <span style={{ fontSize: '0.85rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Simulasi & Transparansi Biaya
          </span>
        </div>

        <h2 style={{ fontSize: '1.9rem', color: 'var(--text-main)', marginBottom: '8px', letterSpacing: '-0.02em' }}>
          Hitung Estimasi Pendapatan & Penghematan di Layanesia
        </h2>
        <p style={{ color: 'var(--text-muted)', fontSize: '0.95rem', marginBottom: '28px', maxWidth: '700px' }}>
          Bandingkan hasil kerja Anda di Layanesia dengan platform konvensional yang memotong komisi 20% dari setiap pekerjaan.
        </p>

        <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1.2fr) minmax(0, 1fr)', gap: '32px' }}>
          {/* Controls */}
          <div>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '8px' }}>
              Pilih Profesi / Keahlian:
            </label>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginBottom: '20px' }}>
              {Object.entries(skillPresets).map(([key, item]) => (
                <button
                  key={key}
                  type="button"
                  onClick={() => handlePresetChange(key)}
                  style={{
                    padding: '8px 14px',
                    borderRadius: '10px',
                    fontSize: '0.85rem',
                    fontWeight: 600,
                    cursor: 'pointer',
                    transition: 'all 0.2s ease',
                    background: selectedSkillType === key ? 'rgba(16, 185, 129, 0.2)' : 'var(--bg-secondary)',
                    color: selectedSkillType === key ? '#10b981' : 'var(--text-main)',
                    border: selectedSkillType === key ? '1px solid #10b981' : '1px solid var(--border-color)',
                  }}
                >
                  {item.label}
                </button>
              ))}
            </div>

            <div style={{ marginBottom: '20px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px', fontSize: '0.85rem' }}>
                <span style={{ color: 'var(--text-muted)', fontWeight: 600 }}>Tarif Harian Rata-rata:</span>
                <span style={{ color: '#10b981', fontWeight: 700, fontFamily: 'monospace' }}>
                  Rp {customDailyRate.toLocaleString('id-ID')} / hari
                </span>
              </div>
              <input
                type="range"
                min="50000"
                max="600000"
                step="10000"
                value={customDailyRate}
                onChange={(e) => setCustomDailyRate(parseInt(e.target.value, 10))}
                style={{ width: '100%', accentColor: '#10b981', cursor: 'pointer' }}
              />
            </div>

            <div style={{ marginBottom: '24px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px', fontSize: '0.85rem' }}>
                <span style={{ color: 'var(--text-muted)', fontWeight: 600 }}>Hari Kerja Per Bulan:</span>
                <span style={{ color: '#6366f1', fontWeight: 700, fontFamily: 'monospace' }}>
                  {workDaysPerMonth} Hari
                </span>
              </div>
              <input
                type="range"
                min="1"
                max="30"
                value={workDaysPerMonth}
                onChange={(e) => setWorkDaysPerMonth(parseInt(e.target.value, 10))}
                style={{ width: '100%', accentColor: '#6366f1', cursor: 'pointer' }}
              />
            </div>
          </div>

          {/* Result Card */}
          <div style={{
            background: 'var(--bg-secondary)',
            border: '1px solid var(--border-color)',
            borderRadius: '16px',
            padding: '24px',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'space-between',
          }}>
            <div>
              <span style={{ fontSize: '0.78rem', color: 'var(--text-dim)', textTransform: 'uppercase', fontWeight: 700 }}>
                HASIL PERHITUNGAN BERSIH
              </span>
              <div style={{ fontSize: '2rem', fontWeight: 800, color: '#10b981', margin: '4px 0 16px 0', fontFamily: 'monospace' }}>
                Rp {monthlyNetIncome.toLocaleString('id-ID')}
                <span style={{ fontSize: '0.9rem', color: 'var(--text-muted)', fontWeight: 400, marginLeft: '6px' }}>/ bulan</span>
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', fontSize: '0.88rem', borderTop: '1px dashed var(--border-color)', paddingTop: '16px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-muted)' }}>Estimasi Pendapatan Kotor:</span>
                  <strong style={{ color: 'var(--text-main)', fontFamily: 'monospace' }}>Rp {monthlyGrossIncome.toLocaleString('id-ID')}</strong>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-muted)' }}>Biaya Langganan Layanesia:</span>
                  <strong style={{ color: '#10b981', fontFamily: 'monospace' }}>- Rp {layanesiaSubscription.toLocaleString('id-ID')} (Flat)</strong>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', color: '#f59e0b', background: 'rgba(245, 158, 11, 0.1)', padding: '8px 12px', borderRadius: '8px' }}>
                  <span>Penghematan vs Platform Lain (Komisi 20%):</span>
                  <strong style={{ fontFamily: 'monospace' }}>+ Rp {savedSavings.toLocaleString('id-ID')}</strong>
                </div>
              </div>
            </div>

            <div style={{ marginTop: '20px' }}>
              <button
                onClick={onOpenPostSkill}
                className="btn btn-primary"
                style={{ width: '100%', padding: '10px', fontSize: '0.9rem' }}
              >
                Mulai Pasang Keahlian Anda Sekarang <ArrowRight size={16} />
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* 5. FAQ SECTION */}
      <section id="faq" style={{ marginBottom: '56px' }}>
        <div style={{ textAlign: 'center', maxWidth: '720px', margin: '0 auto 36px auto' }}>
          <div style={{ fontSize: '0.85rem', fontWeight: 700, color: '#f59e0b', textTransform: 'uppercase', letterSpacing: '0.08em', marginBottom: '8px' }}>
            Pertanyaan Umum
          </div>
          <h2 style={{ fontSize: '2.1rem', color: 'var(--text-main)', letterSpacing: '-0.02em', marginBottom: '12px' }}>
            Pertanyaan yang Sering Diajukan
          </h2>
          <p style={{ color: 'var(--text-muted)', fontSize: '1rem', lineHeight: 1.6 }}>
            Semua hal yang perlu Anda ketahui mengenai keamanan, akun, pembayaran, dan sistem kerja di Layanesia.
          </p>
        </div>

        <div style={{ maxWidth: '840px', margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {faqs.map((faq, idx) => {
            const isOpen = openFaqIndex === idx;
            return (
              <div
                key={idx}
                className="glass-panel"
                style={{
                  borderRadius: '14px',
                  border: isOpen ? '1px solid rgba(16, 185, 129, 0.4)' : '1px solid var(--border-color)',
                  overflow: 'hidden',
                  transition: 'all 0.2s ease',
                }}
              >
                <button
                  type="button"
                  onClick={() => setOpenFaqIndex(isOpen ? null : idx)}
                  style={{
                    width: '100%',
                    padding: '18px 22px',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    background: 'transparent',
                    border: 'none',
                    textAlign: 'left',
                    color: 'var(--text-main)',
                    fontSize: '1.02rem',
                    fontWeight: 600,
                    cursor: 'pointer',
                    gap: '16px',
                  }}
                >
                  <span>{faq.q}</span>
                  <ChevronDown
                    size={20}
                    style={{
                      transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)',
                      transition: 'transform 0.25s ease',
                      color: isOpen ? '#10b981' : 'var(--text-muted)',
                      flexShrink: 0,
                    }}
                  />
                </button>

                {isOpen && (
                  <div style={{
                    padding: '0 22px 18px 22px',
                    fontSize: '0.92rem',
                    color: 'var(--text-muted)',
                    lineHeight: 1.65,
                    borderTop: '1px solid var(--border-color)',
                    paddingTop: '14px',
                  }}>
                    {faq.a}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </section>

      {/* 6. CALL TO ACTION BANNER */}
      <section style={{
        background: 'var(--cta-banner-bg)',
        border: '1px solid rgba(16, 185, 129, 0.3)',
        borderRadius: '20px',
        padding: '44px 32px',
        textAlign: 'center',
        marginBottom: '40px',
      }}>
        <div style={{ maxWidth: '640px', margin: '0 auto' }}>
          <h3 style={{ fontSize: '1.8rem', color: 'var(--text-main)', marginBottom: '12px', letterSpacing: '-0.02em' }}>
            Siap Menemukan Tenaga Kerja Pengganti atau Pekerjaan Baru Hari Ini?
          </h3>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.98rem', marginBottom: '24px', lineHeight: 1.6 }}>
            Bergabunglah dengan ribuan pekerja dan pemberi kerja terverifikasi di Jabodetabek & seluruh Indonesia. Tanpa biaya pendaftaran.
          </p>
          <div style={{ display: 'flex', gap: '14px', justifyContent: 'center', flexWrap: 'wrap' }}>
            <button
              onClick={onOpenPostSkill}
              className="btn btn-primary"
              style={{ padding: '12px 24px', fontSize: '0.95rem' }}
            >
              <Zap size={18} /> Pasang Keahlian Saya
            </button>
            <button
              onClick={onOpenPostJob}
              className="btn btn-outline"
              style={{
                borderColor: 'rgba(99, 102, 241, 0.4)',
                padding: '12px 24px',
                fontSize: '0.95rem',
                fontWeight: 600,
              }}
            >
              <Briefcase size={18} color="#6366f1" /> Buka Lowongan Gratis
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}
