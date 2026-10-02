import React, { useState } from 'react';
import { X, Star, AlertCircle } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export default function RatingModal({ isOpen, onClose, connectionId, partnerName, onSuccess }) {
  const { token } = useAuth();
  const [ratingStars, setRatingStars] = useState(5);
  const [comment, setComment] = useState('');
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  if (!isOpen) return null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token || !connectionId) return;

    setLoading(true);
    setErrorMsg('');

    try {
      const data = await api('/api/v1/ratings', {
        method: 'POST',
        token,
        body: {
          connection_id: connectionId,
          rating_stars: ratingStars,
          comment: comment.trim() ? comment.trim() : null,
        },
      });

      alert(data.message || 'Ulasan tersimpan.');
      if (onSuccess) onSuccess(data.data);
      onClose();
    } catch (err) {
      if (err.code === 'CONNECTION_NOT_COMPLETED') {
        setErrorMsg('Tandai pekerjaan selesai di Kontak Terhubung sebelum memberi ulasan.');
      } else {
        setErrorMsg(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '480px', width: '90%' }} onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(245, 158, 11, 0.15)', color: '#f59e0b' }}>
              <Star size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.2rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Rating & Ulasan Dua Arah</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Evaluasi pengalaman kerja dengan {partnerName || 'Partner'}</p>
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
          <div style={{ textAlign: 'center', marginBottom: '20px' }}>
            <label style={{ display: 'block', fontSize: '0.88rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '10px' }}>
              Beri Bintang (1 - 5 Bintang):
            </label>
            <div style={{ display: 'flex', justifyContent: 'center', gap: '8px' }}>
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  onClick={() => setRatingStars(star)}
                  style={{
                    background: 'none',
                    border: 'none',
                    cursor: 'pointer',
                    transform: ratingStars >= star ? 'scale(1.15)' : 'scale(1)',
                    transition: 'transform 0.15s ease',
                  }}
                >
                  <Star
                    size={32}
                    color="#f59e0b"
                    fill={ratingStars >= star ? '#f59e0b' : 'transparent'}
                  />
                </button>
              ))}
            </div>
          </div>

          <div style={{ marginBottom: '22px' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '6px' }}>Ulasan & Catatan Pengalaman Kerja</label>
            <textarea
              className="input-field"
              rows={3}
              placeholder="Contoh: Sangat tepat waktu, ramah, dan hasil kerja sangat memuaskan."
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
            <button type="button" className="btn btn-outline" onClick={onClose} disabled={loading}>
              Batal
            </button>
            <button type="submit" className="btn btn-primary" style={{ background: 'var(--gradient-primary)' }} disabled={loading}>
              {loading ? 'Menyimpan...' : 'Kirim Rating & Ulasan'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
