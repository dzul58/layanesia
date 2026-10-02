import React, { useState, useEffect } from 'react';
import { X, Users, ShieldCheck, FileText, CheckCircle2, XCircle, AlertCircle, ExternalLink } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api, fileUrl } from '../api';

export default function JobApplicantsModal({ isOpen, onClose, job }) {
  const { token } = useAuth();
  const [applicants, setApplicants] = useState([]);
  const [loading, setLoading] = useState(false);
  const [actionLoadingId, setActionLoadingId] = useState(null);
  const [errorMsg, setErrorMsg] = useState('');

  const fetchApplicants = async () => {
    if (!token || !job) return;
    setLoading(true);
    setErrorMsg('');
    try {
      const data = await api(`/api/v1/jobs/${job.id}/applicants`, { token });
      setApplicants(data.data || []);
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen && job) {
      fetchApplicants();
    }
  }, [isOpen, job, token]);

  const handleRespond = async (applicationId, action) => {
    setActionLoadingId(applicationId);
    setErrorMsg('');
    try {
      const data = await api(`/api/v1/jobs/applications/${applicationId}/${action}`, {
        method: 'PATCH',
        token,
      });
      alert(data.message);
      fetchApplicants();
    } catch (err) {
      setErrorMsg(err.message);
    } finally {
      setActionLoadingId(null);
    }
  };

  if (!isOpen || !job) return null;

  const statusMeta = (status) => {
    if (status === 'ACCEPTED') return { label: 'DITERIMA (TERHUBUNG)', bg: 'rgba(16,185,129,0.15)', color: '#10b981', border: 'rgba(16,185,129,0.4)' };
    if (status === 'REJECTED') return { label: 'DITOLAK', bg: 'rgba(239,68,68,0.15)', color: '#ef4444', border: 'rgba(239,68,68,0.4)' };
    return { label: 'MENUNGGU REVIU', bg: 'rgba(245,158,11,0.15)', color: '#f59e0b', border: 'rgba(245,158,11,0.4)' };
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '680px', width: '90%', maxHeight: '85vh', display: 'flex', flexDirection: 'column' }} onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(99, 102, 241, 0.15)', color: '#6366f1' }}>
              <Users size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.2rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Peninjauan Pelamar Lowongan</h2>
              <p style={{ fontSize: '0.8rem', color: '#10b981', margin: 0, fontWeight: 600 }}>{job.title}</p>
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
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Memuat data pelamar...</div>
          ) : applicants.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>
              <Users size={40} style={{ opacity: 0.3, marginBottom: '12px' }} />
              <p style={{ margin: 0 }}>Belum ada pelamar yang melamar pada lowongan ini.</p>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
              {applicants.map((app) => {
                const applicantName = app.applicant?.name || 'Pelamar';
                const isKtp = app.applicant?.is_ktp_verified;
                const resumeHref = fileUrl(app.applicant?.resume_url);
                const appliedDate = app.created_at ? app.created_at.split('T')[0] : '-';
                const meta = statusMeta(app.status);

                return (
                  <div key={app.id} className="glass-panel" style={{ padding: '18px', border: '1px solid var(--border-color)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <div style={{ width: '40px', height: '40px', borderRadius: '50%', background: 'var(--gradient-primary)', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold' }}>
                          {applicantName.charAt(0)}
                        </div>
                        <div>
                          <h4 style={{ margin: 0, fontSize: '1.05rem', color: 'var(--text-main)', display: 'flex', alignItems: 'center', gap: '6px' }}>
                            {applicantName}
                            {isKtp && <ShieldCheck size={16} color="#10b981" title="KTP Verified" />}
                          </h4>
                          <span style={{ fontSize: '0.78rem', color: 'var(--text-muted)' }}>Dilamar: {appliedDate}</span>
                        </div>
                      </div>

                      <span style={{
                        padding: '4px 10px',
                        borderRadius: '20px',
                        fontSize: '0.75rem',
                        fontWeight: 700,
                        background: meta.bg,
                        color: meta.color,
                        border: `1px solid ${meta.border}`,
                      }}>
                        {meta.label}
                      </span>
                    </div>

                    <div style={{ display: 'flex', gap: '16px', background: 'var(--bg-secondary)', padding: '12px', borderRadius: '8px', marginBottom: '14px', fontSize: '0.85rem' }}>
                      {resumeHref && (
                        <a
                          href={resumeHref}
                          target="_blank"
                          rel="noreferrer"
                          style={{ color: '#6366f1', textDecoration: 'none', display: 'flex', alignItems: 'center', gap: '4px', fontWeight: 600 }}
                        >
                          <FileText size={16} /> Lihat Resume / CV (PDF) <ExternalLink size={12} />
                        </a>
                      )}
                    </div>

                    {app.status === 'PENDING' && (
                      <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
                        <button
                          type="button"
                          className="btn btn-outline"
                          style={{ borderColor: 'rgba(239,68,68,0.4)', color: '#ef4444', padding: '6px 14px', fontSize: '0.82rem' }}
                          onClick={() => handleRespond(app.id, 'reject')}
                          disabled={actionLoadingId === app.id}
                        >
                          <XCircle size={14} /> Tolak
                        </button>
                        <button
                          type="button"
                          className="btn btn-primary"
                          style={{ padding: '6px 16px', fontSize: '0.82rem' }}
                          onClick={() => handleRespond(app.id, 'accept')}
                          disabled={actionLoadingId === app.id}
                        >
                          <CheckCircle2 size={14} /> Terima Pelamar (Buka Kontak)
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
