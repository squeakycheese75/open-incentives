import { useNavigate } from "react-router-dom";

import { Select } from "../../components/ui/Select";
import { useProjects } from "./api";
import { useCurrentProject } from "./useCurrentProject";

export function ProjectSwitcher() {
  const navigate = useNavigate();
  const { data, isLoading } = useProjects();
  const { projectId } = useCurrentProject();

  if (isLoading) return null;

  const projects = data?.projects ?? [];
  if (projects.length === 0) return null;

  return (
    <Select
      className="w-56"
      value={projectId ?? ""}
      onChange={(e) => navigate(`/projects/${e.target.value}/campaigns`)}
    >
      {!projectId && <option value="">Select a project…</option>}
      {projects.map((p) => (
        <option key={p.publicId} value={p.publicId}>
          {p.name}
        </option>
      ))}
    </Select>
  );
}
