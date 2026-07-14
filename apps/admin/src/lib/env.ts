function required(name: keyof ImportMetaEnv): string {
  const value = import.meta.env[name];
  if (!value) {
    throw new Error(`Missing required environment variable: ${name}`);
  }
  return value;
}

export const env = {
  apiBaseUrl: required("VITE_API_BASE_URL"),
  defaultOrgId: required("VITE_DEFAULT_ORG_ID"),
};
