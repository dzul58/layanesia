import React, { useState } from 'react';
import { X, Zap, AlertCircle } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';
import LocationSelector from './LocationSelector';

export default function PostSkillModal({ isOpen, onClose, onSuccess, onOpenSubModal }) {
  const { token, user } = useAuth();

  const [category, setCategory] = useState('SERABUTAN');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [rateType, setRateType] = useState('PER_HOUR');
  const [rateAmount, setRateAmount] = useState('');
  const [useProfileLocation, setUseProfileLocation] = useState(true);
  const [locationState, setLocationState] = useState({
    province: user?.province || '',
    city: user?.city || '',
    district: user?.district || ''
  });

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');
  const [showSubPrompt, setShowSubPrompt] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token) {
      setErrorMsg('Anda harus masuk terlebih dahulu untuk mempublikasikan keahlian.');
      return;
    }

    if (!title.trim() || !description.trim() || !rateAmount) {
      setErrorMsg('Mohon lengkapi judul, deskripsi, dan tarif keahlian Anda.');
      return;
    }

    setLoading(true);
    setErrorMsg('');
    setShowSubPrompt(false);

    try {
      const payload = {
        category,
        title,
        description,
        rate_type: rateType,
        rate_amount: parseFloat(rateAmount),
        availability: 'AVAILABLE',
      };

      if (!useProfileLocation) {
        payload.province = locationState.province;
        payload.city = locationState.city;
        payload.district = locationState.district;
      }

      const data = await api('/api/v1/skills', {
        method: 'POST',
        token,
        body: payload,
      });

      alert(data.message || 'Berhasil! Profil keahlian Anda kini aktif di Etalase Talent Layanesia.');
      if (onSuccess) onSuccess(data.data);
      onClose();
    } catch (err) {
      if (err.code === 'SUBSCRIPTION_REQUIRED') {
        setShowSubPrompt(true);
        setErrorMsg('Anda belum berlangganan Paket Posting Keahlian (Rp 10.000 / Bulan).');
      } else {
        setErrorMsg(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '640px', width: '90%' }} onClick={e => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(16, 185, 129, 0.15)', color: '#10b981' }}>
              <Zap size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.25rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Pasang Profil Keahlian (Talent Showcase)</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Model 1 - Penawaran kerja langsung dari Pemberi Kerja</p>
            </div>
          </div>

          <button onClick={onClose} style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}>
            <X size={20} />
          </button>
        </div>

        {errorMsg && (
          <div style={{ background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', borderRadius: '10px', padding: '12px 16px', marginBottom: '18px', color: '#ef4444', fontSize: '0.88rem', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <AlertCircle size={20} />
            <div style={{ flex: 1 }}>{errorMsg}</div>
          </div>
        )}

        {showSubPrompt && (
          <div style={{ background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.4)', borderRadius: '12px', padding: '16px', marginBottom: '20px', textAlign: 'center' }}>
            <h4 style={{ margin: '0 0 6px 0', color: '#10b981', fontSize: '1rem' }}>Aktifkan Paket Posting Keahlian (Rp 10.000 / 30 Hari)</h4>
            <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '14px' }}>
              Untuk mempublikasikan keahlian Anda di etalase, aktifkan paket 10K yang berlaku selama 30 hari penuh.
            </p>
            <button 
              type="button" 
              className="btn btn-primary"
              onClick={() => {
                onClose();
                if (onOpenSubModal) onOpenSubModal('POST_SKILL_10K');
              }}
              style={{ padding: '8px 20px', fontSize: '0.88rem' }}
            >
              Langganan Rp 10.000 Sekarang
            </button>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          {/* Category Toggle */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '8px' }}>Kategori Keahlian</label>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
              <button
                type="button"
                onClick={() => setCategory('SERABUTAN')}
                style={{
                  padding: '12px',
                  borderRadius: '10px',
                  border: category === 'SERABUTAN' ? '2px solid #10b981' : '1px solid var(--border-color)',
                  background: category === 'SERABUTAN' ? 'rgba(16, 185, 129, 0.12)' : 'var(--bg-secondary)',
                  color: 'var(--text-main)',
                  fontWeight: 600,
                  fontSize: '0.9rem',
                  cursor: 'pointer',
                  textAlign: 'center',
                  transition: 'all 0.2s ease'
                }}
              >
                🛠️ Pekerjaan Serabutan (Harian/Jam)
              </button>

              <button
                type="button"
                onClick={() => setCategory('PROFESIONAL')}
                style={{
                  padding: '12px',
                  borderRadius: '10px',
                  border: category === 'PROFESIONAL' ? '2px solid #6366f1' : '1px solid var(--border-color)',
                  background: category === 'PROFESIONAL' ? 'rgba(99, 102, 241, 0.12)' : 'var(--bg-secondary)',
                  color: 'var(--text-main)',
                  fontWeight: 600,
                  fontSize: '0.9rem',
                  cursor: 'pointer',
                  textAlign: 'center',
                  transition: 'all 0.2s ease'
                }}
              >
                💼 Pekerjaan Profesional / Pengganti (Coverage)
              </button>
            </div>
          </div>

          {/* Title Input */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Judul Keahlian / Jasa</label>
            <input
              type="text"
              className="input-field"
              placeholder="Contoh: Supir Pengganti Cuti / Cuci Setrika Harian / Staff Admin Logistik"
              value={title}
              onChange={e => setTitle(e.target.value)}
              required
            />
          </div>

          {/* Description Textarea */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Deskripsi Lengkap & Pengalaman</label>
            <textarea
              className="input-field"
              rows={3}
              placeholder="Jelaskan keterampilan, sertifikat/SIM yang dimiliki, dan ketersediaan waktu Anda..."
              value={description}
              onChange={e => setDescription(e.target.value)}
              required
            />
          </div>

          {/* Rate Type & Amount */}
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '12px', marginBottom: '18px' }}>
            <div>
              <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Satuan Tarif</label>
              <select 
                className="input-field"
                value={rateType}
                onChange={e => setRateType(e.target.value)}
              >
                <option value="PER_HOUR">Per Jam</option>
                <option value="PER_DAY">Per Hari</option>
              </select>
            </div>

            <div>
              <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Besar Tarif (Rp)</label>
              <input
                type="number"
                className="input-field"
                placeholder="Contoh: 150000"
                value={rateAmount}
                onChange={e => setRateAmount(e.target.value)}
                min="10000"
                required
              />
            </div>
          </div>

          {/* Location Selection */}
          <div style={{ marginBottom: '22px' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '8px' }}>
              <label style={{ fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)' }}>Lokasi Operasional Jasa</label>
              <label style={{ fontSize: '0.8rem', color: '#10b981', display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer' }}>
                <input 
                  type="checkbox"
                  checked={useProfileLocation}
                  onChange={e => setUseProfileLocation(e.target.checked)}
                />
                Gunakan Lokasi Profil Domisili ({user?.district || 'Profil'})
              </label>
            </div>

            {!useProfileLocation && (
              <div style={{ background: 'var(--bg-secondary)', padding: '14px', borderRadius: '10px', border: '1px solid var(--border-color)' }}>
                <LocationSelector onLocationChange={loc => setLocationState(loc)} />
              </div>
            )}
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
            <button type="button" className="btn btn-outline" onClick={onClose} disabled={loading}>
              Batal
            </button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Memproses...' : 'Publikasikan Keahlian'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
