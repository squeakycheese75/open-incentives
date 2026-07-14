import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { apiFetch } from "../../api/client";
import type { ApiKey, CreateApiKeyRequest, CreateApiKeyResponse } from "../../api/types";

function apiKeysKey(projectId: string) {
  return ["apiKeys", projectId] as const;
}

export function useApiKeys(projectId: string) {
  return useQuery({
    queryKey: apiKeysKey(projectId),
    queryFn: () => apiFetch<{ apiKeys: ApiKey[] }>(`/admin/projects/${projectId}/api-keys`),
    enabled: !!projectId,
  });
}

export function useCreateApiKey(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateApiKeyRequest) =>
      apiFetch<CreateApiKeyResponse>(`/admin/projects/${projectId}/api-keys`, {
        method: "POST",
        body: JSON.stringify(body),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: apiKeysKey(projectId) }),
  });
}

export function useRevokeApiKey(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (apiKeyPublicId: string) =>
      apiFetch<ApiKey>(`/admin/projects/${projectId}/api-keys/${apiKeyPublicId}/revoke`, { method: "POST" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: apiKeysKey(projectId) }),
  });
}
