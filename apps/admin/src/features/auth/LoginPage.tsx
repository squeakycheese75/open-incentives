import { useState, type FormEvent } from "react";
import { Navigate, useLocation } from "react-router-dom";

import { env } from "../../lib/env";
import { ApiError } from "../../api/client";
import { useAuth } from "./AuthContext";
import { useLogin } from "./api";

export function LoginPage() {
  const { isAuthenticated, login } = useAuth();
  const location = useLocation();
  const loginMutation = useLogin();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  if (isAuthenticated) {
    const redirectTo = (location.state as { from?: string } | null)?.from ?? "/projects";
    return <Navigate to={redirectTo} replace />;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    try {
      const res = await loginMutation.mutateAsync({
        email,
        password,
        organization: env.defaultOrgId,
      });
      login({ token: res.token, expiresAt: res.expiresAt });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to log in");
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-sm rounded-lg border border-gray-200 bg-white p-8 shadow-sm">
        <h1 className="mb-1 text-xl font-semibold text-gray-900">Open Incentives</h1>
        <p className="mb-6 text-sm text-gray-500">Sign in to manage your campaigns</p>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="email" className="mb-1 block text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              id="email"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full rounded border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
            />
          </div>

          <div>
            <label htmlFor="password" className="mb-1 block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              id="password"
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full rounded border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
            />
          </div>

          {error && <p className="text-sm text-red-600">{error}</p>}

          <button
            type="submit"
            disabled={loginMutation.isPending}
            className="w-full rounded bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 disabled:opacity-50"
          >
            {loginMutation.isPending ? "Signing in..." : "Sign in"}
          </button>
        </form>
      </div>
    </div>
  );
}
