const BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:3001';

type TokenStore = {
  access: string;
  refresh: string;
  tenantId?: string;
};

// ponytail: sessionStorage for dev; swap to httpOnly cookie in prod
function getTokens(): TokenStore | null {
  if (typeof window === 'undefined') return null;
  const raw = sessionStorage.getItem('chiguire_tokens');
  return raw ? JSON.parse(raw) : null;
}

function setTokens(t: TokenStore) {
  sessionStorage.setItem('chiguire_tokens', JSON.stringify(t));
}

function clearTokens() {
  sessionStorage.removeItem('chiguire_tokens');
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const tokens = getTokens();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init.headers as Record<string, string>),
  };
  if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;

  const res = await fetch(`${BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${res.status}: ${text}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  auth: {
    async register(email: string, password: string, fullName: string) {
      return request<{ user_id: string }>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password, full_name: fullName }),
      });
    },
    async login(email: string, password: string, tenantId?: string) {
      const data = await request<{ access_token: string; refresh_token: string }>(
        '/auth/login',
        {
          method: 'POST',
          body: JSON.stringify({ email, password, tenant_id: tenantId }),
        },
      );
      setTokens({ access: data.access_token, refresh: data.refresh_token, tenantId });
      return data;
    },
    logout() {
      clearTokens();
    },
    getAccessToken() {
      return getTokens()?.access ?? null;
    },
  },
  tenants: {
    list() {
      return request<{ id: string; name: string; slug: string; plan: string }[]>('/tenants');
    },
    create(name: string, slug: string) {
      return request<{ id: string; name: string; slug: string }>('/tenants', {
        method: 'POST',
        body: JSON.stringify({ name, slug }),
      });
    },
  },
  powersync: {
    getToken() {
      return request<{ token: string; expires_at: number }>('/powersync/token');
    },
  },
  testItems: {
    list() {
      return request<{ id: string; label: string; created_at: string }[]>('/test-items');
    },
    create(label: string) {
      return request<{ id: string; label: string; created_at: string }>('/test-items', {
        method: 'POST',
        body: JSON.stringify({ label }),
      });
    },
  },
};
