import { env } from "../lib/env";
import { clearStoredToken, getStoredToken } from "../lib/authStorage";

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

function handleUnauthorized(): void {
  clearStoredToken();
  window.location.href = "/login";
}

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers);

  const stored = getStoredToken();
  if (stored) {
    headers.set("Authorization", `Bearer ${stored.token}`);
  }

  if (options.body) {
    headers.set("Content-Type", "application/json");
  }

  const res = await fetch(`${env.apiBaseUrl}${path}`, {
    ...options,
    headers,
  });

  if (res.status === 204) {
    return undefined as T;
  }

  const body = await res.json().catch(() => null);

  if (!res.ok) {
    const message = body && typeof body.error === "string" ? body.error : `request failed with status ${res.status}`;

    // A 401 on the login call itself just means bad credentials - surface it
    // as a normal error rather than treating it as an expired session.
    if (res.status === 401 && path !== "/admin/auth/login") {
      handleUnauthorized();
    }

    throw new ApiError(res.status, message);
  }

  return body as T;
}
