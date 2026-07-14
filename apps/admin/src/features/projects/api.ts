import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { apiFetch } from "../../api/client";
import type { CreateProjectRequest, Project, UpdateProjectRequest } from "../../api/types";

const projectsKey = ["projects"] as const;

export function useProjects() {
  return useQuery({
    queryKey: projectsKey,
    queryFn: () => apiFetch<{ projects: Project[] }>("/admin/projects"),
  });
}

export function useCreateProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateProjectRequest) =>
      apiFetch<Project>("/admin/projects", { method: "POST", body: JSON.stringify(body) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: projectsKey }),
  });
}

export function useUpdateProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, body }: { projectId: string; body: UpdateProjectRequest }) =>
      apiFetch<Project>(`/admin/projects/${projectId}`, { method: "PATCH", body: JSON.stringify(body) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: projectsKey }),
  });
}

export function useDeleteProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (projectId: string) => apiFetch<void>(`/admin/projects/${projectId}`, { method: "DELETE" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: projectsKey }),
  });
}
