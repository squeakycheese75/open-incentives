import { useEffect, useState, type FormEvent } from "react";

import { Dialog } from "../../components/ui/Dialog";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { ApiError } from "../../api/client";
import { useToast } from "../../components/ui/ToastProvider";
import type { Project } from "../../api/types";
import { useCreateProject, useUpdateProject } from "./api";

interface ProjectFormDialogProps {
  open: boolean;
  onClose: () => void;
  project?: Project | null;
}

export function ProjectFormDialog({ open, onClose, project }: ProjectFormDialogProps) {
  const [name, setName] = useState(project?.name ?? "");
  const [error, setError] = useState<string | null>(null);
  const { showToast } = useToast();

  const createProject = useCreateProject();
  const updateProject = useUpdateProject();

  const isEdit = Boolean(project);
  const isPending = createProject.isPending || updateProject.isPending;

  useEffect(() => {
    if (open) {
      setName(project?.name ?? "");
      setError(null);
    }
  }, [open, project]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    try {
      if (isEdit && project) {
        await updateProject.mutateAsync({ projectId: project.publicId, body: { name } });
        showToast("Project updated");
      } else {
        await createProject.mutateAsync({ name });
        showToast("Project created");
      }
      onClose();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong");
    }
  }

  return (
    <Dialog open={open} onClose={onClose} title={isEdit ? "Rename project" : "Create project"}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="project-name" className="mb-1 block text-sm font-medium text-gray-700">
            Name
          </label>
          <Input id="project-name" required value={name} onChange={(e) => setName(e.target.value)} autoFocus />
        </div>

        {error && <p className="text-sm text-red-600">{error}</p>}

        <div className="flex justify-end gap-2 pt-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? "Saving..." : "Save"}
          </Button>
        </div>
      </form>
    </Dialog>
  );
}
