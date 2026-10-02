import React, { useState, useEffect } from 'react';
import { X, CheckCircle2, MessageSquare, MapPin, Star, AlertCircle, ShieldCheck } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api, whatsappUrl } from '../api';

export default function ActiveConnectionsModal({ isOpen, onClose, onOpenRatingModal }) {
  const { token } = useAuth();
  const [connections, setConnections] = useState([]);
  const [loading, setLoading] = useState(false);
  const [actionId, setActionId] = useState(null);
  const [errorMsg, setErrorMsg] = useState('');

  const fetchConnections = async () => {
    if (!token) return;
    setLoading(true);
    setErrorMsg('');
    try {
      const data = await api('/api/v1/connections', { token });
      setConnections(data.data || []);
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchConnections();
    }
  }, [isOpen, token]);

  const handleComplete = async (connectionId) => {
    setActionId(connectionId);
    setErrorMsg('');
    try {
      const data = await api(`/api/v1/connections/${connectionId}/complete`, {
        method: 'PATCH',
        token,
      });
      alert(data.message || 'Pekerjaan ditandai selesai.');
      await fetchConnections();
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setActionId(null);
    }
  };

  if (!isOpen) return null;

  const active = connections.filter((item) => item.connection?.status === 'ACTIVE');
  const completed = connections.filter((item) => item.connection?.status === 'COMPLETED');

  const renderCard = (item) => {
    const conn = item.connection;
    const contact = item.contact;
    if (!conn || !contact) return null;
    const isActive = conn.status === 'ACTIVE';
    const waHref = whatsappUrl(contact.whatsapp_link || contact.phone);

    return (
      <div
        key={conn.id}
        className="glass-panel"
        style={{
          padding: '18px',
          border: isActive ? '1px solid rgba(16, 185, 129, 0.4)' : '1px solid var(--border-color)',
          background: isActive ? 'rgba(16, 185, 129, 0.04)' : 'transparent',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ width: '42px', height: '42px', borderRadius: '50%', background: 'var(--gradient-primary)', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold' }}>
              {(contact.name || '?').charAt(0)}
            </div>
            <div>
              <h4 style={{ margin: 0, fontSize: '1.05rem', color: 'var(--text-main)', display: 'flex', alignItems: 'center', gap: '6px' }}>
                {contact.name}
                {contact.is_ktp_verified && <ShieldCheck size={16} color="#10b981" title="KTP Verified" />}
              </h4>
              <span style={{ fontSize: '0.78rem', color: '#10b981', fontWeight: 600 }}>
                {conn.source_type === 'JOB_OFFER' ? 'Tawaran Pekerjaan Disetujui' : 'Lamaran Pekerjaan Diterima'}
              </span>
            </div>
          </div>

          <span className="badge badge-seeker" style={{ background: isActive ? 'rgba(16, 185, 129, 0.2)' : 'rgba(99, 102, 241, 0.2)', color: isActive ? '#10b981' : '#818cf8' }}>
            {isActive ? 'TERHUBUNG' : 'SELESAI'}
          </span>
        </div>

        <div style={{ background: 'var(--bg-secondary)', padding: '12px', borderRadius: '8px', marginBottom: '14px', fontSize: '0.85rem' }}>
          <div style={{ marginBottom: '6px', color: 'var(--text-main)' }}>
            <strong>Nomor Telepon/HP:</strong> {contact.phone}
          </div>
          <div style={{ color: 'var(--text-muted)' }}>
            <MapPin size={14} color="#10b981" style={{ verticalAlign: 'middle', marginRight: '4px' }} />
            {contact.district}, {contact.city}, {contact.province} {contact.address_detail ? `(${contact.address_detail})` : ''}
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '8px' }}>
          <a
            href={waHref}
            target="_blank"
            rel="noreferrer"
            className="btn btn-primary"
            style={{ padding: '6px 14px', fontSize: '0.82rem', background: '#25D366', borderColor: '#25D366' }}
          >
            <MessageSquare size={14} /> Chat via WhatsApp Direct
          </a>

          {isActive ? (
            <button
              type="button"
              className="btn btn-outline"
              style={{ padding: '6px 14px', fontSize: '0.82rem' }}
              disabled={actionId === conn.id}
              onClick={() => handleComplete(conn.id)}
            >
              <CheckCircle2 size={14} /> {actionId === conn.id ? 'Menyimpan...' : 'Tandai Selesai'}
            </button>
          ) : (
            <button
              type="button"
              className="btn btn-outline"
              style={{ padding: '6px 14px', fontSize: '0.82rem', borderColor: '#f59e0b', color: '#f59e0b' }}
              onClick={() => {
                onClose();
                if (onOpenRatingModal) onOpenRatingModal(conn.id, contact.name);
              }}
            >
              <Star size={14} /> Beri Rating & Review
            </button>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '680px', width: '90%', maxHeight: '85vh', display: 'flex', flexDirection: 'column' }} onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(16, 185, 129, 0.15)', color: '#10b981' }}>
              <CheckCircle2 size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.25rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Kontak Terhubung</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Tandai selesai dulu, baru beri rating. Hubungi partner via WhatsApp.</p>
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
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Memuat kontak terhubung...</div>
          ) : connections.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>
              <CheckCircle2 size={40} style={{ opacity: 0.3, marginBottom: '12px' }} />
              <p style={{ margin: 0 }}>Belum ada kontak terhubung saat ini.</p>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '18px' }}>
              {active.length > 0 && (
                <div>
                  <h3 style={{ fontSize: '0.85rem', color: 'var(--text-muted)', margin: '0 0 10px 0' }}>Sedang berjalan</h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
                    {active.map(renderCard)}
                  </div>
                </div>
              )}
              {completed.length > 0 && (
                <div>
                  <h3 style={{ fontSize: '0.85rem', color: 'var(--text-muted)', margin: '0 0 10px 0' }}>Selesai — siap diulas</h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
                    {completed.map(renderCard)}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
