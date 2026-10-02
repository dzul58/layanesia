import React, { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { api, ApiError, loadMidtransSnap } from '../api';
import { CheckCircle2, DollarSign } from 'lucide-react';

const PLAN_LABEL = {
  APPLY_JOB_5K: 'Melamar Lowongan 5K',
  POST_SKILL_10K: 'Etalase Keahlian 10K',
};

const SubscriptionModal = ({ isOpen, onClose, selectedPlan = 'APPLY_JOB_5K', onSuccess }) => {
  const { token, refreshUser } = useAuth();
  const [plan, setPlan] = useState(selectedPlan);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    if (isOpen) {
      setPlan(selectedPlan || 'APPLY_JOB_5K');
      setError('');
      setSuccess('');
    }
  }, [isOpen, selectedPlan]);

  if (!isOpen) return null;

  const syncSubscription = async (subscriptionId) => {
    if (!subscriptionId) return null;
    try {
      const data = await api(`/api/v1/subscriptions/${subscriptionId}/sync`, {
        method: 'POST',
        token,
      });
      return data.subscription || null;
    } catch {
      return null;
    }
  };

  const finishIfActive = async (subscription, checkoutData) => {
    const status = subscription?.status || checkoutData?.status;
    if (status === 'ACTIVE') {
      setSuccess(`Pembayaran berhasil. Paket ${PLAN_LABEL[plan] || plan} aktif selama 30 hari.`);
      if (refreshUser) await refreshUser();
      if (onSuccess) onSuccess(subscription || checkoutData);
      setTimeout(() => onClose(), 1500);
      return true;
    }
    return false;
  };

  const payWithSnap = async (checkout) => {
    const snap = await loadMidtransSnap();
    const tokenLooksFake = !checkout.snap_token || String(checkout.snap_token).startsWith('SNAP-SIM');

    if (snap && !tokenLooksFake) {
      await new Promise((resolve) => {
        snap.pay(checkout.snap_token, {
          onSuccess: async () => {
            const sub = await syncSubscription(checkout.subscription_id);
            await finishIfActive(sub, checkout);
            resolve();
          },
          onPending: async () => {
            await syncSubscription(checkout.subscription_id);
            setSuccess('Pembayaran masih diproses. Paket aktif setelah Midtrans mengonfirmasi.');
            if (onSuccess) onSuccess(checkout);
            resolve();
          },
          onError: () => {
            setError('Pembayaran gagal atau dibatalkan. Silakan coba lagi.');
            resolve();
          },
          onClose: async () => {
            const sub = await syncSubscription(checkout.subscription_id);
            if (!(await finishIfActive(sub, checkout))) {
              setSuccess('Jendela pembayaran ditutup. Jika sudah bayar, status akan menyusul.');
            }
            resolve();
          },
        });
      });
      return;
    }

    if (checkout.redirect_url && !String(checkout.redirect_url).includes('SNAP-SIM')) {
      window.location.href = checkout.redirect_url;
      return;
    }

    throw new Error('Gateway pembayaran belum dikonfigurasi. Isi kunci Midtrans di server dan VITE_MIDTRANS_CLIENT_KEY.');
  };

  const tryDevSimulation = async () => {
    const data = await api('/api/v1/subscriptions/simulate-activate', {
      method: 'POST',
      token,
      body: { plan_type: plan },
    });
    setSuccess(data.message || 'Mode pengembangan: paket diaktifkan tanpa pembayaran.');
    if (refreshUser) await refreshUser();
    if (onSuccess) onSuccess(data.subscription);
    setTimeout(() => onClose(), 1500);
  };

  const handleCheckout = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      if (!token) {
        throw new Error('Silakan masuk / mendaftar akun terlebih dahulu.');
      }

      const res = await api('/api/v1/subscriptions/checkout', {
        method: 'POST',
        token,
        body: { plan_type: plan },
      });
      const checkout = res.data || {};
      await payWithSnap(checkout);
    } catch (err) {
      const canSimulate = import.meta.env.DEV && err instanceof ApiError && err.code === 'SERVICE_UNAVAILABLE';
      if (canSimulate) {
        try {
          await tryDevSimulation();
        } catch (simErr) {
          setError(simErr.message || err.message);
        }
      } else {
        setError(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content subscription-modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close-btn" onClick={onClose}>&times;</button>

        <div className="auth-header">
          <div className="auth-logo">
            <DollarSign color="#10b981" size={28} />
            <h2>Midtrans Subscription</h2>
          </div>
          <p className="auth-subtitle">Pilih Paket Langganan Terjangkau Layanesia (Aktif 30 Hari)</p>
        </div>

        {error && <div className="auth-error-banner">⚠️ {error}</div>}
        {success && <div className="auth-success-banner">✅ {success}</div>}

        <div className="subscription-plan-grid">
          <div
            className={`plan-card ${plan === 'APPLY_JOB_5K' ? 'selected' : ''}`}
            onClick={() => setPlan('APPLY_JOB_5K')}
          >
            <div className="plan-badge">POPULER PEKERJA</div>
            <h3 className="plan-title">Paket Melamar Lowongan</h3>
            <div className="plan-price">
              <span className="amount">Rp 5.000</span>
              <span className="period">/ 30 Hari</span>
            </div>
            <ul className="plan-features">
              <li><CheckCircle2 size={16} color="#10b981" /> Apply unlimited lowongan harian di Model 2</li>
              <li><CheckCircle2 size={16} color="#10b981" /> Dapatkan notifikasi saat pemberi kerja menerima lamaran</li>
              <li><CheckCircle2 size={16} color="#10b981" /> Lencana pelamar terpercaya</li>
            </ul>
          </div>

          <div
            className={`plan-card ${plan === 'POST_SKILL_10K' ? 'selected' : ''}`}
            onClick={() => setPlan('POST_SKILL_10K')}
          >
            <div className="plan-badge employer">ETALASE TALENT</div>
            <h3 className="plan-title">Paket Pasang Keahlian</h3>
            <div className="plan-price">
              <span className="amount">Rp 10.000</span>
              <span className="period">/ 30 Hari</span>
            </div>
            <ul className="plan-features">
              <li><CheckCircle2 size={16} color="#6366f1" /> Pasang etalase jasa/keahlian di Model 1</li>
              <li><CheckCircle2 size={16} color="#6366f1" /> Terima penawaran kerja (Direct Offer) dari pemberi kerja</li>
              <li><CheckCircle2 size={16} color="#6366f1" /> Tampil teratas pada pencarian hiperlokal</li>
            </ul>
          </div>
        </div>

        <form onSubmit={handleCheckout} style={{ marginTop: '20px' }}>
          <button type="submit" className="btn btn-primary auth-submit-btn" disabled={loading}>
            {loading ? 'Memproses Midtrans Snap Checkout...' : `Bayar via Midtrans Snap (${plan === 'APPLY_JOB_5K' ? 'Rp 5.000' : 'Rp 10.000'})`}
          </button>
        </form>
      </div>
    </div>
  );
};

export default SubscriptionModal;
