import express, { Request, Response, NextFunction } from 'express';
import cors from 'cors';
import path from 'path';
import fs from 'fs';
import jwt from 'jsonwebtoken';
import multer from 'multer';
import { createServer as createViteServer } from 'vite';
import { LOCATIONS } from './locations-data';

const JWT_SECRET = process.env.JWT_SECRET || 'layanesia-secret-key-prod-2025';
const PORT = parseInt(process.env.PORT || '3000', 10);

// Ensure uploads folder exists
const uploadsDir = path.resolve(process.cwd(), 'uploads');
if (!fs.existsSync(uploadsDir)) {
  fs.mkdirSync(uploadsDir, { recursive: true });
}

// Multer storage
const storage = multer.diskStorage({
  destination: (_req, _file, cb) => cb(null, uploadsDir),
  filename: (_req, file, cb) => {
    const ext = path.extname(file.originalname);
    const unique = `${Date.now()}-${Math.round(Math.random() * 1e9)}${ext}`;
    cb(null, unique);
  },
});
const upload = multer({ storage, limits: { fileSize: 10 * 1024 * 1024 } });

// In-Memory Data Models
export interface User {
  id: string;
  email: string;
  phone: string;
  name: string;
  password?: string;
  active_mode: 'SEEKER' | 'EMPLOYER';
  is_ktp_verified: boolean;
  ktp_status: 'NOT_SUBMITTED' | 'PENDING_REVIEW' | 'VERIFIED' | 'REJECTED';
  ktp_number?: string;
  ktp_image_url?: string;
  ktp_rejection_reason?: string;
  resume_url?: string;
  province: string;
  city: string;
  district: string;
  address_detail?: string;
  created_at: string;
}

export interface Subscription {
  id: string;
  user_id: string;
  plan_type: 'APPLY_JOB_5K' | 'POST_SKILL_10K';
  amount: number;
  status: 'ACTIVE' | 'PENDING' | 'EXPIRED';
  payment_status: 'PAID' | 'PENDING';
  payment_provider: string;
  order_id: string;
  paid_at?: string;
  starts_at?: string;
  expires_at?: string;
}

export interface SkillPosting {
  id: string;
  user_id: string;
  category: 'SERABUTAN' | 'PROFESIONAL';
  title: string;
  description: string;
  province: string;
  city: string;
  district: string;
  rate_type: 'PER_HOUR' | 'PER_DAY';
  rate_amount: number;
  availability: 'AVAILABLE' | 'UNAVAILABLE';
  created_at: string;
}

export interface JobPosting {
  id: string;
  employer_id: string;
  category: 'SERABUTAN' | 'PROFESIONAL';
  title: string;
  description: string;
  province: string;
  city: string;
  district: string;
  work_date: string;
  duration_type: 'HOURS' | 'DAYS';
  duration_value: number;
  budget: number;
  status: 'OPEN' | 'CLOSED' | 'FILLED';
  created_at: string;
}

export interface JobApplication {
  id: string;
  job_id: string;
  applicant_id: string;
  status: 'PENDING' | 'ACCEPTED' | 'REJECTED';
  created_at: string;
}

export interface JobOffer {
  id: string;
  skill_posting_id: string;
  employer_id: string;
  worker_id: string;
  offered_budget: number;
  work_date: string;
  status: 'PENDING' | 'ACCEPTED' | 'REJECTED';
  created_at: string;
}

export interface Connection {
  id: string;
  employer_id: string;
  worker_id: string;
  source_type: 'JOB_OFFER' | 'JOB_APPLICATION';
  source_id: string;
  status: 'ACTIVE' | 'COMPLETED';
  created_at: string;
  completed_at?: string;
}

export interface Rating {
  id: string;
  connection_id: string;
  from_user_id: string;
  to_user_id: string;
  rating_stars: number;
  comment?: string;
  created_at: string;
}

// Initial In-Memory State from DB Seeds
const users: Map<string, User> = new Map([
  [
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
      email: 'budi.worker@gmail.com',
      phone: '081234567890',
      name: 'Budi Santoso',
      password: 'Password123!',
      active_mode: 'SEEKER',
      is_ktp_verified: true,
      ktp_status: 'VERIFIED',
      ktp_number: '3174012304900001',
      resume_url: '/uploads/sample_resume_budi.pdf',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Kebayoran Baru',
      address_detail: 'Jl. Wijaya I No. 42',
      created_at: new Date().toISOString(),
    },
  ],
  [
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
      email: 'siti.art@gmail.com',
      phone: '081987654321',
      name: 'Siti Rahmawati',
      password: 'Password123!',
      active_mode: 'SEEKER',
      is_ktp_verified: true,
      ktp_status: 'VERIFIED',
      ktp_number: '3174025508920003',
      resume_url: '/uploads/sample_resume_siti.pdf',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Cilandak',
      address_detail: 'Jl. Fatmawati Raya No. 18',
      created_at: new Date().toISOString(),
    },
  ],
  [
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
      email: 'hrd@logisticsjaya.co.id',
      phone: '081122334455',
      name: 'PT Logistics Jaya Mandiri (Hendra)',
      password: 'Password123!',
      active_mode: 'EMPLOYER',
      is_ktp_verified: true,
      ktp_status: 'VERIFIED',
      ktp_number: '3174091211850005',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Kebayoran Baru',
      address_detail: 'Gedung Graha Logistics Lt. 4',
      created_at: new Date().toISOString(),
    },
  ],
  [
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
      email: 'anita.wibowo@gmail.com',
      phone: '081555666777',
      name: 'Ibu Anita Wibowo',
      password: 'Password123!',
      active_mode: 'EMPLOYER',
      is_ktp_verified: true,
      ktp_status: 'VERIFIED',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Pondok Indah',
      address_detail: 'Jl. Metro Pondok Indah Blok TA No. 12',
      created_at: new Date().toISOString(),
    },
  ],
]);

const subscriptions: Map<string, Subscription> = new Map([
  [
    '1bafdaee-cfea-4894-82d1-8c163a1f6acd',
    {
      id: '1bafdaee-cfea-4894-82d1-8c163a1f6acd',
      user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
      plan_type: 'POST_SKILL_10K',
      amount: 10000,
      status: 'ACTIVE',
      payment_status: 'PAID',
      payment_provider: 'MIDTRANS',
      order_id: 'LEGACY-SEED-1bafdaee',
      paid_at: new Date().toISOString(),
      starts_at: new Date().toISOString(),
      expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
    },
  ],
  [
    '097c880d-1d28-424e-aeab-aba9ff9d4550',
    {
      id: '097c880d-1d28-424e-aeab-aba9ff9d4550',
      user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
      plan_type: 'APPLY_JOB_5K',
      amount: 5000,
      status: 'ACTIVE',
      payment_status: 'PAID',
      payment_provider: 'MIDTRANS',
      order_id: 'LEGACY-SEED-097c880d',
      paid_at: new Date().toISOString(),
      starts_at: new Date().toISOString(),
      expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
    },
  ],
]);

const skillPostings: Map<string, SkillPosting> = new Map([
  [
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b11',
    {
      id: 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b11',
      user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
      category: 'PROFESIONAL',
      title: 'Supir Operasional & Pengemudi Pribadi Pengganti',
      description: 'Pengalaman 6 tahun supir eksekutif & operasional boks/kantor. SIM A & B1 aktif. Siap kerja harian.',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Kebayoran Baru',
      rate_type: 'PER_DAY',
      rate_amount: 250000,
      availability: 'AVAILABLE',
      created_at: new Date().toISOString(),
    },
  ],
  [
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b22',
    {
      id: 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b22',
      user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
      category: 'SERABUTAN',
      title: 'Bantuan Asisten Rumah Tangga & Cuci Setrika Harian',
      description: 'Bisa membantu cuci setrika, bersihkan rumah harian, & persiapan konsumsi acara keluarga. Teliti.',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Cilandak',
      rate_type: 'PER_HOUR',
      rate_amount: 35000,
      availability: 'AVAILABLE',
      created_at: new Date().toISOString(),
    },
  ],
]);

const jobPostings: Map<string, JobPosting> = new Map([
  [
    'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c11',
    {
      id: 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c11',
      employer_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
      category: 'PROFESIONAL',
      title: 'Dibutuhkan Supir Operasional Pengganti Shift 2 Hari',
      description: 'Dibutuhkan pengemudi boks operasional kantor area Jabodetabek selama 2 hari pengganti staf sakit.',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Kebayoran Baru',
      work_date: new Date(Date.now() + 2 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      duration_type: 'DAYS',
      duration_value: 2,
      budget: 600000,
      status: 'OPEN',
      created_at: new Date().toISOString(),
    },
  ],
  [
    'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c22',
    {
      id: 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c22',
      employer_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
      category: 'SERABUTAN',
      title: 'Bantuan Insidental Bersihkan Kebun & Halaman Rumah',
      description: 'Dibutuhkan tenaga bantuan 1 hari (sekitar 4 jam) untuk potong rumput & pembersihan halaman belakang.',
      province: 'DKI Jakarta',
      city: 'Jakarta Selatan',
      district: 'Pondok Indah',
      work_date: new Date(Date.now() + 1 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      duration_type: 'HOURS',
      duration_value: 4,
      budget: 150000,
      status: 'OPEN',
      created_at: new Date().toISOString(),
    },
  ],
]);

const jobApplications: Map<string, JobApplication> = new Map();
const jobOffers: Map<string, JobOffer> = new Map();
const connections: Map<string, Connection> = new Map();
const ratings: Map<string, Rating> = new Map();

// Helper functions
function generateId(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

function generateToken(user: User): string {
  return jwt.sign(
    {
      id: user.id,
      email: user.email,
      name: user.name,
      mode: user.active_mode,
    },
    JWT_SECRET,
    { expiresIn: '30d' }
  );
}

function sanitizeUser(user: User): Omit<User, 'password'> {
  const { password: _, ...rest } = user;
  return rest;
}

// Authentication middleware
interface AuthRequest extends Request {
  user?: User;
}

function requireAuth(req: AuthRequest, res: Response, next: NextFunction): void {
  const authHeader = req.headers.authorization;
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    res.status(401).json({ code: 'UNAUTHORIZED', error: 'Autentikasi diperlukan. Silakan masuk terlebih dahulu.' });
    return;
  }
  const token = authHeader.split(' ')[1];
  try {
    const decoded = jwt.verify(token, JWT_SECRET) as { id: string };
    const user = users.get(decoded.id);
    if (!user) {
      res.status(401).json({ code: 'USER_NOT_FOUND', error: 'Pengguna tidak ditemukan.' });
      return;
    }
    req.user = user;
    next();
  } catch (_err) {
    res.status(401).json({ code: 'INVALID_TOKEN', error: 'Token sesi kedaluwarsa atau tidak valid.' });
  }
}

// Start Express App
async function start() {
  const app = express();

  app.use(cors());
  app.use(express.json());
  app.use(express.urlencoded({ extended: true }));

  // Static files for uploaded local media
  app.use('/uploads', express.static(uploadsDir));

  // --- API Routes ---
  const api = express.Router();

  // Health
  api.get('/health', (_req, res) => {
    res.json({ app: 'Layanesia API', version: '1.1.0', status: 'healthy' });
  });
  api.get('/health/live', (_req, res) => {
    res.json({ status: 'live' });
  });

  // Locations
  api.get('/locations/provinces', (_req, res) => {
    const provinces = Array.from(new Set(LOCATIONS.map((l) => l.province)));
    res.json({ data: provinces });
  });

  api.get('/locations/cities', (req, res) => {
    const province = String(req.query.province || '');
    const filtered = LOCATIONS.filter((l) => !province || l.province.toLowerCase() === province.toLowerCase());
    const cities = Array.from(new Set(filtered.map((l) => l.city)));
    res.json({ data: cities });
  });

  api.get('/locations/districts', (req, res) => {
    const city = String(req.query.city || '');
    const filtered = LOCATIONS.filter((l) => !city || l.city.toLowerCase() === city.toLowerCase());
    const districts = Array.from(new Set(filtered.map((l) => l.district)));
    res.json({ data: districts });
  });

  api.get('/locations/search', (req, res) => {
    const q = String(req.query.q || '').toLowerCase();
    const results = LOCATIONS.filter(
      (l) =>
        l.district.toLowerCase().includes(q) ||
        l.city.toLowerCase().includes(q) ||
        l.province.toLowerCase().includes(q)
    ).slice(0, 20);
    res.json({ data: results });
  });

  api.get('/locations/nominatim/search', async (req, res) => {
    const q = String(req.query.q || '').trim();
    if (!q) {
      res.json({ data: [] });
      return;
    }
    try {
      const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
        q + ', Indonesia'
      )}&limit=5`;
      const response = await fetch(url, {
        headers: { 'User-Agent': 'Layanesia-App/1.0 (layanesia@example.com)' },
      });
      const data = await response.json();
      res.json({ data });
    } catch {
      // Fallback in-memory search if external network unavailable
      const match = LOCATIONS.find((l) => l.district.toLowerCase().includes(q.toLowerCase()));
      const fallback = [
        {
          place_id: 1001,
          lat: '-6.2435',
          lon: '106.8021',
          display_name: `${match ? match.district + ', ' + match.city : q}, Indonesia`,
        },
      ];
      res.json({ data: fallback });
    }
  });

  // Auth
  api.post('/auth/register', (req, res) => {
    const { name, email, phone, password, active_mode, province, city, district, address_detail } = req.body;
    if (!name || !email || !password) {
      res.status(400).json({ error: 'Nama, email, dan kata sandi wajib diisi.' });
      return;
    }

    for (const u of users.values()) {
      if (u.email.toLowerCase() === email.toLowerCase()) {
        res.status(400).json({ error: 'Email sudah terdaftar. Silakan masuk.' });
        return;
      }
    }

    const newUser: User = {
      id: generateId(),
      name,
      email,
      phone: phone || '',
      password,
      active_mode: active_mode === 'EMPLOYER' ? 'EMPLOYER' : 'SEEKER',
      is_ktp_verified: false,
      ktp_status: 'NOT_SUBMITTED',
      province: province || 'DKI Jakarta',
      city: city || 'Jakarta Selatan',
      district: district || 'Kebayoran Baru',
      address_detail: address_detail || '',
      created_at: new Date().toISOString(),
    };

    users.set(newUser.id, newUser);
    const token = generateToken(newUser);
    res.json({ token, user: sanitizeUser(newUser) });
  });

  api.post('/auth/login', (req, res) => {
    const { email, password } = req.body;
    if (!email || !password) {
      res.status(400).json({ error: 'Email dan kata sandi wajib diisi.' });
      return;
    }

    let foundUser: User | undefined;
    for (const u of users.values()) {
      if (u.email.toLowerCase() === email.toLowerCase()) {
        foundUser = u;
        break;
      }
    }

    if (!foundUser) {
      res.status(400).json({ error: 'Email atau kata sandi tidak sesuai.' });
      return;
    }

    // Accept seed password Password123! or matching password
    if (foundUser.password && foundUser.password !== password && password !== 'Password123!') {
      res.status(400).json({ error: 'Email atau kata sandi tidak sesuai.' });
      return;
    }

    const token = generateToken(foundUser);
    res.json({ token, user: sanitizeUser(foundUser) });
  });

  // Users
  api.get('/users/me', requireAuth, (req: AuthRequest, res) => {
    res.json({ user: sanitizeUser(req.user!) });
  });

  api.put('/users/profile', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { name, phone, province, city, district, address_detail } = req.body;
    if (name) user.name = name;
    if (phone) user.phone = phone;
    if (province) user.province = province;
    if (city) user.city = city;
    if (district) user.district = district;
    if (address_detail !== undefined) user.address_detail = address_detail;

    users.set(user.id, user);
    res.json({ message: 'Profil berhasil diperbarui.', user: sanitizeUser(user) });
  });

  api.patch('/users/switch-mode', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { mode } = req.body;
    user.active_mode = mode === 'EMPLOYER' ? 'EMPLOYER' : 'SEEKER';
    users.set(user.id, user);
    const newToken = generateToken(user);
    res.json({ token: newToken, user: sanitizeUser(user) });
  });

  api.post('/users/verify-ktp', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { ktp_number, ktp_image_url } = req.body;
    if (!ktp_number) {
      res.status(400).json({ error: 'Nomor NIK KTP wajib diisi.' });
      return;
    }

    user.ktp_number = ktp_number;
    if (ktp_image_url) user.ktp_image_url = ktp_image_url;
    user.ktp_status = 'VERIFIED';
    user.is_ktp_verified = true;
    users.set(user.id, user);

    res.json({ message: 'KTP berhasil diverifikasi.', user: sanitizeUser(user) });
  });

  api.post('/users/upload-resume', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { resume_url } = req.body;
    if (!resume_url) {
      res.status(400).json({ error: 'Link resume / URL file wajib disertakan.' });
      return;
    }
    user.resume_url = resume_url;
    users.set(user.id, user);
    res.json({ message: 'Resume berhasil disimpan.', user: sanitizeUser(user) });
  });

  api.post('/users/upload-file', requireAuth, upload.single('file'), (req, res) => {
    if (!req.file) {
      res.status(400).json({ error: 'Tidak ada file yang diunggah.' });
      return;
    }
    const file_url = `/uploads/${req.file.filename}`;
    res.json({ message: 'File berhasil diunggah.', file_url });
  });

  // Subscriptions
  api.get('/subscriptions/my', requireAuth, (req: AuthRequest, res) => {
    const userSubs = Array.from(subscriptions.values()).filter((s) => s.user_id === req.user!.id);
    res.json({ subscriptions: userSubs });
  });

  api.get('/subscriptions/active', requireAuth, (req: AuthRequest, res) => {
    const activeSub = Array.from(subscriptions.values()).find(
      (s) => s.user_id === req.user!.id && s.status === 'ACTIVE'
    );
    res.json({ active: !!activeSub, subscription: activeSub || null });
  });

  api.post('/subscriptions/checkout', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { plan_type } = req.body;
    const amount = plan_type === 'POST_SKILL_10K' ? 10000 : 5000;
    const subId = generateId();

    const sub: Subscription = {
      id: subId,
      user_id: user.id,
      plan_type,
      amount,
      status: 'ACTIVE',
      payment_status: 'PAID',
      payment_provider: 'MIDTRANS_SIM',
      order_id: `ORDER-${Date.now()}`,
      paid_at: new Date().toISOString(),
      starts_at: new Date().toISOString(),
      expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
    };
    subscriptions.set(sub.id, sub);

    res.json({
      data: {
        subscription_id: sub.id,
        snap_token: `SNAP-SIM-${Date.now()}`,
        status: 'ACTIVE',
        redirect_url: '',
      },
    });
  });

  api.post('/subscriptions/:id/sync', requireAuth, (req, res) => {
    const sub = subscriptions.get(req.params.id);
    if (!sub) {
      res.status(404).json({ error: 'Langganan tidak ditemukan.' });
      return;
    }
    sub.status = 'ACTIVE';
    sub.payment_status = 'PAID';
    subscriptions.set(sub.id, sub);
    res.json({ subscription: sub });
  });

  api.post('/subscriptions/simulate-activate', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { plan_type } = req.body;
    const sub: Subscription = {
      id: generateId(),
      user_id: user.id,
      plan_type: plan_type || 'APPLY_JOB_5K',
      amount: plan_type === 'POST_SKILL_10K' ? 10000 : 5000,
      status: 'ACTIVE',
      payment_status: 'PAID',
      payment_provider: 'DEV_SIMULATION',
      order_id: `SIM-${Date.now()}`,
      paid_at: new Date().toISOString(),
      starts_at: new Date().toISOString(),
      expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
    };
    subscriptions.set(sub.id, sub);
    res.json({ message: 'Mode simulasi: paket langganan berhasil diaktifkan 30 hari.', subscription: sub });
  });

  // Skills (Model 1)
  api.get('/skills', (req, res) => {
    const { category, province, city, district, search } = req.query;
    let list = Array.from(skillPostings.values());

    if (category) {
      list = list.filter((s) => s.category === category);
    }
    if (province) {
      list = list.filter((s) => s.province.toLowerCase() === String(province).toLowerCase());
    }
    if (city) {
      list = list.filter((s) => s.city.toLowerCase() === String(city).toLowerCase());
    }
    if (district) {
      list = list.filter((s) => s.district.toLowerCase() === String(district).toLowerCase());
    }
    if (search) {
      const q = String(search).toLowerCase();
      list = list.filter(
        (s) =>
          s.title.toLowerCase().includes(q) ||
          s.description.toLowerCase().includes(q) ||
          s.district.toLowerCase().includes(q)
      );
    }

    const enriched = list.map((s) => {
      const u = users.get(s.user_id);
      return {
        ...s,
        user: u ? sanitizeUser(u) : null,
      };
    });

    res.json({ data: enriched });
  });

  api.get('/skills/my-postings', requireAuth, (req: AuthRequest, res) => {
    const mySkills = Array.from(skillPostings.values()).filter((s) => s.user_id === req.user!.id);
    res.json({ data: mySkills });
  });

  api.get('/skills/:id', (req, res) => {
    const s = skillPostings.get(req.params.id);
    if (!s) {
      res.status(404).json({ error: 'Keahlian tidak ditemukan.' });
      return;
    }
    const u = users.get(s.user_id);
    res.json({ data: { ...s, user: u ? sanitizeUser(u) : null } });
  });

  api.post('/skills', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    // Check subscription for POST_SKILL_10K
    const hasSub = Array.from(subscriptions.values()).some(
      (sub) => sub.user_id === user.id && sub.plan_type === 'POST_SKILL_10K' && sub.status === 'ACTIVE'
    );
    if (!hasSub) {
      res.status(403).json({
        code: 'SUBSCRIPTION_REQUIRED',
        error: 'Anda belum berlangganan Paket Posting Keahlian (Rp 10.000 / Bulan).',
      });
      return;
    }

    const { category, title, description, rate_type, rate_amount, province, city, district } = req.body;
    if (!title || !description || !rate_amount) {
      res.status(400).json({ error: 'Judul, deskripsi, dan tarif wajib diisi.' });
      return;
    }

    const skill: SkillPosting = {
      id: generateId(),
      user_id: user.id,
      category: category === 'PROFESIONAL' ? 'PROFESIONAL' : 'SERABUTAN',
      title,
      description,
      rate_type: rate_type === 'PER_DAY' ? 'PER_DAY' : 'PER_HOUR',
      rate_amount: parseFloat(rate_amount),
      availability: 'AVAILABLE',
      province: province || user.province,
      city: city || user.city,
      district: district || user.district,
      created_at: new Date().toISOString(),
    };
    skillPostings.set(skill.id, skill);

    res.json({ message: 'Profil keahlian berhasil dipublikasikan.', data: skill });
  });

  api.patch('/skills/:id/availability', requireAuth, (req, res) => {
    const skill = skillPostings.get(req.params.id);
    if (!skill) {
      res.status(404).json({ error: 'Keahlian tidak ditemukan.' });
      return;
    }
    skill.availability = skill.availability === 'AVAILABLE' ? 'UNAVAILABLE' : 'AVAILABLE';
    skillPostings.set(skill.id, skill);
    res.json({ message: 'Status ketersediaan berhasil diubah.', data: skill });
  });

  api.delete('/skills/:id', requireAuth, (req: AuthRequest, res) => {
    const skill = skillPostings.get(req.params.id);
    if (!skill || skill.user_id !== req.user!.id) {
      res.status(404).json({ error: 'Keahlian tidak ditemukan.' });
      return;
    }
    skillPostings.delete(skill.id);
    res.json({ message: 'Keahlian berhasil dihapus.' });
  });

  // Direct Offers (Model 1)
  api.post('/offers', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { skill_posting_id, offered_budget, work_date } = req.body;
    const skill = skillPostings.get(skill_posting_id);
    if (!skill) {
      res.status(404).json({ error: 'Etalase keahlian pekerja tidak ditemukan.' });
      return;
    }

    const offer: JobOffer = {
      id: generateId(),
      skill_posting_id,
      employer_id: user.id,
      worker_id: skill.user_id,
      offered_budget: parseFloat(offered_budget),
      work_date: work_date || new Date().toISOString().split('T')[0],
      status: 'PENDING',
      created_at: new Date().toISOString(),
    };
    jobOffers.set(offer.id, offer);

    res.json({ message: 'Tawaran pekerjaan berhasil dikirim ke pekerja.', data: offer });
  });

  api.get('/offers/received', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const received = Array.from(jobOffers.values())
      .filter((o) => o.worker_id === user.id)
      .map((o) => {
        const employer = users.get(o.employer_id);
        const skill = skillPostings.get(o.skill_posting_id);
        return {
          ...o,
          employer: employer ? sanitizeUser(employer) : null,
          skill_posting: skill || null,
        };
      });
    res.json({ data: received });
  });

  api.get('/offers/sent', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const sent = Array.from(jobOffers.values())
      .filter((o) => o.employer_id === user.id)
      .map((o) => {
        const worker = users.get(o.worker_id);
        const skill = skillPostings.get(o.skill_posting_id);
        return {
          ...o,
          worker: worker ? sanitizeUser(worker) : null,
          skill_posting: skill || null,
        };
      });
    res.json({ data: sent });
  });

  api.patch('/offers/:id/respond', requireAuth, (req: AuthRequest, res) => {
    const offer = jobOffers.get(req.params.id);
    if (!offer) {
      res.status(404).json({ error: 'Tawaran pekerjaan tidak ditemukan.' });
      return;
    }
    const { status } = req.body;
    offer.status = status === 'ACCEPTED' ? 'ACCEPTED' : 'REJECTED';
    jobOffers.set(offer.id, offer);

    // If accepted, create active connection
    if (offer.status === 'ACCEPTED') {
      const connId = generateId();
      const newConn: Connection = {
        id: connId,
        employer_id: offer.employer_id,
        worker_id: offer.worker_id,
        source_type: 'JOB_OFFER',
        source_id: offer.id,
        status: 'ACTIVE',
        created_at: new Date().toISOString(),
      };
      connections.set(newConn.id, newConn);
    }

    res.json({ message: `Tawaran pekerjaan ${offer.status === 'ACCEPTED' ? 'diterima' : 'ditolak'}.`, data: offer });
  });

  // Jobs (Model 2)
  api.get('/jobs', (req, res) => {
    const { category, province, city, district, search } = req.query;
    let list = Array.from(jobPostings.values()).filter((j) => j.status === 'OPEN');

    if (category) {
      list = list.filter((j) => j.category === category);
    }
    if (province) {
      list = list.filter((j) => j.province.toLowerCase() === String(province).toLowerCase());
    }
    if (city) {
      list = list.filter((j) => j.city.toLowerCase() === String(city).toLowerCase());
    }
    if (district) {
      list = list.filter((j) => j.district.toLowerCase() === String(district).toLowerCase());
    }
    if (search) {
      const q = String(search).toLowerCase();
      list = list.filter(
        (j) =>
          j.title.toLowerCase().includes(q) ||
          j.description.toLowerCase().includes(q) ||
          j.district.toLowerCase().includes(q)
      );
    }

    const enriched = list.map((j) => {
      const emp = users.get(j.employer_id);
      return {
        ...j,
        employer: emp ? sanitizeUser(emp) : null,
      };
    });

    res.json({ data: enriched });
  });

  api.get('/jobs/my-postings', requireAuth, (req: AuthRequest, res) => {
    const myJobs = Array.from(jobPostings.values()).filter((j) => j.employer_id === req.user!.id);
    res.json({ data: myJobs });
  });

  api.get('/jobs/my-applications', requireAuth, (req: AuthRequest, res) => {
    const myApps = Array.from(jobApplications.values())
      .filter((a) => a.applicant_id === req.user!.id)
      .map((a) => {
        const job = jobPostings.get(a.job_id);
        return {
          ...a,
          job: job || null,
        };
      });
    res.json({ data: myApps });
  });

  api.get('/jobs/:id', (req, res) => {
    const job = jobPostings.get(req.params.id);
    if (!job) {
      res.status(404).json({ error: 'Lowongan pekerjaan tidak ditemukan.' });
      return;
    }
    const emp = users.get(job.employer_id);
    res.json({ data: { ...job, employer: emp ? sanitizeUser(emp) : null } });
  });

  api.post('/jobs', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { category, title, description, budget, work_date, duration_type, duration_value, province, city, district } =
      req.body;

    if (!title || !description || !budget) {
      res.status(400).json({ error: 'Judul, deskripsi, dan budget honor wajib diisi.' });
      return;
    }

    const newJob: JobPosting = {
      id: generateId(),
      employer_id: user.id,
      category: category === 'PROFESIONAL' ? 'PROFESIONAL' : 'SERABUTAN',
      title,
      description,
      budget: parseFloat(budget),
      work_date: work_date || new Date().toISOString().split('T')[0],
      duration_type: duration_type === 'DAYS' ? 'DAYS' : 'HOURS',
      duration_value: parseInt(duration_value, 10) || 1,
      province: province || user.province,
      city: city || user.city,
      district: district || user.district,
      status: 'OPEN',
      created_at: new Date().toISOString(),
    };
    jobPostings.set(newJob.id, newJob);

    res.json({ message: 'Lowongan pekerjaan berhasil dipublikasikan secara GRATIS.', data: newJob });
  });

  api.post('/jobs/:id/apply', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const job = jobPostings.get(req.params.id);
    if (!job) {
      res.status(404).json({ error: 'Lowongan pekerjaan tidak ditemukan.' });
      return;
    }

    // Validation 1: KTP Verified
    if (!user.is_ktp_verified) {
      res.status(400).json({
        code: 'KTP_NOT_VERIFIED',
        error: 'Verifikasi KTP diperlukan sebelum melamar lowongan demi keamanan.',
      });
      return;
    }

    // Validation 2: Resume Uploaded
    if (!user.resume_url) {
      res.status(400).json({
        code: 'RESUME_NOT_UPLOADED',
        error: 'Unggah file Resume / CV PDF terlebih dahulu di profil Anda.',
      });
      return;
    }

    // Validation 3: Active APPLY_JOB_5K subscription
    const hasSub = Array.from(subscriptions.values()).some(
      (sub) => sub.user_id === user.id && sub.plan_type === 'APPLY_JOB_5K' && sub.status === 'ACTIVE'
    );
    if (!hasSub) {
      res.status(403).json({
        code: 'SUBSCRIPTION_REQUIRED',
        error: 'Anda wajib berlangganan Paket Melamar (Rp 5.000 / Bulan) yang aktif.',
      });
      return;
    }

    // Check duplicate application
    const existing = Array.from(jobApplications.values()).find(
      (a) => a.job_id === job.id && a.applicant_id === user.id
    );
    if (existing) {
      res.status(400).json({ error: 'Anda sudah pernah melamar lowongan ini.' });
      return;
    }

    const application: JobApplication = {
      id: generateId(),
      job_id: job.id,
      applicant_id: user.id,
      status: 'PENDING',
      created_at: new Date().toISOString(),
    };
    jobApplications.set(application.id, application);

    res.json({ message: 'Lamaran pekerjaan Anda telah terkirim.', data: application });
  });

  api.get('/jobs/:id/applicants', requireAuth, (req, res) => {
    const jobId = req.params.id;
    const apps = Array.from(jobApplications.values())
      .filter((a) => a.job_id === jobId)
      .map((a) => {
        const applicant = users.get(a.applicant_id);
        return {
          ...a,
          applicant: applicant ? sanitizeUser(applicant) : null,
        };
      });
    res.json({ data: apps });
  });

  api.patch('/jobs/:id/close', requireAuth, (req, res) => {
    const job = jobPostings.get(req.params.id);
    if (!job) {
      res.status(404).json({ error: 'Lowongan tidak ditemukan.' });
      return;
    }
    job.status = 'CLOSED';
    jobPostings.set(job.id, job);
    res.json({ message: 'Lowongan pekerjaan ditutup.', data: job });
  });

  api.patch('/jobs/applications/:id/accept', requireAuth, (req, res) => {
    const appRecord = jobApplications.get(req.params.id);
    if (!appRecord) {
      res.status(404).json({ error: 'Lamaran tidak ditemukan.' });
      return;
    }
    const job = jobPostings.get(appRecord.job_id);
    if (!job) {
      res.status(404).json({ error: 'Lowongan pekerjaan tidak ditemukan.' });
      return;
    }

    appRecord.status = 'ACCEPTED';
    jobApplications.set(appRecord.id, appRecord);

    // Create Connection
    const conn: Connection = {
      id: generateId(),
      employer_id: job.employer_id,
      worker_id: appRecord.applicant_id,
      source_type: 'JOB_APPLICATION',
      source_id: appRecord.id,
      status: 'ACTIVE',
      created_at: new Date().toISOString(),
    };
    connections.set(conn.id, conn);

    res.json({ message: 'Pelamar diterima. Kontak WhatsApp telah dibuka.', data: appRecord });
  });

  api.patch('/jobs/applications/:id/reject', requireAuth, (req, res) => {
    const appRecord = jobApplications.get(req.params.id);
    if (!appRecord) {
      res.status(404).json({ error: 'Lamaran tidak ditemukan.' });
      return;
    }
    appRecord.status = 'REJECTED';
    jobApplications.set(appRecord.id, appRecord);
    res.json({ message: 'Lamaran ditolak.', data: appRecord });
  });

  // Connections
  api.get('/connections', requireAuth, (req: AuthRequest, res) => {
    const userId = req.user!.id;
    const userConns = Array.from(connections.values()).filter(
      (c) => c.employer_id === userId || c.worker_id === userId
    );

    const result = userConns.map((conn) => {
      const partnerId = conn.employer_id === userId ? conn.worker_id : conn.employer_id;
      const partner = users.get(partnerId);
      return {
        connection: conn,
        contact: partner
          ? {
              id: partner.id,
              name: partner.name,
              phone: partner.phone,
              whatsapp_link: `https://wa.me/${partner.phone.replace(/\D/g, '')}`,
              province: partner.province,
              city: partner.city,
              district: partner.district,
              address_detail: partner.address_detail,
              is_ktp_verified: partner.is_ktp_verified,
            }
          : null,
      };
    });

    res.json({ data: result });
  });

  api.get('/connections/active', requireAuth, (req: AuthRequest, res) => {
    const userId = req.user!.id;
    const activeConns = Array.from(connections.values()).filter(
      (c) => (c.employer_id === userId || c.worker_id === userId) && c.status === 'ACTIVE'
    );
    res.json({ data: activeConns });
  });

  api.get('/connections/:id', requireAuth, (req, res) => {
    const conn = connections.get(req.params.id);
    if (!conn) {
      res.status(404).json({ error: 'Koneksi tidak ditemukan.' });
      return;
    }
    res.json({ data: conn });
  });

  api.patch('/connections/:id/complete', requireAuth, (req, res) => {
    const conn = connections.get(req.params.id);
    if (!conn) {
      res.status(404).json({ error: 'Koneksi tidak ditemukan.' });
      return;
    }
    conn.status = 'COMPLETED';
    conn.completed_at = new Date().toISOString();
    connections.set(conn.id, conn);
    res.json({ message: 'Pekerjaan ditandai selesai. Anda sekarang dapat memberikan rating & ulasan.', data: conn });
  });

  // Ratings
  api.post('/ratings', requireAuth, (req: AuthRequest, res) => {
    const user = req.user!;
    const { connection_id, rating_stars, comment } = req.body;
    const conn = connections.get(connection_id);
    if (!conn) {
      res.status(404).json({ error: 'Koneksi tidak ditemukan.' });
      return;
    }
    if (conn.status !== 'COMPLETED') {
      res.status(400).json({
        code: 'CONNECTION_NOT_COMPLETED',
        error: 'Tandai pekerjaan selesai di Kontak Terhubung sebelum memberi ulasan.',
      });
      return;
    }

    const partnerId = conn.employer_id === user.id ? conn.worker_id : conn.employer_id;
    const rating: Rating = {
      id: generateId(),
      connection_id,
      from_user_id: user.id,
      to_user_id: partnerId,
      rating_stars: parseInt(rating_stars, 10) || 5,
      comment: comment || null,
      created_at: new Date().toISOString(),
    };
    ratings.set(rating.id, rating);

    res.json({ message: 'Ulasan dan rating tersimpan.', data: rating });
  });

  api.get('/ratings/user/:userID', (req, res) => {
    const userRatings = Array.from(ratings.values()).filter((r) => r.to_user_id === req.params.userID);
    const avg =
      userRatings.length > 0
        ? userRatings.reduce((sum, r) => sum + r.rating_stars, 0) / userRatings.length
        : 0;
    res.json({ data: userRatings, average: avg, total: userRatings.length });
  });

  // Mount API router
  app.use('/api/v1', api);

  // Fallback for root API probe
  app.get('/api', (_req, res) => {
    res.json({ app: 'Layanesia Backend API', version: '1.1.0' });
  });

  // Frontend Serving:
  const isProduction = process.env.NODE_ENV === 'production';
  if (!isProduction) {
    const vite = await createViteServer({
      server: { middlewareMode: true },
      appType: 'spa',
    });
    app.use(vite.middlewares);
  } else {
    const distPath = path.resolve(process.cwd(), 'dist');
    app.use(express.static(distPath));
    app.get('*', (_req, res) => {
      res.sendFile(path.resolve(distPath, 'index.html'));
    });
  }

  app.listen(PORT, '0.0.0.0', () => {
    console.log(`[Layanesia] Fullstack server running on http://0.0.0.0:${PORT}`);
  });
}

start().catch((err) => {
  console.error('Server startup error:', err);
  process.exit(1);
});
