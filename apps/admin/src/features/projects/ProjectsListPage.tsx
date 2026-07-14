import { useState } from "react";
import { Link } from "react-router-dom";

import { Button } from "../../components/ui/Button";
import { Card } from "../../components/ui/Card";
import { Spinner } from "../../components/ui/Spinner";
import { Table, TableBody, TableCell, TableHead, TableHeaderCell, TableRow } from "../../components/ui/Table";
import { useToast } from "../../components/ui/ToastProvider";
import { ApiError } from "../../api/client";
import type { Project } from "../../api/types";
import { useDeleteProject, useProjects } from "./api";
import { ProjectFormDialog } from "./ProjectFormDialog";

export function ProjectsListPage() {
  const { data, isLoading, isError } = useProjects();
  const deleteProject = useDeleteProject();
  const { showToast } = useToast();

  const [dialogState, setDialogState] = useState<{ open: boolean; project: Project | null }>({
    open: false,
    project: null,
  });

  async function handleDelete(project: Project) {
    if (!window.confirm(`Delete project "${project.name}"? This cannot be undone.`)) return;

    try {
      await deleteProject.mutateAsync(project.publicId);
      showToast("Project deleted");
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : "Failed to delete project", "error");
    }
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-gray-900">Projects</h1>
        <Button onClick={() => setDialogState({ open: true, project: null })}>New project</Button>
      </div>

      <Card>
        {isLoading && (
          <div className="flex justify-center py-10">
            <Spinner />
          </div>
        )}

        {isError && <p className="p-6 text-sm text-red-600">Failed to load projects.</p>}

        {data && data.projects.length === 0 && (
          <p className="p-6 text-sm text-gray-500">No projects yet. Create one to get started.</p>
        )}

        {data && data.projects.length > 0 && (
          <Table>
            <TableHead>
              <TableRow>
                <TableHeaderCell>Name</TableHeaderCell>
                <TableHeaderCell>Public ID</TableHeaderCell>
                <TableHeaderCell>Created</TableHeaderCell>
                <TableHeaderCell></TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data.projects.map((project) => (
                <TableRow key={project.publicId}>
                  <TableCell className="font-medium text-gray-900">
                    <Link to={`/projects/${project.publicId}/campaigns`} className="hover:text-primary-600">
                      {project.name}
                    </Link>
                  </TableCell>
                  <TableCell className="font-mono text-xs text-gray-500">{project.publicId}</TableCell>
                  <TableCell>{new Date(project.createdAt).toLocaleDateString()}</TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-3 text-sm">
                      <button
                        className="text-gray-500 hover:text-gray-900"
                        onClick={() => setDialogState({ open: true, project })}
                      >
                        Rename
                      </button>
                      <button className="text-red-500 hover:text-red-700" onClick={() => handleDelete(project)}>
                        Delete
                      </button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </Card>

      <ProjectFormDialog
        open={dialogState.open}
        project={dialogState.project}
        onClose={() => setDialogState({ open: false, project: null })}
      />
    </div>
  );
}
