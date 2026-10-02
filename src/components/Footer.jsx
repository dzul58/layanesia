import React from 'react';
import { Zap, ShieldCheck } from 'lucide-react';

export default function Footer() {
  return (
    <footer style={{
      borderTop: '1px solid var(--border-color)',
      padding: '48px 24px 32px 24px',
      marginTop: '60px',
      background: 'var(--bg-secondary)',
      color: 'var(--text-muted)',
      fontSize: '0.9rem',
      transition: 'background-color 0.3s ease',
    }}>
      <div className="container" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '32px' }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--text-main)', fontSize: '1.25rem', fontWeight: 700, marginBottom: '10px' }}>
            <Zap color="#10b981" size={22} /> Layanesia
          </div>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.88rem', lineHeight: 1.6, marginBottom: '14px' }}>
            Platform marketplace tenaga kerja lokal terpercaya dan penghubung jasa harian/insidental terverifikasi KTP di Indonesia.
          </p>
          <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '0.8rem', color: '#10b981', fontWeight: 600 }}>
            <ShieldCheck size={16} /> 100% Verifikasi KTP Resmi
          </div>
        </div>

        <div>
          <h4 style={{ color: 'var(--text-main)', marginBottom: '14px', fontSize: '0.95rem', fontWeight: 700 }}>Navigasi Cepat</h4>
          <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: '8px', fontSize: '0.88rem' }}>
            <li><a href="#beranda" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Beranda Utama</a></li>
            <li><a href="#keunggulan" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Keunggulan Layanesia</a></li>
            <li><a href="#cara-kerja" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Cara Kerja Dua Arah</a></li>
            <li><a href="#katalog" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Bursa Jasa & Lowongan</a></li>
            <li><a href="#kalkulator" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Kalkulator Pendapatan</a></li>
            <li><a href="#faq" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Pertanyaan Umum (FAQ)</a></li>
          </ul>
        </div>

        <div>
          <h4 style={{ color: 'var(--text-main)', marginBottom: '14px', fontSize: '0.95rem', fontWeight: 700 }}>Model Bisnis Adil</h4>
          <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: '8px', fontSize: '0.88rem' }}>
            <li>Pasang Lowongan Kerja: <strong style={{ color: '#10b981' }}>100% GRATIS</strong></li>
            <li>Paket Melamar (Pekerja): <strong>Rp 5.000 / 30 Hari</strong></li>
            <li>Paket Etalase Jasa: <strong>Rp 10.000 / 30 Hari</strong></li>
            <li>Potongan Komisi Gaji: <strong style={{ color: '#10b981' }}>0% (Utuh)</strong></li>
            <li>Pembayaran Langganan Resmi via Midtrans</li>
          </ul>
        </div>

        <div>
          <h4 style={{ color: 'var(--text-main)', marginBottom: '14px', fontSize: '0.95rem', fontWeight: 700 }}>Cakupan Layanan</h4>
          <p style={{ fontSize: '0.88rem', color: 'var(--text-muted)', lineHeight: 1.6, margin: 0 }}>
            Jabodetabek (DKI Jakarta, Bogor, Depok, Tangerang, Tangerang Selatan, Bekasi) serta kota-kota besar utama di seluruh Indonesia.
          </p>
        </div>
      </div>

      <div style={{
        textAlign: 'center',
        marginTop: '36px',
        paddingTop: '20px',
        borderTop: '1px solid var(--border-color)',
        fontSize: '0.82rem',
        color: 'var(--text-dim)',
      }}>
        © 2026 Layanesia. Platform Tenaga Kerja Harian & Insidental Terpercaya. Didukung Peta Interaktif OpenStreetMap & Nominatim.
      </div>
    </footer>
  );
}
