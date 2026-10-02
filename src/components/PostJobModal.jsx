import React, { useState } from 'react';
import { X, Briefcase, AlertCircle } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';
import LocationSelector from './LocationSelector';

export default function PostJobModal({ isOpen, onClose, onSuccess }) {
  const { token, user } = useAuth();

  const [category, setCategory] = useState('SERABUTAN');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [budget, setBudget] = useState('');
  const [workDate, setWorkDate] = useState(() => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    return tomorrow.toISOString().split('T')[0];
  });
  const [durationType, setDurationType] = useState('HOURS');
  const [durationValue, setDurationValue] = useState(4);
  const [useProfileLocation, setUseProfileLocation] = useState(true);
  const [locationState, setLocationState] = useState({
    province: user?.province || '',
    city: user?.city || '',
    district: user?.district || ''
  });

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  if (!isOpen) return null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token) {
      setErrorMsg('Anda harus masuk terlebih dahulu sebagai Pemberi Kerja.');
      return;
    }

    if (!title.trim() || !description.trim() || !budget) {
      setErrorMsg('Lengkapi judul, deskripsi, dan total budget honor pekerjaan.');
      return;
    }

    setLoading(true);
    setErrorMsg('');

    try {
      const payload = {
        category,
        title,
        description,
        budget: parseFloat(budget),
        work_date: workDate,
        duration_type: durationType,
        duration_value: parseInt(durationValue, 10),
      };

      if (!useProfileLocation) {
        payload.province = locationState.province;
        payload.city = locationState.city;
        payload.district = locationState.district;
      }

      const data = await api('/api/v1/jobs', {
        method: 'POST',
        token,
        body: payload,
      });

      alert(data.message || 'Berhasil! Lowongan pekerjaan insidental Anda telah dipublikasikan secara GRATIS.');
      if (onSuccess) onSuccess(data.data);
      onClose();
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '640px', width: '90%' }} onClick={e => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(99, 102, 241, 0.15)', color: '#6366f1' }}>
              <Briefcase size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.25rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Buka Lowongan Pekerjaan Insidental (Model 2)</h2>
              <p style={{ fontSize: '0.8rem', color: '#10b981', margin: 0, fontWeight: 600 }}>🎉 Post Lowongan 100% GRATIS untuk Pemberi Kerja</p>
            </div>
          </div>

          <button onClick={onClose} style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}>
            <X size={20} />
          </button>
        </div>

        {errorMsg && (
          <div style={{ background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', borderRadius: '10px', padding: '12px 16px', marginBottom: '18px', color: '#ef4444', fontSize: '0.88rem', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <AlertCircle size={20} />
            <div>{errorMsg}</div>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          {/* Category Selector */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '8px' }}>Kategori Lowongan</label>
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
                  textAlign: 'center'
                }}
              >
                🛠️ Pekerjaan Serabutan
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
                  textAlign: 'center'
                }}
              >
                💼 Pekerjaan Profesional / Temporary
              </button>
            </div>
          </div>

          {/* Title */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Judul Lowongan Pekerjaan</label>
            <input
              type="text"
              className="input-field"
              placeholder="Contoh: Dibutuhkan Supir Boks Pengganti 2 Hari / Bantuan ART Kebun"
              value={title}
              onChange={e => setTitle(e.target.value)}
              required
            />
          </div>

          {/* Description */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Deskripsi Pekerjaan & Persyaratan</label>
            <textarea
              className="input-field"
              rows={3}
              placeholder="Jelaskan rincian tugas, jam kerja, dan persyaratan tenaga kerja yang dibutuhkan..."
              value={description}
              onChange={e => setDescription(e.target.value)}
              required
            />
          </div>

          {/* Budget, Work Date, Duration */}
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '12px', marginBottom: '18px' }}>
            <div>
              <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Total Budget Honor (Rp)</label>
              <input
                type="number"
                className="input-field"
                placeholder="250000"
                value={budget}
                onChange={e => setBudget(e.target.value)}
                required
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Tanggal Pelaksanaan</label>
              <input
                type="date"
                className="input-field"
                value={workDate}
                onChange={e => setWorkDate(e.target.value)}
                min={new Date().toISOString().split('T')[0]}
                required
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Durasi Kerja</label>
              <div style={{ display: 'flex', gap: '6px' }}>
                <input
                  type="number"
                  className="input-field"
                  style={{ width: '60px' }}
                  value={durationValue}
                  onChange={e => setDurationValue(e.target.value)}
                  min="1"
                  required
                />
                <select className="input-field" value={durationType} onChange={e => setDurationType(e.target.value)}>
                  <option value="HOURS">Jam</option>
                  <option value="DAYS">Hari</option>
                </select>
              </div>
            </div>
          </div>

          {/* Location */}
          <div style={{ marginBottom: '22px' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '8px' }}>
              <label style={{ fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)' }}>Lokasi Pelaksanaan Pekerjaan</label>
              <label style={{ fontSize: '0.8rem', color: '#10b981', display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer' }}>
                <input 
                  type="checkbox"
                  checked={useProfileLocation}
                  onChange={e => setUseProfileLocation(e.target.checked)}
                />
                Gunakan Lokasi Profil ({user?.district || 'Profil'})
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
            <button type="submit" className="btn btn-primary" style={{ background: 'var(--gradient-secondary)' }} disabled={loading}>
              {loading ? 'Mempublikasikan...' : 'Publikasikan Lowongan (Gratis)'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
