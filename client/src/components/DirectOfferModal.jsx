import React, { useState } from 'react';
import { X, Send, AlertCircle, ShieldCheck } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export default function DirectOfferModal({ isOpen, onClose, talent, onSuccess }) {
  const { token } = useAuth();
  const [offeredBudget, setOfferedBudget] = useState(talent?.rate_amount || talent?.rateAmount || '');
  const [workDate, setWorkDate] = useState(() => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    return tomorrow.toISOString().split('T')[0];
  });

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  if (!isOpen || !talent) return null;

  const talentName = talent.user?.name || talent.name || 'Pekerja';
  const skillTitle = talent.title;
  const talentRate = talent.rate_amount || talent.rateAmount;
  const rateTypeLabel = (talent.rate_type || talent.rateType) === 'PER_HOUR' ? 'per jam' : 'per hari';

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token) {
      setErrorMsg('Silakan login terlebih dahulu sebagai Pemberi Kerja untuk mengirim tawaran.');
      return;
    }

    if (!offeredBudget || parseFloat(offeredBudget) <= 0) {
      setErrorMsg('Tentukan nominal budget honor yang valid.');
      return;
    }

    setLoading(true);
    setErrorMsg('');

    try {
      const data = await api('/api/v1/offers', {
        method: 'POST',
        token,
        body: {
          skill_posting_id: talent.id,
          offered_budget: parseFloat(offeredBudget),
          work_date: workDate,
        },
      });

      alert(data.message || `Berhasil! Tawaran pekerjaan telah dikirimkan kepada ${talentName}.`);
      if (onSuccess) onSuccess(data.data);
      onClose();
    } catch (err) {
      if (err.code === 'KTP_NOT_VERIFIED') {
        setErrorMsg('Verifikasi KTP wajib sebelum mengirim tawaran. Ajukan KTP di profil Anda.');
      } else {
        setErrorMsg(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '520px', width: '90%' }} onClick={e => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(99, 102, 241, 0.15)', color: '#6366f1' }}>
              <Send size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.2rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Kirim Tawaran Pekerjaan (Direct Offer)</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Model 1 - Penawaran kerja langsung ke Pekerja</p>
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

        {/* Talent Info Card */}
        <div style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: '12px', padding: '16px', marginBottom: '20px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '8px' }}>
            <div style={{ width: '36px', height: '36px', borderRadius: '50%', background: 'var(--gradient-primary)', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold' }}>
              {talentName.charAt(0)}
            </div>
            <div>
              <h4 style={{ margin: 0, fontSize: '0.98rem', color: 'var(--text-main)', display: 'flex', alignItems: 'center', gap: '6px' }}>
                {talentName}
                {(talent.user?.is_ktp_verified || talent.isKtpVerified) && (
                  <ShieldCheck size={16} color="#10b981" title="KTP Verified" />
                )}
              </h4>
              <span style={{ fontSize: '0.8rem', color: '#10b981', fontWeight: 600 }}>{skillTitle}</span>
            </div>
          </div>

          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '10px', borderTop: '1px dashed var(--border-color)', fontSize: '0.85rem' }}>
            <span style={{ color: 'var(--text-muted)' }}>Tarif Dasar Pekerja:</span>
            <span style={{ fontWeight: 700, color: 'var(--text-main)' }}>
              Rp {typeof talentRate === 'number' ? talentRate.toLocaleString('id-ID') : talentRate} / {rateTypeLabel}
            </span>
          </div>
        </div>

        <form onSubmit={handleSubmit}>
          {/* Offered Budget Input */}
          <div style={{ marginBottom: '18px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Nominal Honor yang Ditawarkan (Rp)</label>
            <input
              type="number"
              className="input-field"
              placeholder="Masukkan total budget honor"
              value={offeredBudget}
              onChange={e => setOfferedBudget(e.target.value)}
              required
            />
          </div>

          {/* Work Date Picker */}
          <div style={{ marginBottom: '22px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Tanggal Pelaksanaan Pekerjaan</label>
            <input
              type="date"
              className="input-field"
              value={workDate}
              onChange={e => setWorkDate(e.target.value)}
              min={new Date().toISOString().split('T')[0]}
              required
            />
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
            <button type="button" className="btn btn-outline" onClick={onClose} disabled={loading}>
              Batal
            </button>
            <button type="submit" className="btn btn-primary" style={{ background: 'var(--gradient-secondary)' }} disabled={loading}>
              {loading ? 'Mengirim...' : 'Kirim Tawaran Sekarang'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
