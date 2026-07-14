export interface LoginRequest {
  email: string;
  password: string;
  organization: string;
}

export interface LoginResponse {
  token: string;
  expiresAt: string;
}

export interface Project {
  publicId: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProjectRequest {
  name: string;
}

export interface UpdateProjectRequest {
  name: string;
}

export type CampaignStatus = "active" | "inactive";

export interface Campaign {
  publicId: string;
  name: string;
  status: CampaignStatus;
  rules: unknown;
}

export interface CreateCampaignRequest {
  name: string;
  status: CampaignStatus;
  rules: unknown;
}

export interface UpdateCampaignRequest {
  name: string;
  status: CampaignStatus;
  rules: unknown;
}

export type ApiKeyStatus = "active" | "revoked";

export interface ApiKey {
  publicId: string;
  name: string;
  prefix: string;
  status: ApiKeyStatus;
  createdAt: string;
  lastUsedAt?: string;
  revokedAt?: string;
}

export interface CreateApiKeyRequest {
  name: string;
  description: string;
}

export interface CreateApiKeyResponse {
  apiKey: string;
  apiKeyPublicId: string;
  createdAt: string;
}
