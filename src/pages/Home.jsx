import React, { useState, useEffect } from 'react';
import { 
  Search, MapPin, ShieldCheck, Clock, 
  Briefcase, User, ChevronRight, PlusCircle, Users
} from 'lucide-react';
import LocationSelector from '../components/LocationSelector';
import LocationPickerModal from '../components/LocationPickerModal';
import SubscriptionModal from '../components/SubscriptionModal';
import UserProfileModal from '../components/UserProfileModal';
import PostSkillModal from '../components/PostSkillModal';
import DirectOfferModal from '../components/DirectOfferModal';
import PostJobModal from '../components/PostJobModal';
import ApplyJobModal from '../components/ApplyJobModal';
import JobApplicantsModal from '../components/JobApplicantsModal';
import LandingSections from '../components/LandingSections';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export default function Home({ activeMode: _activeMode }) {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState('MODEL_1'); // 'MODEL_1' (Talent Showcase) or 'MODEL_2' (Job Board)
  
  // Search & Location Filter State
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');
  const [locationState, setLocationState] = useState({
    province: '',
    city: '',
    district: ''
  });

  // Location Map Modal State
  const [isMapModalOpen, setIsMapModalOpen] = useState(false);
  const [pinnedLocation, setPinnedLocation] = useState(null);

  // Subscription Modal State
  const [isSubModalOpen, setIsSubModalOpen] = useState(false);
  const [subModalPlan, setSubModalPlan] = useState('APPLY_JOB_5K');
  const [isProfileModalOpen, setIsProfileModalOpen] = useState(false);

  // Sprint 4 Modals State
  const [isPostSkillOpen, setIsPostSkillOpen] = useState(false);
  const [isOfferModalOpen, setIsOfferModalOpen] = useState(false);
  const [selectedTalent, setSelectedTalent] = useState(null);

  // Sprint 5 Modals State
  const [isPostJobOpen, setIsPostJobOpen] = useState(false);
  const [isApplyJobOpen, setIsApplyJobOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState(null);
  const [isApplicantsModalOpen, setIsApplicantsModalOpen] = useState(false);
  const [selectedJobForApplicants, setSelectedJobForApplicants] = useState(null);

  // Real Data States
  const [realTalents, setRealTalents] = useState([]);
  const [loadingTalents, setLoadingTalents] = useState(false);

  const [realJobs, setRealJobs] = useState([]);
  const [loadingJobs, setLoadingJobs] = useState(false);
  const [listError, setListError] = useState('');

  const fetchTalents = async () => {
    setLoadingTalents(true);
    setListError('');
    try {
      const params = new URLSearchParams();
      if (locationState.province) params.append('province', locationState.province);
      if (locationState.city) params.append('city', locationState.city);
      if (locationState.district) params.append('district', locationState.district);
      if (categoryFilter) params.append('category', categoryFilter);
      if (searchQuery) params.append('search', searchQuery);

      const qs = params.toString();
      const data = await api(`/api/v1/skills${qs ? `?${qs}` : ''}`);
      setRealTalents(data.data || []);
    } catch (err) {
      console.error('Error fetching talents:', err);
      setRealTalents([]);
      setListError(err.message || 'Gagal memuat etalase keahlian.');
    } finally {
      setLoadingTalents(false);
    }
  };

  const fetchJobs = async () => {
    setLoadingJobs(true);
    setListError('');
    try {
      const params = new URLSearchParams();
      if (locationState.province) params.append('province', locationState.province);
      if (locationState.city) params.append('city', locationState.city);
      if (locationState.district) params.append('district', locationState.district);
      if (categoryFilter) params.append('category', categoryFilter);
      if (searchQuery) params.append('search', searchQuery);

      const qs = params.toString();
      const data = await api(`/api/v1/jobs${qs ? `?${qs}` : ''}`);
      setRealJobs(data.data || []);
    } catch (err) {
      console.error('Error fetching jobs:', err);
      setRealJobs([]);
      setListError(err.message || 'Gagal memuat lowongan.');
    } finally {
      setLoadingJobs(false);
    }
  };

  useEffect(() => {
    if (activeTab === 'MODEL_1') {
      fetchTalents();
    } else {
      fetchJobs();
    }
  }, [activeTab, categoryFilter, locationState.province, locationState.city, locationState.district]);

  const handleLocationChange = (loc) => {
    setLocationState(loc);
  };

  const handleConfirmPinnedLocation = (pinData) => {
    setPinnedLocation(pinData);
    if (pinData) {
      setLocationState(prev => ({
        ...prev,
        province: pinData.province || prev.province,
        city: pinData.city || prev.city,
        district: pinData.district || prev.district
      }));
    }
  };

  const openSubModal = (planType) => {
    setSubModalPlan(planType);
    setIsSubModalOpen(true);
  };

  const handleOpenOfferModal = (talent) => {
    setSelectedTalent(talent);
    setIsOfferModalOpen(true);
  };

  const handleOpenApplyModal = (job) => {
    setSelectedJob(job);
    setIsApplyJobOpen(true);
  };

  const handleOpenApplicantsModal = (job) => {
    setSelectedJobForApplicants(job);
    setIsApplicantsModalOpen(true);
  };

  return (
    <div className="container" style={{ paddingTop: '20px' }}>
      {/* Informative Landing Page Experience: Hero, Advantages, Workflow, Calculator, FAQ, CTA */}
      <LandingSections
        onOpenPostSkill={() => setIsPostSkillOpen(true)}
        onOpenPostJob={() => setIsPostJobOpen(true)}
        onOpenMap={() => setIsMapModalOpen(true)}
        onScrollToCatalog={() => {
          document.getElementById('katalog')?.scrollIntoView({ behavior: 'smooth' });
        }}
        activeMode={user?.active_mode || 'SEEKER'}
      />

      {/* Hiperlokal Location Filter & Live Marketplace Section */}
      <div id="katalog" style={{ scrollMarginTop: '80px', marginBottom: '32px' }}>
        <div style={{ marginBottom: '18px' }}>
          <div style={{ fontSize: '0.85rem', fontWeight: 700, color: '#10b981', textTransform: 'uppercase', letterSpacing: '0.06em', marginBottom: '6px' }}>
            Bursa Langsung Hiperlokal
          </div>
          <h2 style={{ fontSize: '1.9rem', color: 'var(--text-main)', letterSpacing: '-0.02em', margin: 0 }}>
            Eksplorasi Etalase Jasa & Lowongan Pekerjaan
          </h2>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.92rem', margin: '4px 0 0 0' }}>
            Filter tenaga kerja terampil atau lowongan insidental berdasarkan kecamatan domisili Anda.
          </p>
        </div>

        {/* Hiperlokal Location Filter Bar */}
        <div className="glass-panel" style={{ padding: '24px', marginBottom: '24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px', flexWrap: 'wrap', gap: '10px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#10b981', fontSize: '0.95rem', fontWeight: 700 }}>
            <MapPin size={20} /> Filter Hiperlokal Presisi (Backend Master Data Locations)
          </div>

          <button 
            onClick={() => setIsMapModalOpen(true)}
            className="btn btn-outline"
            style={{ padding: '6px 14px', fontSize: '0.8rem', borderColor: 'rgba(16,185,129,0.4)', color: '#10b981' }}
          >
            📍 {pinnedLocation ? 'Titik Peta Dipilih' : 'Pilih Titik di Peta (Leaflet.js)'}
          </button>
        </div>

        {pinnedLocation && (
          <div style={{ background: 'rgba(16,185,129,0.1)', border: '1px solid rgba(16,185,129,0.3)', borderRadius: '8px', padding: '8px 14px', marginBottom: '16px', fontSize: '0.85rem', color: '#10b981' }}>
            <strong>Titik Presisi Peta:</strong> {pinnedLocation.formattedAddress} (Lat: {pinnedLocation.latitude.toFixed(4)}, Lng: {pinnedLocation.longitude.toFixed(4)})
          </div>
        )}

        <div style={{ display: 'flex', gap: '14px', alignItems: 'flex-end', flexWrap: 'wrap' }}>
          <div style={{ flex: 1, minWidth: '240px' }}>
            <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-muted)', marginBottom: '4px' }}>Kata Kunci Pekerjaan / Keahlian</label>
            <div style={{ position: 'relative' }}>
              <input
                type="text"
                className="input-field"
                placeholder="Supir, ART, Admin Logistik, Tukang..."
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                style={{ paddingLeft: '36px' }}
              />
              <Search size={16} style={{ position: 'absolute', left: '12px', top: '13px', color: 'var(--text-muted)' }} />
            </div>
          </div>

          <div style={{ flex: 2, minWidth: '300px' }}>
            <LocationSelector onLocationChange={handleLocationChange} />
          </div>

          <button className="btn btn-primary" onClick={activeTab === 'MODEL_1' ? fetchTalents : fetchJobs} style={{ height: '42px', minWidth: '140px' }}>
            <Search size={18} /> Cari Data
          </button>
        </div>
      </div>

      {/* Model Selection Tabs */}
      <div style={{ display: 'flex', gap: '16px', marginBottom: '24px', borderBottom: '1px solid var(--border-color)', paddingBottom: '12px' }}>
        <button 
          onClick={() => setActiveTab('MODEL_1')}
          style={{
            padding: '12px 24px',
            borderRadius: '12px',
            background: activeTab === 'MODEL_1' ? 'var(--gradient-primary)' : 'var(--bg-glass)',
            color: activeTab === 'MODEL_1' ? '#ffffff' : 'var(--text-main)',
            border: activeTab === 'MODEL_1' ? 'none' : '1px solid var(--border-color)',
            fontWeight: 700,
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transition: 'all 0.2s ease'
          }}
        >
          <User size={18} /> Model 1: Talent Showcase (Cari Pekerja)
        </button>

        <button 
          onClick={() => setActiveTab('MODEL_2')}
          style={{
            padding: '12px 24px',
            borderRadius: '12px',
            background: activeTab === 'MODEL_2' ? 'var(--gradient-secondary)' : 'var(--bg-glass)',
            color: activeTab === 'MODEL_2' ? '#ffffff' : 'var(--text-main)',
            border: activeTab === 'MODEL_2' ? 'none' : '1px solid var(--border-color)',
            fontWeight: 700,
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transition: 'all 0.2s ease'
          }}
        >
          <Briefcase size={18} /> Model 2: Job Board (Cari Lowongan - Free)
        </button>
      </div>

      {/* Model 1: Talent Showcase List */}
      {activeTab === 'MODEL_1' && (
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', flexWrap: 'wrap', gap: '12px' }}>
            <div>
              <h2 style={{ fontSize: '1.4rem', color: 'var(--text-main)', margin: 0 }}>Etalase Keahlian Pekerja (Worker Showcase)</h2>
              <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', margin: 0 }}>Pilih keahlian pekerja dan kirim tawaran pekerjaan secara langsung</p>
            </div>

            <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
              <button
                className={`btn ${categoryFilter === '' ? 'btn-primary' : 'btn-outline'}`}
                style={{ padding: '6px 12px', fontSize: '0.78rem' }}
                onClick={() => setCategoryFilter('')}
              >
                Semua Kategori
              </button>
              <button
                className={`btn ${categoryFilter === 'SERABUTAN' ? 'btn-primary' : 'btn-outline'}`}
                style={{ padding: '6px 12px', fontSize: '0.78rem' }}
                onClick={() => setCategoryFilter('SERABUTAN')}
              >
                Serabutan
              </button>
              <button
                className={`btn ${categoryFilter === 'PROFESIONAL' ? 'btn-primary' : 'btn-outline'}`}
                style={{ padding: '6px 12px', fontSize: '0.78rem' }}
                onClick={() => setCategoryFilter('PROFESIONAL')}
              >
                Profesional
              </button>

              <button
                className="btn btn-primary"
                onClick={() => setIsPostSkillOpen(true)}
                style={{ padding: '6px 14px', fontSize: '0.82rem', marginLeft: '8px' }}
              >
                <PlusCircle size={16} /> Post Keahlian
              </button>
            </div>
          </div>

          {listError && (
            <div style={{ textAlign: 'center', padding: '16px', color: '#ef4444', marginBottom: '12px' }}>{listError}</div>
          )}
          {loadingTalents ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Memuat etalase keahlian pekerja...</div>
          ) : realTalents.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Belum ada etalase keahlian pada filter ini.</div>
          ) : (
            <div className="grid-3">
              {realTalents.map(talent => {
                const name = talent.user?.name || talent.name || 'Pekerja Layanesia';
                const isKtp = talent.user?.is_ktp_verified || talent.isKtpVerified;
                const locStr = `${talent.district || ''}${talent.city ? ', ' + talent.city : ''}`;
                const rateAmount = typeof talent.rate_amount === 'number' ? talent.rate_amount.toLocaleString('id-ID') : (talent.rateAmount || talent.rate_amount);
                const rateTypeLabel = (talent.rate_type || talent.rateType) === 'PER_HOUR' ? 'jam' : 'hari';

                return (
                  <div key={talent.id} className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
                    <div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                        <span className={`badge ${talent.category === 'PROFESIONAL' ? 'badge-employer' : 'badge-seeker'}`}>{talent.category}</span>
                        {isKtp && (
                          <span className="badge badge-verified" title="KTP Verified">
                            <ShieldCheck size={12} /> Verified
                          </span>
                        )}
                      </div>

                      <h3 style={{ fontSize: '1.15rem', marginBottom: '8px', color: 'var(--text-main)' }}>{talent.title}</h3>
                      <p style={{ fontSize: '0.85rem', color: '#10b981', fontWeight: 600, marginBottom: '6px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                        <User size={14} /> {name}
                      </p>
                      <p style={{ fontSize: '0.88rem', color: 'var(--text-muted)', marginBottom: '16px', lineHeight: '1.5' }}>
                        {talent.description}
                      </p>
                    </div>

                    <div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '12px' }}>
                        <MapPin size={14} color="#10b981" /> {locStr || 'Lokasi Terdaftar'}
                      </div>

                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '14px', borderTop: '1px solid var(--border-color)' }}>
                        <div>
                          <span style={{ fontSize: '1.2rem', fontWeight: 800, color: '#10b981' }}>Rp {rateAmount}</span>
                          <span style={{ fontSize: '0.75rem', color: 'var(--text-dim)' }}> / {rateTypeLabel}</span>
                        </div>

                        <button 
                          className="btn btn-primary" 
                          style={{ padding: '6px 14px', fontSize: '0.85rem' }}
                          onClick={() => handleOpenOfferModal(talent)}
                        >
                          Kirim Offer <ChevronRight size={14} />
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* Model 2: Demand Job Board List */}
      {activeTab === 'MODEL_2' && (
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', flexWrap: 'wrap', gap: '12px' }}>
            <div>
              <h2 style={{ fontSize: '1.4rem', color: 'var(--text-main)', margin: 0 }}>Lowongan Kerja Insidental (Job Board - Free Post)</h2>
              <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', margin: 0 }}>Pemberi Kerja membuka lowongan gratis, Pekerja melamar (Rp 5k/bln)</p>
            </div>

            <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
              <button
                className="btn btn-primary"
                onClick={() => setIsPostJobOpen(true)}
                style={{ padding: '6px 14px', fontSize: '0.82rem', background: 'var(--gradient-secondary)' }}
              >
                <PlusCircle size={16} /> Post Lowongan (Free)
              </button>
            </div>
          </div>

          {listError && (
            <div style={{ textAlign: 'center', padding: '16px', color: '#ef4444', marginBottom: '12px' }}>{listError}</div>
          )}
          {loadingJobs ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Memuat lowongan pekerjaan...</div>
          ) : realJobs.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>Belum ada lowongan pada filter ini.</div>
          ) : (
            <div className="grid-2">
              {realJobs.map(job => {
                const employerName = job.employer?.name || job.employerName || 'Pemberi Kerja';
                const isKtp = job.employer?.is_ktp_verified;
                const budgetStr = typeof job.budget === 'number' ? `Rp ${job.budget.toLocaleString('id-ID')}` : (job.budget || '-');
                const workDateStr = job.workDate || job.work_date?.split('T')[0] || '-';
                const durationStr = `${job.duration_value || 1} ${job.duration_type === 'DAYS' ? 'Hari' : 'Jam'}`;
                const locationStr = `${job.district || ''}${job.city ? ', ' + job.city : ''}`;
                const isMyJob = user && job.employer_id === user.id;

                return (
                  <div key={job.id} className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
                    <div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                        <span className="badge badge-employer">{job.category}</span>
                        <span className="badge badge-free">100% GRATIS POST</span>
                      </div>

                      <h3 style={{ fontSize: '1.2rem', marginBottom: '8px', color: 'var(--text-main)' }}>{job.title}</h3>
                      <p style={{ fontSize: '0.85rem', color: '#6366f1', fontWeight: 600, marginBottom: '6px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                        <Briefcase size={14} /> {employerName}
                        {isKtp && <ShieldCheck size={14} color="#10b981" title="KTP Verified" />}
                      </p>
                      <p style={{ fontSize: '0.88rem', color: 'var(--text-muted)', marginBottom: '16px', lineHeight: '1.5' }}>
                        {job.description}
                      </p>
                    </div>

                    <div>
                      <div style={{ display: 'flex', gap: '16px', fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '14px', flexWrap: 'wrap' }}>
                        <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                          <MapPin size={14} color="#6366f1" /> {locationStr || 'Lokasi Terdaftar'}
                        </span>
                        <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                          <Clock size={14} color="#f59e0b" /> Pelaksanaan: {workDateStr} ({durationStr})
                        </span>
                      </div>

                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '14px', borderTop: '1px solid var(--border-color)' }}>
                        <div>
                          <span style={{ fontSize: '0.75rem', color: 'var(--text-dim)', display: 'block' }}>Total Honor:</span>
                          <span style={{ fontSize: '1.25rem', fontWeight: 800, color: '#6366f1' }}>{budgetStr}</span>
                        </div>

                        {isMyJob ? (
                          <button
                            onClick={() => handleOpenApplicantsModal(job)}
                            className="btn btn-outline"
                            style={{ padding: '8px 16px', fontSize: '0.85rem', borderColor: '#6366f1', color: '#6366f1' }}
                          >
                            <Users size={16} /> Reviu Pelamar
                          </button>
                        ) : (
                          <button 
                            onClick={() => handleOpenApplyModal(job)}
                            className="btn btn-secondary" 
                            style={{ padding: '8px 18px', fontSize: '0.85rem' }}
                          >
                            Apply Job (Rp 5k/bln)
                          </button>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
      </div>

      {/* Location Picker Map Modal */}
      <LocationPickerModal 
        isOpen={isMapModalOpen}
        onClose={() => setIsMapModalOpen(false)}
        onConfirmLocation={handleConfirmPinnedLocation}
      />

      {/* User Profile Modal */}
      <UserProfileModal
        isOpen={isProfileModalOpen}
        onClose={() => setIsProfileModalOpen(false)}
      />

      {/* Midtrans Subscription Modal */}
      <SubscriptionModal
        isOpen={isSubModalOpen}
        onClose={() => setIsSubModalOpen(false)}
        selectedPlan={subModalPlan}
      />

      {/* Post Skill Modal (Sprint 4) */}
      <PostSkillModal
        isOpen={isPostSkillOpen}
        onClose={() => setIsPostSkillOpen(false)}
        onSuccess={fetchTalents}
        onOpenSubModal={openSubModal}
      />

      {/* Direct Offer Modal (Sprint 4) */}
      <DirectOfferModal
        isOpen={isOfferModalOpen}
        onClose={() => setIsOfferModalOpen(false)}
        talent={selectedTalent}
      />

      {/* Post Job Modal (Sprint 5) */}
      <PostJobModal
        isOpen={isPostJobOpen}
        onClose={() => setIsPostJobOpen(false)}
        onSuccess={fetchJobs}
      />

      {/* Apply Job Modal (Sprint 5) */}
      <ApplyJobModal
        isOpen={isApplyJobOpen}
        onClose={() => setIsApplyJobOpen(false)}
        job={selectedJob}
        onOpenProfileModal={() => setIsProfileModalOpen(true)}
        onOpenSubModal={openSubModal}
      />

      {/* Job Applicants Modal (Sprint 5) */}
      <JobApplicantsModal
        isOpen={isApplicantsModalOpen}
        onClose={() => setIsApplicantsModalOpen(false)}
        job={selectedJobForApplicants}
      />
    </div>
  );
}
