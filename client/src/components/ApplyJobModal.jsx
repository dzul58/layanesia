import React, { useState } from 'react';
import { X, Send, ShieldCheck, FileText, Zap, AlertCircle, CheckCircle2 } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export default function ApplyJobModal({ isOpen, onClose, job, onSuccess, onOpenProfileModal, onOpenSubModal }) {
  const { token, user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');
  const [promptType, setPromptType] = useState(null);

  if (!isOpen || !job) return null;

  const employerName = job.employerName || job.employer?.name || 'Pemberi Kerja';
  const jobTitle = job.title;
  const budgetStr = typeof job.budget === 'number' ? `Rp ${job.budget.toLocaleString('id-ID')}` : (job.budget || '-');
  const workDateStr = job.workDate || job.work_date?.split('T')[0] || '-';

  const ktpPending = user?.ktp_status === 'PENDING_REVIEW';
  const ktpOk = user?.is_ktp_verified;

  const handleApply = async () => {
    if (!token || !user) {
      setErrorMsg('Silakan masuk ke akun Anda terlebih dahulu untuk melamar pekerjaan.');
      return;
    }

    setLoading(true);
    setErrorMsg('');
    setPromptType(null);

    try {
      const data = await api(`/api/v1/jobs/${job.id}/apply`, {
        method: 'POST',
        token,
      });

      alert(data.message || 'Berhasil! Lamaran Anda telah terkirim.');
      if (onSuccess) onSuccess(data.data);
      onClose();
    } catch (err) {
      if (err.code === 'KTP_NOT_VERIFIED') {
        setPromptType(ktpPending ? 'KTP_PENDING' : 'KTP');
        setErrorMsg(ktpPending
          ? 'Pengajuan KTP masih ditinjau. Tunggu verifikasi admin sebelum melamar.'
          : 'Anda belum memverifikasi KTP. Verifikasi KTP diperlukan untuk keamanan transaksi pekerjaan.');
      } else if (err.code === 'RESUME_NOT_UPLOADED') {
        setPromptType('RESUME');
        setErrorMsg('Anda belum mengunggah Resume/CV. File Resume diperlukan agar Pemberi Kerja dapat melihat kualifikasi Anda.');
      } else if (err.code === 'SUBSCRIPTION_REQUIRED') {
        setPromptType('SUB_5K');
        setErrorMsg('Anda wajib berlangganan Paket Melamar (Rp 5.000 / Bulan) yang aktif.');
      } else {
        setErrorMsg(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content glass-panel" style={{ maxWidth: '520px', width: '90%' }} onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', borderBottom: '1px solid var(--border-color)', paddingBottom: '14px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <div style={{ padding: '8px', borderRadius: '10px', background: 'rgba(16, 185, 129, 0.15)', color: '#10b981' }}>
              <Send size={22} />
            </div>
            <div>
              <h2 style={{ fontSize: '1.2rem', fontWeight: 700, margin: 0, color: 'var(--text-main)' }}>Melamar Lowongan Pekerjaan</h2>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: 0 }}>Model 2 - Demand Application-Driven</p>
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

        {(promptType === 'KTP' || promptType === 'KTP_PENDING') && (
          <div style={{ background: 'rgba(245, 158, 11, 0.12)', border: '1px solid rgba(245, 158, 11, 0.4)', borderRadius: '12px', padding: '14px', marginBottom: '20px', textAlign: 'center' }}>
            <h4 style={{ margin: '0 0 6px 0', color: '#f59e0b', fontSize: '0.95rem' }}>
              {promptType === 'KTP_PENDING' ? 'KTP Sedang Ditinjau' : 'Verifikasi KTP Diperlukan'}
            </h4>
            <p style={{ fontSize: '0.82rem', color: 'var(--text-muted)', marginBottom: '12px' }}>
              {promptType === 'KTP_PENDING'
                ? 'Admin belum menyelesaikan tinjauan KTP Anda.'
                : 'Buka profil Anda untuk mengunggah foto & nomor KTP.'}
            </p>
            <button
              className="btn btn-outline"
              onClick={() => { onClose(); if (onOpenProfileModal) onOpenProfileModal(); }}
              style={{ borderColor: '#f59e0b', color: '#f59e0b', padding: '6px 16px', fontSize: '0.82rem' }}
            >
              <ShieldCheck size={14} /> Buka Verifikasi KTP Profil
            </button>
          </div>
        )}

        {promptType === 'RESUME' && (
          <div style={{ background: 'rgba(99, 102, 241, 0.12)', border: '1px solid rgba(99, 102, 241, 0.4)', borderRadius: '12px', padding: '14px', marginBottom: '20px', textAlign: 'center' }}>
            <h4 style={{ margin: '0 0 6px 0', color: '#6366f1', fontSize: '0.95rem' }}>Upload Resume/CV Diperlukan</h4>
            <p style={{ fontSize: '0.82rem', color: 'var(--text-muted)', marginBottom: '12px' }}>Unggah file Resume PDF Anda di menu profil pengguna.</p>
            <button
              className="btn btn-outline"
              onClick={() => { onClose(); if (onOpenProfileModal) onOpenProfileModal(); }}
              style={{ borderColor: '#6366f1', color: '#6366f1', padding: '6px 16px', fontSize: '0.82rem' }}
            >
              <FileText size={14} /> Buka Upload Resume Profil
            </button>
          </div>
        )}

        {promptType === 'SUB_5K' && (
          <div style={{ background: 'rgba(16, 185, 129, 0.12)', border: '1px solid rgba(16, 185, 129, 0.4)', borderRadius: '12px', padding: '14px', marginBottom: '20px', textAlign: 'center' }}>
            <h4 style={{ margin: '0 0 6px 0', color: '#10b981', fontSize: '0.95rem' }}>Aktifkan Paket Melamar (Rp 5.000 / Bulan)</h4>
            <p style={{ fontSize: '0.82rem', color: 'var(--text-muted)', marginBottom: '12px' }}>Langganan murah Rp 5.000 per bulan untuk melamar ke seluruh lowongan kerja.</p>
            <button
              className="btn btn-primary"
              onClick={() => { onClose(); if (onOpenSubModal) onOpenSubModal('APPLY_JOB_5K'); }}
              style={{ padding: '6px 16px', fontSize: '0.82rem' }}
            >
              <Zap size={14} /> Beli Paket Rp 5.000
            </button>
          </div>
        )}

        <div style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: '12px', padding: '16px', marginBottom: '20px' }}>
          <h3 style={{ fontSize: '1.1rem', margin: '0 0 6px 0', color: 'var(--text-main)' }}>{jobTitle}</h3>
          <p style={{ fontSize: '0.85rem', color: '#6366f1', fontWeight: 600, margin: '0 0 14px 0' }}>{employerName}</p>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px', fontSize: '0.85rem', color: 'var(--text-muted)' }}>
            <div>
              <span style={{ display: 'block', fontSize: '0.75rem' }}>Honor Pekerjaan:</span>
              <strong style={{ color: '#10b981', fontSize: '1.05rem' }}>{budgetStr}</strong>
            </div>
            <div>
              <span style={{ display: 'block', fontSize: '0.75rem' }}>Tanggal Pelaksanaan:</span>
              <strong style={{ color: 'var(--text-main)' }}>📅 {workDateStr}</strong>
            </div>
          </div>
        </div>

        <div style={{ background: 'var(--bg-glass)', border: '1px solid var(--border-color)', borderRadius: '10px', padding: '14px', marginBottom: '22px', fontSize: '0.82rem' }}>
          <div style={{ fontWeight: 600, color: 'var(--text-main)', marginBottom: '8px' }}>Cek Persyaratan Melamar:</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: ktpOk ? '#10b981' : '#ef4444' }}>
              <CheckCircle2 size={16} /> Verifikasi Identitas KTP: {ktpOk ? 'Terverifikasi' : ktpPending ? 'Menunggu tinjauan' : 'Belum'}
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: user?.resume_url ? '#10b981' : '#ef4444' }}>
              <CheckCircle2 size={16} /> File Resume/CV PDF: {user?.resume_url ? 'Sudah Diunggah' : 'Belum'}
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#10b981' }}>
              <CheckCircle2 size={16} /> Paket Langganan Melamar (Rp 5k/Bulan)
            </div>
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <button type="button" className="btn btn-outline" onClick={onClose} disabled={loading}>
            Batal
          </button>
          <button type="button" className="btn btn-primary" onClick={handleApply} disabled={loading}>
            {loading ? 'Mengirim Lamaran...' : 'Kirim Lamaran Pekerjaan'}
          </button>
        </div>
      </div>
    </div>
  );
}
