import { createContext, useContext, useMemo, useState, type ReactNode } from "react";

import { clearStoredToken, getStoredToken, isTokenExpired, setStoredToken, type StoredToken } from "../../lib/authStorage";

// Today the API is single-tenant: `organization` is hardcoded from
// VITE_DEFAULT_ORG_ID at login time (see LoginPage). Adding a real org
// switcher later means: (1) adding an orgId to the URL param structure
// alongside :projectId, (2) an useOrgs()/OrgSwitcher mirroring
// ProjectSwitcher, (3) sourcing `organization` from that selection instead
// of the env default.

interface AuthContextValue {
  isAuthenticated: boolean;
  login: (token: StoredToken) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function hasValidToken(): boolean {
  const stored = getStoredToken();
  return stored !== null && !isTokenExpired(stored);
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(hasValidToken);

  const value = useMemo<AuthContextValue>(
    () => ({
      isAuthenticated,
      login: (token) => {
        setStoredToken(token);
        setIsAuthenticated(true);
      },
      logout: () => {
        clearStoredToken();
        setIsAuthenticated(false);
      },
    }),
    [isAuthenticated],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
