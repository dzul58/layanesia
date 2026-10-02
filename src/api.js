const API_BASE = String(import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')

export class ApiError extends Error {
  constructor(message, { status = 0, code = '', details } = {}) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}

export function apiUrl(path = '') {
  if (/^https?:\/\//i.test(path)) return path
  const p = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${p}`
}

/** Ubah path relatif hasil upload lokal (`/uploads/...`) menjadi URL yang bisa dibuka browser. */
export function fileUrl(path) {
  if (!path) return ''
  if (/^https?:\/\//i.test(path) || path.startsWith('blob:')) return path
  return apiUrl(path)
}

export function whatsappUrl(linkOrPhone) {
  if (!linkOrPhone) return '#'
  if (/^https?:\/\//i.test(linkOrPhone)) return linkOrPhone
  const digits = String(linkOrPhone).replace(/\D/g, '')
  return digits ? `https://wa.me/${digits}` : '#'
}

async function parseBody(res) {
  const text = await res.text()
  if (!text) return {}
  try {
    return JSON.parse(text)
  } catch {
    return { error: text }
  }
}

/**
 * @param {string} path  misalnya `/api/v1/users/me`
 * @param {{ method?: string, token?: string, body?: any, formData?: FormData }} [opts]
 */
export async function api(path, opts = {}) {
  const { method = 'GET', token, body, formData } = opts
  const headers = {}
  if (token) headers.Authorization = `Bearer ${token}`

  const init = { method, headers }
  if (formData) {
    init.body = formData
  } else if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
    init.body = JSON.stringify(body)
  }

  const res = await fetch(apiUrl(path), init)
  const data = await parseBody(res)
  if (!res.ok) {
    throw new ApiError(data.error || data.message || 'Permintaan gagal.', {
      status: res.status,
      code: data.code || '',
      details: data.details,
    })
  }
  return data
}

let snapLoader = null

export function loadMidtransSnap() {
  if (typeof window !== 'undefined' && window.snap) {
    return Promise.resolve(window.snap)
  }
  const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY
  if (!clientKey) return Promise.resolve(null)
  if (snapLoader) return snapLoader

  snapLoader = new Promise((resolve, reject) => {
    const isProd = String(import.meta.env.VITE_MIDTRANS_IS_PRODUCTION) === 'true'
    const script = document.createElement('script')
    script.src = isProd
      ? 'https://app.midtrans.com/snap/snap.js'
      : 'https://app.sandbox.midtrans.com/snap/snap.js'
    script.setAttribute('data-client-key', clientKey)
    script.onload = () => resolve(window.snap || null)
    script.onerror = () => {
      snapLoader = null
      reject(new Error('Gagal memuat Midtrans Snap.'))
    }
    document.head.appendChild(script)
  })
  return snapLoader
}

export function ktpStatusLabel(user) {
  const status = user?.ktp_status
  if (user?.is_ktp_verified || status === 'VERIFIED') return 'Terverifikasi'
  if (status === 'PENDING_REVIEW') return 'Menunggu tinjauan'
  if (status === 'REJECTED') return 'Ditolak'
  return 'Belum diajukan'
}
