import React, { useState } from 'react';
import { Briefcase, Repeat, ShieldCheck, Zap, User, Sun, Moon, LogOut, FileText, Settings, Clock, PlusCircle, CheckCircle2, Star } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import AuthModal from './AuthModal';
import UserProfileModal from './UserProfileModal';
import PostSkillModal from './PostSkillModal';
import ReceivedOffersModal from './ReceivedOffersModal';
import SubscriptionModal from './SubscriptionModal';
import PostJobModal from './PostJobModal';
import ActiveConnectionsModal from './ActiveConnectionsModal';
import RatingModal from './RatingModal';

export default function Navbar({ theme, toggleTheme }) {
  const { user, token, logout, switchMode } = useAuth();
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);
  const [authModalTab, setAuthModalTab] = useState('login');
  const [isProfileModalOpen, setIsProfileModalOpen] = useState(false);
  const [isPostSkillOpen, setIsPostSkillOpen] = useState(false);
  const [isPostJobOpen, setIsPostJobOpen] = useState(false);
  const [isReceivedOffersOpen, setIsReceivedOffersOpen] = useState(false);
  const [isActiveConnectionsOpen, setIsActiveConnectionsOpen] = useState(false);
  const [isRatingModalOpen, setIsRatingModalOpen] = useState(false);
  const [ratingTarget, setRatingTarget] = useState({ connectionId: null, partnerName: '' });
  const [isSubModalOpen, setIsSubModalOpen] = useState(false);
  const [subModalPlan, setSubModalPlan] = useState('POST_SKILL_10K');
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [switching, setSwitching] = useState(false);

  const activeMode = user?.active_mode || 'SEEKER';

  const handleToggleMode = async () => {
    if (!token) {
      setAuthModalTab('login');
      setIsAuthModalOpen(true);
      return;
    }

    setSwitching(true);
    const newMode = activeMode === 'SEEKER' ? 'EMPLOYER' : 'SEEKER';
    try {
      await switchMode(newMode);
    } catch (err) {
      alert(err.message);
    } finally {
      setSwitching(false);
    }
  };

  const openAuth = (tabName) => {
    setAuthModalTab(tabName);
    setIsAuthModalOpen(true);
  };

  const openSubModal = (planType) => {
    setSubModalPlan(planType);
    setIsSubModalOpen(true);
  };

  const handleOpenRating = (connectionId, partnerName) => {
    setRatingTarget({ connectionId, partnerName });
    setIsRatingModalOpen(true);
  };

  return (
    <>
      <nav className="navbar">
        <a href="/" className="brand-logo">
          <Zap className="w-6 h-6 text-emerald-400" color="#10b981" size={26} />
          <span>Layanesia</span>
        </a>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
          {/* Dual Role Mode Switcher Widget */}
          <div 
            onClick={handleToggleMode}
            className={`mode-switcher-widget ${switching ? 'disabled' : ''}`}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '10px',
              padding: '6px 14px',
              borderRadius: '24px',
              background: activeMode === 'SEEKER' ? 'rgba(16, 185, 129, 0.12)' : 'rgba(99, 102, 241, 0.12)',
              border: activeMode === 'SEEKER' ? '1px solid rgba(16, 185, 129, 0.3)' : '1px solid rgba(99, 102, 241, 0.3)',
              cursor: 'pointer',
              transition: 'all 0.3s ease',
              userSelect: 'none',
            }}
            title="Klik untuk berganti secara langsung antara Mode Pencari Kerja & Pemberi Kerja"
          >
            <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)', fontWeight: 600 }}>Mode:</span>
            {activeMode === 'SEEKER' ? (
              <span className="badge badge-seeker" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                <User size={14} /> Pencari Kerja
              </span>
            ) : (
              <span className="badge badge-employer" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                <Briefcase size={14} /> Pemberi Kerja
              </span>
            )}
            <Repeat size={14} style={{ color: 'var(--text-muted)' }} />
          </div>

          {/* Quick Action: Active Connections (Terhubung) */}
          {user && (
            <button
              onClick={() => setIsActiveConnectionsOpen(true)}
              className="btn btn-outline"
              style={{
                padding: '6px 12px',
                borderRadius: '20px',
                fontSize: '0.8rem',
                borderColor: 'rgba(16, 185, 129, 0.4)',
                color: '#10b981',
                display: 'flex',
                alignItems: 'center',
                gap: '6px'
              }}
              title="Lihat Kontak Partner Kerja (Status TERHUBUNG)"
            >
              <CheckCircle2 size={16} /> Kontak Terhubung
            </button>
          )}

          {/* Theme Toggle Button (Light vs Dark) */}
          <button 
            onClick={toggleTheme}
            className="btn btn-outline"
            style={{
              padding: '8px 12px',
              borderRadius: '20px',
              fontSize: '0.8rem',
              borderColor: 'var(--border-color)',
              color: 'var(--text-main)',
              display: 'flex',
              alignItems: 'center',
              gap: '6px'
            }}
            title={`Beralih ke Mode ${theme === 'dark' ? 'Cerah' : 'Gelap'}`}
          >
            {theme === 'dark' ? (
              <>
                <Sun size={16} color="#f59e0b" /> <span style={{ fontSize: '0.78rem' }}>Mode Cerah</span>
              </>
            ) : (
              <>
                <Moon size={16} color="#6366f1" /> <span style={{ fontSize: '0.78rem' }}>Mode Gelap</span>
              </>
            )}
          </button>

          {/* User Auth Buttons or Profile Menu */}
          {user ? (
            <div style={{ position: 'relative' }}>
              <div 
                className="user-profile-button"
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  padding: '6px 12px',
                  borderRadius: '20px',
                  background: 'var(--bg-secondary)',
                  border: '1px solid var(--border-color)',
                  cursor: 'pointer',
                }}
              >
                <div 
                  style={{
                    width: '28px',
                    height: '28px',
                    borderRadius: '50%',
                    background: 'var(--accent-gradient)',
                    color: '#fff',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontWeight: 'bold',
                    fontSize: '0.85rem'
                  }}
                >
                  {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
                </div>
                <span style={{ fontSize: '0.88rem', fontWeight: 600, color: 'var(--text-main)' }}>
                  {user.name.split(' ')[0]}
                </span>
                {user.is_ktp_verified && (
                  <ShieldCheck size={16} color="#10b981" title="KTP Terverifikasi" />
                )}
              </div>

              {/* Dropdown Menu */}
              {isDropdownOpen && (
                <div 
                  className="profile-dropdown-menu"
                  onClick={(e) => e.stopPropagation()}
                  style={{
                    position: 'absolute',
                    right: 0,
                    top: '46px',
                    width: '220px',
                    background: 'var(--bg-secondary)',
                    border: '1px solid var(--border-color)',
                    borderRadius: '12px',
                    boxShadow: '0 10px 25px rgba(0,0,0,0.2)',
                    zIndex: 100,
                    padding: '8px 0',
                  }}
                >
                  <div style={{ padding: '8px 16px', borderBottom: '1px solid var(--border-color)' }}>
                    <div style={{ fontWeight: 600, fontSize: '0.9rem', color: 'var(--text-main)' }}>{user.name}</div>
                    <div style={{ fontSize: '0.78rem', color: 'var(--text-muted)' }}>{user.email}</div>
                  </div>

                  <button
                    onClick={() => { setIsPostSkillOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item text-emerald"
                    style={{ color: '#10b981', fontWeight: 600 }}
                  >
                    <PlusCircle size={16} /> Post Keahlian Saya
                  </button>

                  <button
                    onClick={() => { setIsPostJobOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item"
                    style={{ color: '#6366f1', fontWeight: 600 }}
                  >
                    <Briefcase size={16} /> Buka Lowongan Kerja (Gratis)
                  </button>

                  <button
                    onClick={() => { setIsActiveConnectionsOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item"
                  >
                    <CheckCircle2 size={16} /> Kontak Matched (Terhubung)
                  </button>

                  <button
                    onClick={() => { setIsReceivedOffersOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item"
                  >
                    <Clock size={16} /> Tawaran Pekerjaan Masuk
                  </button>

                  <button
                    onClick={() => { setIsProfileModalOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item"
                  >
                    <Settings size={16} /> Pengaturan Profil
                  </button>

                  <button
                    onClick={() => { setIsProfileModalOpen(true); setIsDropdownOpen(false); }}
                    className="dropdown-item"
                  >
                    <ShieldCheck size={16} /> Verifikasi KTP
                  </button>

                  <button
                    onClick={() => { logout(); setIsDropdownOpen(false); }}
                    className="dropdown-item text-danger"
                  >
                    <LogOut size={16} /> Keluar Akun
                  </button>
                </div>
              )}
            </div>
          ) : (
            <div style={{ display: 'flex', gap: '8px' }}>
              <button 
                onClick={() => openAuth('login')}
                className="btn btn-outline" 
                style={{ padding: '8px 14px', fontSize: '0.85rem' }}
              >
                Masuk
              </button>
              <button 
                onClick={() => openAuth('register')}
                className="btn btn-primary" 
                style={{ padding: '8px 14px', fontSize: '0.85rem' }}
              >
                Daftar
              </button>
            </div>
          )}
        </div>
      </nav>

      {/* Auth Modal */}
      <AuthModal
        isOpen={isAuthModalOpen}
        onClose={() => setIsAuthModalOpen(false)}
        initialTab={authModalTab}
      />

      {/* User Profile & KTP Modal */}
      <UserProfileModal
        isOpen={isProfileModalOpen}
        onClose={() => setIsProfileModalOpen(false)}
      />

      {/* Post Skill Modal (Sprint 4) */}
      <PostSkillModal
        isOpen={isPostSkillOpen}
        onClose={() => setIsPostSkillOpen(false)}
        onOpenSubModal={openSubModal}
      />

      {/* Post Job Modal (Sprint 5) */}
      <PostJobModal
        isOpen={isPostJobOpen}
        onClose={() => setIsPostJobOpen(false)}
      />

      {/* Received Offers Modal (Sprint 4) */}
      <ReceivedOffersModal
        isOpen={isReceivedOffersOpen}
        onClose={() => setIsReceivedOffersOpen(false)}
      />

      {/* Active Connections Modal (Sprint 6) */}
      <ActiveConnectionsModal
        isOpen={isActiveConnectionsOpen}
        onClose={() => setIsActiveConnectionsOpen(false)}
        onOpenRatingModal={handleOpenRating}
      />

      {/* Reciprocal Rating Modal (Sprint 6) */}
      <RatingModal
        isOpen={isRatingModalOpen}
        onClose={() => setIsRatingModalOpen(false)}
        connectionId={ratingTarget.connectionId}
        partnerName={ratingTarget.partnerName}
      />

      {/* Midtrans Subscription Modal */}
      <SubscriptionModal
        isOpen={isSubModalOpen}
        onClose={() => setIsSubModalOpen(false)}
        selectedPlan={subModalPlan}
      />
    </>
  );
}
