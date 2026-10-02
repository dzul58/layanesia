import React from 'react';
import { MapPin, ShieldCheck, Heart } from 'lucide-react';

export default function Footer() {
  return (
    <footer style={{
      borderTop: '1px solid var(--border-color)',
      padding: '40px 24px',
      marginTop: '60px',
      background: 'var(--bg-secondary)',
      color: 'var(--text-muted)',
      fontSize: '0.9rem',
      transition: 'background-color 0.3s ease'
    }}>
      <div className="container" style={{ display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap', gap: '24px' }}>
        <div>
          <h3 style={{ color: 'var(--text-main)', fontSize: '1.2rem', marginBottom: '8px' }}>Layanesia</h3>
          <p style={{ maxWidth: '360px', color: 'var(--text-muted)' }}>
            Platform marketplace kerja sementara & insidental dua arah. Menghubungkan Pekerja Serabutan & Profesional dengan Pemberi Kerja secara transparan.
          </p>
        </div>

        <div>
          <h4 style={{ color: 'var(--text-main)', marginBottom: '12px', fontSize: '1rem' }}>Fitur & Model</h4>
          <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <li>Model 1: Talent Showcase (Worker Offer-Driven)</li>
            <li>Model 2: Job Board (Demand Apply-Driven - Gratis)</li>
            <li>Verifikasi KTP & Portfolio Cloudinary</li>
            <li>Peta Hiperlokal OpenStreetMap + Leaflet.js</li>
          </ul>
        </div>

        <div>
          <h4 style={{ color: 'var(--text-main)', marginBottom: '12px', fontSize: '1rem' }}>Monetisasi</h4>
          <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <li>Paket Apply Job (Rp 5.000 / bln)</li>
            <li>Paket Skill Showcase (Rp 10.000 / bln)</li>
            <li>Post Job Vacancy: <strong>100% GRATIS</strong></li>
          </ul>
        </div>
      </div>

      <div style={{ textAlign: 'center', marginTop: '32px', paddingTop: '20px', borderTop: '1px solid var(--border-color)', fontSize: '0.8rem', color: 'var(--text-dim)' }}>
        © 2026 Layanesia Platform. Open-source Geocoding by OpenStreetMap & Nominatim API.
      </div>
    </footer>
  );
}
