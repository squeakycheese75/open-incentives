import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { apiFetch } from "../../api/client";
import type { Campaign, CreateCampaignRequest, UpdateCampaignRequest } from "../../api/types";

function campaignsKey(projectId: string) {
  return ["campaigns", projectId] as const;
}

function campaignKey(projectId: string, campaignId: string) {
  return ["campaigns", projectId, campaignId] as const;
}

export function useCampaigns(projectId: string) {
  return useQuery({
    queryKey: campaignsKey(projectId),
    queryFn: () => apiFetch<{ campaigns: Campaign[] }>(`/admin/projects/${projectId}/campaigns`),
    enabled: !!projectId,
  });
}

export function useCampaign(projectId: string, campaignId: string) {
  return useQuery({
    queryKey: campaignKey(projectId, campaignId),
    queryFn: () => apiFetch<Campaign>(`/admin/projects/${projectId}/campaigns/${campaignId}`),
    enabled: !!projectId && !!campaignId,
  });
}

export function useCreateCampaign(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateCampaignRequest) =>
      apiFetch<{ campaignPublicId: string }>(`/admin/projects/${projectId}/campaigns`, {
        method: "POST",
        body: JSON.stringify(body),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: campaignsKey(projectId) }),
  });
}

export function useUpdateCampaign(projectId: string, campaignId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: UpdateCampaignRequest) =>
      apiFetch<Campaign>(`/admin/projects/${projectId}/campaigns/${campaignId}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: campaignsKey(projectId) });
      qc.invalidateQueries({ queryKey: campaignKey(projectId, campaignId) });
    },
  });
}

export function useDeleteCampaign(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (campaignId: string) =>
      apiFetch<void>(`/admin/projects/${projectId}/campaigns/${campaignId}`, { method: "DELETE" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: campaignsKey(projectId) }),
  });
}
