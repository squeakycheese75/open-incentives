import { useParams } from "react-router-dom";

import { useProjects } from "./api";

export function useCurrentProject() {
  const { projectId } = useParams<{ projectId: string }>();
  const { data } = useProjects();

  const project = data?.projects.find((p) => p.publicId === projectId) ?? null;

  return { projectId: projectId ?? null, project };
}
