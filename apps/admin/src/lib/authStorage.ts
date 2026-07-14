const STORAGE_KEY = "oi_admin_token";

export interface StoredToken {
  token: string;
  expiresAt: string;
}

export function getStoredToken(): StoredToken | null {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return null;

  try {
    return JSON.parse(raw) as StoredToken;
  } catch {
    return null;
  }
}

export function setStoredToken(value: StoredToken): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
}

export function clearStoredToken(): void {
  localStorage.removeItem(STORAGE_KEY);
}

export function isTokenExpired(stored: StoredToken): boolean {
  return new Date(stored.expiresAt).getTime() <= Date.now();
}
