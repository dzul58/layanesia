import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { api, fileUrl, ktpStatusLabel } from '../api';
import LocationSelector from './LocationSelector';
import SubscriptionModal from './SubscriptionModal';

const ktpBadgeClass = (user) => {
  if (user?.is_ktp_verified || user?.ktp_status === 'VERIFIED') return 'verified';
  if (user?.ktp_status === 'PENDING_REVIEW') return 'pending';
  if (user?.ktp_status === 'REJECTED') return 'rejected';
  return 'unverified';
};

const UserProfileModal = ({ isOpen, onClose }) => {
  const { user, token, updateProfile, verifyKTP, uploadResume, uploadFile } = useAuth();
  const [activeTab, setActiveTab] = useState('info');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  const [subscriptions, setSubscriptions] = useState([]);
  const [isSubModalOpen, setIsSubModalOpen] = useState(false);

  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [location, setLocation] = useState({
    province: '',
    city: '',
    district: '',
    addressDetail: '',
  });

  const [ktpNumber, setKtpNumber] = useState('');
  const [ktpFile, setKtpFile] = useState(null);
  const [ktpPreviewURL, setKtpPreviewURL] = useState('');

  const [resumeFile, setResumeFile] = useState(null);
  const [resumeURLInput, setResumeURLInput] = useState('');

  const fetchSubscriptions = async () => {
    if (!token) return;
    try {
      const data = await api('/api/v1/subscriptions/my', { token });
      setSubscriptions(data.subscriptions || []);
    } catch (err) {
      console.error('Error fetching subscriptions:', err);
    }
  };

  useEffect(() => {
    if (user) {
      setName(user.name || '');
      setPhone(user.phone || '');
      setLocation({
        province: user.province || 'DKI Jakarta',
        city: user.city || 'Jakarta Selatan',
        district: user.district || 'Kebayoran Baru',
        addressDetail: user.address_detail || '',
      });
      setKtpNumber(user.ktp_number || '');
      setKtpPreviewURL(user.ktp_image_url || '');
      setResumeURLInput(user.resume_url || '');
      fetchSubscriptions();
    }
  }, [user, isOpen]);

  if (!isOpen || !user) return null;

  const ktpPreviewSrc = ktpPreviewURL?.startsWith('blob:') ? ktpPreviewURL : fileUrl(ktpPreviewURL);
  const resumeHref = fileUrl(user.resume_url);
  const ktpLocked = user.is_ktp_verified || user.ktp_status === 'VERIFIED';

  const handleUpdateProfile = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      await updateProfile({
        name,
        phone,
        province: location.province,
        city: location.city,
        district: location.district,
        address_detail: location.addressDetail || null,
      });
      setSuccess('Profil berhasil diperbarui!');
      setTimeout(() => setActiveTab('info'), 1200);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyKTP = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      let finalImageUrl = ktpPreviewURL?.startsWith('blob:') ? '' : ktpPreviewURL;
      if (ktpFile) {
        finalImageUrl = await uploadFile(ktpFile, 'ktp');
      }

      if (!finalImageUrl) {
        throw new Error('Foto / File KTP wajib diunggah');
      }

      const data = await verifyKTP(ktpNumber, finalImageUrl);
      const status = data?.user?.ktp_status || data?.ktp_status;
      if (status === 'VERIFIED' || data?.user?.is_ktp_verified) {
        setSuccess('KTP berhasil diverifikasi!');
      } else if (status === 'PENDING_REVIEW') {
        setSuccess('Pengajuan KTP diterima dan sedang ditinjau.');
      } else {
        setSuccess(data?.message || 'Pengajuan KTP dikirim.');
      }
      setTimeout(() => setActiveTab('info'), 1400);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleUploadResume = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      let finalResumeUrl = resumeURLInput;
      if (resumeFile) {
        finalResumeUrl = await uploadFile(resumeFile, 'resume');
      }

      if (!finalResumeUrl) {
        throw new Error('File Resume PDF wajib diunggah');
      }

      await uploadResume(finalResumeUrl);
      setSuccess('Resume PDF berhasil disimpan!');
      setTimeout(() => setActiveTab('info'), 1200);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div className="modal-overlay" onClick={onClose}>
        <div className="modal-content profile-modal" onClick={(e) => e.stopPropagation()}>
          <button className="modal-close-btn" onClick={onClose}>&times;</button>

          <div className="profile-header">
            <div className="profile-avatar">
              {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
            </div>
            <div className="profile-header-info">
              <h3>{user.name}</h3>
              <p>{user.email} • {user.phone}</p>
              <div className="profile-badges">
                <span className={`badge-mode ${user.active_mode === 'SEEKER' ? 'badge-seeker' : 'badge-employer'}`}>
                  Mode: {user.active_mode === 'SEEKER' ? 'Pencari Kerja' : 'Pemberi Kerja'}
                </span>
                <span className={`badge-ktp ${ktpBadgeClass(user)}`}>
                  {user.is_ktp_verified ? '✓ KTP Terverifikasi' : `KTP: ${ktpStatusLabel(user)}`}
                </span>
              </div>
            </div>
          </div>

          <div className="profile-nav-tabs">
            <button
              className={`profile-nav-tab ${activeTab === 'info' ? 'active' : ''}`}
              onClick={() => { setActiveTab('info'); setError(''); setSuccess(''); }}
            >
              📋 Ringkasan
            </button>
            <button
              className={`profile-nav-tab ${activeTab === 'edit' ? 'active' : ''}`}
              onClick={() => { setActiveTab('edit'); setError(''); setSuccess(''); }}
            >
              ✏️ Edit Profil & Domisili
            </button>
            <button
              className={`profile-nav-tab ${activeTab === 'subscription' ? 'active' : ''}`}
              onClick={() => { setActiveTab('subscription'); setError(''); setSuccess(''); }}
            >
              💳 Paket Langganan
            </button>
            <button
              className={`profile-nav-tab ${activeTab === 'ktp' ? 'active' : ''}`}
              onClick={() => { setActiveTab('ktp'); setError(''); setSuccess(''); }}
            >
              🪪 Verifikasi KTP
            </button>
            <button
              className={`profile-nav-tab ${activeTab === 'resume' ? 'active' : ''}`}
              onClick={() => { setActiveTab('resume'); setError(''); setSuccess(''); }}
            >
              📄 Resume / CV
            </button>
          </div>

          {error && <div className="auth-error-banner">⚠️ {error}</div>}
          {success && <div className="auth-success-banner">✅ {success}</div>}

          {activeTab === 'info' && (
            <div className="profile-info-content">
              <div className="info-grid">
                <div className="info-item">
                  <span className="info-label">Nama Lengkap</span>
                  <span className="info-value">{user.name}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Email</span>
                  <span className="info-value">{user.email}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Nomor Telepon</span>
                  <span className="info-value">{user.phone}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Lokasi Domisili</span>
                  <span className="info-value">
                    📍 {user.district}, {user.city}, {user.province}
                  </span>
                </div>
                {user.address_detail && (
                  <div className="info-item full-width">
                    <span className="info-label">Detail Alamat Jalan</span>
                    <span className="info-value">{user.address_detail}</span>
                  </div>
                )}
                <div className="info-item">
                  <span className="info-label">Status Verifikasi KTP</span>
                  <span className="info-value">
                    {user.is_ktp_verified ? (
                      <span style={{ color: 'var(--success-color)', fontWeight: 600 }}>✓ Terverifikasi ({user.ktp_number})</span>
                    ) : user.ktp_status === 'PENDING_REVIEW' ? (
                      <span style={{ color: 'var(--warning-color)', fontWeight: 600 }}>Menunggu tinjauan admin</span>
                    ) : user.ktp_status === 'REJECTED' ? (
                      <span style={{ color: '#ef4444', fontWeight: 600 }}>
                        Ditolak{user.ktp_rejection_reason ? `: ${user.ktp_rejection_reason}` : ''}
                      </span>
                    ) : (
                      <span style={{ color: 'var(--warning-color)', fontWeight: 600 }}>Belum Diajukan</span>
                    )}
                  </span>
                </div>
                <div className="info-item">
                  <span className="info-label">File Resume PDF</span>
                  <span className="info-value">
                    {resumeHref ? (
                      <a href={resumeHref} target="_blank" rel="noopener noreferrer" className="link-action">
                        📄 Lihat Resume PDF
                      </a>
                    ) : (
                      <span className="text-muted">Belum Diunggah</span>
                    )}
                  </span>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'edit' && (
            <form onSubmit={handleUpdateProfile} className="auth-form">
              <div className="form-group">
                <label>Nama Lengkap</label>
                <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
              </div>

              <div className="form-group">
                <label>Nomor Whatsapp / HP</label>
                <input type="tel" value={phone} onChange={(e) => setPhone(e.target.value)} required />
              </div>

              <div className="form-group">
                <label>Lokasi Domisili</label>
                <LocationSelector
                  value={location}
                  onChange={(loc) => setLocation((prev) => ({ ...prev, ...loc }))}
                />
              </div>

              <div className="form-group">
                <label>Detail Alamat Jalan (opsional)</label>
                <input
                  type="text"
                  value={location.addressDetail}
                  onChange={(e) => setLocation((prev) => ({ ...prev, addressDetail: e.target.value }))}
                  placeholder="Nama jalan, nomor rumah, patokan"
                />
              </div>

              <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading}>
                {loading ? 'Simpan Perubahan...' : 'Simpan Profil'}
              </button>
            </form>
          )}

          {activeTab === 'subscription' && (
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                <p className="tab-description" style={{ margin: 0 }}>
                  Daftar Paket Langganan Aktif & Riwayat Pembayaran Midtrans.
                </p>
                <button
                  className="btn btn-primary"
                  style={{ padding: '6px 14px', fontSize: '0.82rem' }}
                  onClick={() => setIsSubModalOpen(true)}
                >
                  + Beli Paket Langganan
                </button>
              </div>

              {subscriptions.length === 0 ? (
                <div style={{ textOverflow: 'ellipsis', background: 'var(--input-bg)', padding: '24px', borderRadius: '12px', textAlign: 'center', color: 'var(--text-muted)' }}>
                  Belum ada paket langganan aktif. Klik tombol "+ Beli Paket Langganan" untuk melamar lowongan (5K) atau pasang keahlian (10K).
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  {subscriptions.map((sub) => (
                    <div key={sub.id} style={{ background: 'var(--input-bg)', padding: '14px', borderRadius: '12px', border: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <div>
                        <strong style={{ fontSize: '0.95rem', display: 'block', color: 'var(--text-main)' }}>
                          {sub.plan_type === 'APPLY_JOB_5K' ? 'Paket Melamar Lowongan (5K)' : 'Paket Etalase Keahlian Pekerja (10K)'}
                        </strong>
                        <span style={{ fontSize: '0.78rem', color: 'var(--text-muted)' }}>
                          {sub.status === 'ACTIVE' && sub.expires_at
                            ? `Aktif s/d ${new Date(sub.expires_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}`
                            : sub.status === 'PENDING'
                              ? 'Menunggu pembayaran Midtrans'
                              : `Status: ${sub.status}`}
                        </span>
                      </div>
                      <span className="badge badge-seeker">
                        Status: {sub.status}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {activeTab === 'ktp' && (
            <form onSubmit={handleVerifyKTP} className="auth-form">
              <p className="tab-description">
                Upload foto KTP resmi Anda. Pengajuan ditinjau admin (status PENDING_REVIEW) sebelum badge terverifikasi aktif.
              </p>
              {user.ktp_status === 'REJECTED' && user.ktp_rejection_reason && (
                <div className="auth-error-banner">Alasan penolakan: {user.ktp_rejection_reason}</div>
              )}

              <div className="form-group">
                <label>Nomor Induk Kependudukan (NIK 16 Digit)</label>
                <input
                  type="text"
                  placeholder="3174012304900001"
                  maxLength={16}
                  value={ktpNumber}
                  onChange={(e) => setKtpNumber(e.target.value)}
                  required
                  disabled={ktpLocked}
                />
              </div>

              <div className="form-group">
                <label>Upload Foto KTP (JPG/PNG)</label>
                <input
                  type="file"
                  accept="image/*"
                  disabled={ktpLocked}
                  onChange={(e) => {
                    if (e.target.files[0]) {
                      setKtpFile(e.target.files[0]);
                      setKtpPreviewURL(URL.createObjectURL(e.target.files[0]));
                    }
                  }}
                />
                {ktpPreviewSrc && (
                  <div className="image-preview-box">
                    <img src={ktpPreviewSrc} alt="KTP Preview" className="ktp-img-preview" />
                  </div>
                )}
              </div>

              <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading || ktpLocked}>
                {ktpLocked ? 'KTP Sudah Terverifikasi' : loading ? 'Mengunggah & Mengajukan...' : 'Ajukan Verifikasi KTP'}
              </button>
            </form>
          )}

          {activeTab === 'resume' && (
            <form onSubmit={handleUploadResume} className="auth-form">
              <p className="tab-description">
                Unggah file Resume / CV terbaru Anda dalam format PDF untuk ditunjukkan kepada calon pemberi kerja.
              </p>

              <div className="form-group">
                <label>Pilih File PDF Resume</label>
                <input
                  type="file"
                  accept=".pdf"
                  onChange={(e) => {
                    if (e.target.files[0]) {
                      setResumeFile(e.target.files[0]);
                    }
                  }}
                />
              </div>

              <div className="form-group">
                <label>Atau Input Link Resume (Google Drive / Cloudinary)</label>
                <input
                  type="url"
                  placeholder="https://res.cloudinary.com/demo/image/upload/v123/cv.pdf"
                  value={resumeURLInput}
                  onChange={(e) => setResumeURLInput(e.target.value)}
                />
              </div>

              <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading}>
                {loading ? 'Simpan Resume...' : 'Unggah Resume PDF'}
              </button>
            </form>
          )}
        </div>
      </div>

      <SubscriptionModal
        isOpen={isSubModalOpen}
        onClose={() => setIsSubModalOpen(false)}
        onSuccess={() => fetchSubscriptions()}
      />
    </>
  );
};

export default UserProfileModal;
