import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import LocationSelector from './LocationSelector';

const AuthModal = ({ isOpen, onClose, initialTab = 'login' }) => {
  const { login, register } = useAuth();
  const [tab, setTab] = useState(initialTab); // 'login' | 'register'
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Login State
  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');

  // Register State
  const [regName, setRegName] = useState('');
  const [regEmail, setRegEmail] = useState('');
  const [regPhone, setRegPhone] = useState('');
  const [regPassword, setRegPassword] = useState('');
  const [regActiveMode, setRegActiveMode] = useState('SEEKER');
  const [regLocation, setRegLocation] = useState({
    province: 'DKI Jakarta',
    city: 'Jakarta Selatan',
    district: 'Kebayoran Baru',
    addressDetail: '',
  });

  if (!isOpen) return null;

  const handleLoginSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(loginEmail, loginPassword);
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRegisterSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    if (!regLocation.province || !regLocation.city || !regLocation.district) {
      setError('Silakan pilih lokasi domisili lengkap (Provinsi, Kota/Kabupaten, Kecamatan).');
      setLoading(false);
      return;
    }

    try {
      await register({
        name: regName,
        email: regEmail,
        phone: regPhone,
        password: regPassword,
        active_mode: regActiveMode,
        province: regLocation.province,
        city: regLocation.city,
        district: regLocation.district,
        address_detail: regLocation.addressDetail || null,
      });
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content auth-modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close-btn" onClick={onClose}>&times;</button>
        
        <div className="auth-header">
          <div className="auth-logo">
            <span className="logo-icon">🤝</span>
            <h2>Layanesia</h2>
          </div>
          <p className="auth-subtitle">Marketplace Pekerja Harian & Jasa Insidental</p>
        </div>

        {/* Tab Switcher */}
        <div className="auth-tabs">
          <button
            type="button"
            className={`auth-tab ${tab === 'login' ? 'active' : ''}`}
            onClick={() => { setTab('login'); setError(''); }}
          >
            Masuk Akun
          </button>
          <button
            type="button"
            className={`auth-tab ${tab === 'register' ? 'active' : ''}`}
            onClick={() => { setTab('register'); setError(''); }}
          >
            Daftar Akun Baru
          </button>
        </div>

        {error && <div className="auth-error-banner">⚠️ {error}</div>}

        {/* LOGIN FORM */}
        {tab === 'login' && (
          <form onSubmit={handleLoginSubmit} className="auth-form">
            <div className="form-group">
              <label htmlFor="login-email">Alamat Email</label>
              <input
                id="login-email"
                type="email"
                placeholder="contoh: budi@gmail.com"
                value={loginEmail}
                onChange={(e) => setLoginEmail(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="login-password">Kata Sandi</label>
              <input
                id="login-password"
                type="password"
                placeholder="Masukkan kata sandi..."
                value={loginPassword}
                onChange={(e) => setLoginPassword(e.target.value)}
                required
              />
            </div>

            <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading}>
              {loading ? 'Memproses...' : 'Masuk Sekarang'}
            </button>
          </form>
        )}

        {/* REGISTER FORM */}
        {tab === 'register' && (
          <form onSubmit={handleRegisterSubmit} className="auth-form">
            <div className="form-group">
              <label htmlFor="reg-name">Nama Lengkap</label>
              <input
                id="reg-name"
                type="text"
                placeholder="Sesuai KTP"
                value={regName}
                onChange={(e) => setRegName(e.target.value)}
                required
              />
            </div>

            <div className="form-row">
              <div className="form-group flex-1">
                <label htmlFor="reg-email">Alamat Email</label>
                <input
                  id="reg-email"
                  type="email"
                  placeholder="budi@gmail.com"
                  value={regEmail}
                  onChange={(e) => setRegEmail(e.target.value)}
                  required
                />
              </div>

              <div className="form-group flex-1">
                <label htmlFor="reg-phone">Nomor Whatsapp / HP</label>
                <input
                  id="reg-phone"
                  type="tel"
                  placeholder="081234567890"
                  value={regPhone}
                  onChange={(e) => setRegPhone(e.target.value)}
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="reg-password">Kata Sandi</label>
              <input
                id="reg-password"
                type="password"
                placeholder="Minimal 6 karakter"
                value={regPassword}
                onChange={(e) => setRegPassword(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label>Peran Awal Utama Akun Anda</label>
              <div className="role-choice-grid">
                <button
                  type="button"
                  className={`role-choice-card ${regActiveMode === 'SEEKER' ? 'selected' : ''}`}
                  onClick={() => setRegActiveMode('SEEKER')}
                >
                  <span className="role-choice-icon">🧰</span>
                  <div>
                    <strong>Pencari Kerja (Seeker)</strong>
                    <p>Menawarkan keahlian & mencari lowongan kerja harian</p>
                  </div>
                </button>
                <button
                  type="button"
                  className={`role-choice-card ${regActiveMode === 'EMPLOYER' ? 'selected' : ''}`}
                  onClick={() => setRegActiveMode('EMPLOYER')}
                >
                  <span className="role-choice-icon">🏢</span>
                  <div>
                    <strong>Pemberi Kerja (Employer)</strong>
                    <p>Membuka lowongan harian & merekrut tenaga kerja</p>
                  </div>
                </button>
              </div>
              <small className="form-hint">💡 Anda dapat berganti peran kapan saja secara gratis lewat menu navbar.</small>
            </div>

            <div className="form-group">
              <label>Lokasi Domisili</label>
              <LocationSelector
                value={regLocation}
                onChange={setRegLocation}
                showMapButton={true}
              />
            </div>

            <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading}>
              {loading ? 'Mendaftarkan Akun...' : 'Daftar Sekarang'}
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

export default AuthModal;
