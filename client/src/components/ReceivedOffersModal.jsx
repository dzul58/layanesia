import React, { useState, useEffect } from 'react';
import { X, CheckCircle2, XCircle, Clock, User, ShieldCheck, AlertCircle } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export default function ReceivedOffersModal({ isOpen, onClose }) {
  const { token } = useAuth();
  const [offers, setOffers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [actionLoadingId, setActionLoadingId] = useState(null);
  const [errorMsg, setErrorMsg] = useState('');

  const fetchReceivedOffers = async () => {
    if (!token) return;
    setLoading(true);
    setErrorMsg('');
    try {
      const data = await api('/api/v1/offers/received', { token });
      setOffers(data.data || []);
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchReceivedOffers();
    }
  }, [isOpen, token]);

  const handleRespond = async (offerId, status) => {
    setActionLoadingId(offerId);
    setErrorMsg('');
    try {
      const data = await api(`/api/v1/offers/${offerId}/respond`, {
        method: 'PATCH',
        token,
        body: { status },
      });
      alert(data.message);
      fetchReceivedOffers();
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setActionLoadingId(null);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '680px', width: '90%', maxHeight: '85vh', display: 'flex', flexDirection: 'column' }} onClick={e => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(16, 185, 129, 0.15)', color: '#10b981' }}>
              <Clock size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.25rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Daftar Tawaran Pekerjaan Masuk</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Model 1 - Penawaran kerja langsung dari Pemberi Kerja</p>
            </div>
          </div>

          <button onClick={onClose} style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}>
            <X size={20} />
          </button>
        </div>

        {errorMsg && (
          <div style={{ background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', borderRadius: '10px', padding: '12px 16px', marginBottom: '16px', color: '#ef4444', fontSize: '0.88rem', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <AlertCircle size={20} />
            <div>{errorMsg}</div>
          </div>
        )}

        <div style={{ flex: 1, overflowY: 'auto', paddingRight: '4px' }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Memuat tawaran masuk...</div>
          ) : offers.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>
              <Clock size={40} style={{ opacity: 0.3, marginBottom: '12px' }} />
              <p style={{ margin: 0 }}>Belum ada tawaran pekerjaan masuk untuk profil Anda.</p>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
              {offers.map(offer => {
                const employerName = offer.employer?.name || 'Pemberi Kerja';
                const skillTitle = offer.skill_posting?.title || 'Keahlian';
                const isKtp = offer.employer?.is_ktp_verified;
                const formattedDate = offer.work_date ? offer.work_date.split('T')[0] : '-';

                return (
                  <div key={offer.id} className="glass-panel" style={{ padding: '18px', border: '1px solid var(--border-color)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                      <div>
                        <span className="badge badge-seeker" style={{ marginBottom: '6px', display: 'inline-block' }}>{skillTitle}</span>
                        <h4 style={{ margin: 0, fontSize: '1.05rem', color: 'var(--text-main)', display: 'flex', alignItems: 'center', gap: '6px' }}>
                          <User size={16} color="#6366f1" /> {employerName}
                          {isKtp && <ShieldCheck size={16} color="#10b981" title="KTP Verified" />}
                        </h4>
                      </div>

                      <span style={{
                        padding: '4px 10px',
                        borderRadius: '20px',
                        fontSize: '0.75rem',
                        fontWeight: 700,
                        background: offer.status === 'ACCEPTED' ? 'rgba(16,185,129,0.15)' : offer.status === 'REJECTED' ? 'rgba(239,68,68,0.15)' : 'rgba(245,158,11,0.15)',
                        color: offer.status === 'ACCEPTED' ? '#10b981' : offer.status === 'REJECTED' ? '#ef4444' : '#f59e0b',
                        border: `1px solid ${offer.status === 'ACCEPTED' ? 'rgba(16,185,129,0.4)' : offer.status === 'REJECTED' ? 'rgba(239,68,68,0.4)' : 'rgba(245,158,11,0.4)'}`
                      }}>
                        {offer.status === 'ACCEPTED' ? 'DITERIMA (TERHUBUNG)' : offer.status === 'REJECTED' ? 'DITOLAK' : 'PENDING (MENUNGGU RESPON)'}
                      </span>
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', background: 'var(--bg-secondary)', padding: '12px', borderRadius: '8px', marginBottom: '14px', fontSize: '0.85rem' }}>
                      <div>
                        <span style={{ color: 'var(--text-muted)', display: 'block', fontSize: '0.78rem' }}>Honor Ditawarkan:</span>
                        <strong style={{ color: '#10b981', fontSize: '1.1rem' }}>Rp {offer.offered_budget?.toLocaleString('id-ID')}</strong>
                      </div>
                      <div>
                        <span style={{ color: 'var(--text-muted)', display: 'block', fontSize: '0.78rem' }}>Tanggal Kerja:</span>
                        <strong style={{ color: 'var(--text-main)' }}>📅 {formattedDate}</strong>
                      </div>
                    </div>

                    {offer.status === 'PENDING' && (
                      <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
                        <button
                          type="button"
                          className="btn btn-outline"
                          style={{ borderColor: 'rgba(239,68,68,0.4)', color: '#ef4444', padding: '6px 14px', fontSize: '0.82rem' }}
                          onClick={() => handleRespond(offer.id, 'REJECTED')}
                          disabled={actionLoadingId === offer.id}
                        >
                          <XCircle size={14} /> Tolak Tawaran
                        </button>
                        <button
                          type="button"
                          className="btn btn-primary"
                          style={{ padding: '6px 16px', fontSize: '0.82rem' }}
                          onClick={() => handleRespond(offer.id, 'ACCEPTED')}
                          disabled={actionLoadingId === offer.id}
                        >
                          <CheckCircle2 size={14} /> Terima (Buka Kontak)
                        </button>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
