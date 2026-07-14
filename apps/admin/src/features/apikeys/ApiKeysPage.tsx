import { useState } from "react";
import { useParams } from "react-router-dom";

import { Badge } from "../../components/ui/Badge";
import { Button } from "../../components/ui/Button";
import { Card } from "../../components/ui/Card";
import { Spinner } from "../../components/ui/Spinner";
import { Table, TableBody, TableCell, TableHead, TableHeaderCell, TableRow } from "../../components/ui/Table";
import { useToast } from "../../components/ui/ToastProvider";
import { ApiError } from "../../api/client";
import type { ApiKey } from "../../api/types";
import { useApiKeys, useRevokeApiKey } from "./api";
import { CreateApiKeyDialog } from "./CreateApiKeyDialog";

export function ApiKeysPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const { data, isLoading, isError } = useApiKeys(projectId!);
  const revokeApiKey = useRevokeApiKey(projectId!);
  const { showToast } = useToast();

  const [dialogOpen, setDialogOpen] = useState(false);

  async function handleRevoke(key: ApiKey) {
    if (!window.confirm(`Revoke API key "${key.name}"? This cannot be undone.`)) return;

    try {
      await revokeApiKey.mutateAsync(key.publicId);
      showToast("API key revoked");
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : "Failed to revoke API key", "error");
    }
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-gray-900">API Keys</h1>
        <Button onClick={() => setDialogOpen(true)}>New API key</Button>
      </div>

      <Card>
        {isLoading && (
          <div className="flex justify-center py-10">
            <Spinner />
          </div>
        )}

        {isError && <p className="p-6 text-sm text-red-600">Failed to load API keys.</p>}

        {data && data.apiKeys.length === 0 && (
          <p className="p-6 text-sm text-gray-500">No API keys yet. Create one to call the evaluate endpoint.</p>
        )}

        {data && data.apiKeys.length > 0 && (
          <Table>
            <TableHead>
              <TableRow>
                <TableHeaderCell>Name</TableHeaderCell>
                <TableHeaderCell>Prefix</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
                <TableHeaderCell>Created</TableHeaderCell>
                <TableHeaderCell></TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data.apiKeys.map((key) => (
                <TableRow key={key.publicId}>
                  <TableCell className="font-medium text-gray-900">{key.name}</TableCell>
                  <TableCell className="font-mono text-xs text-gray-500">{key.prefix}</TableCell>
                  <TableCell>
                    <Badge tone={key.status === "active" ? "green" : "red"}>{key.status}</Badge>
                  </TableCell>
                  <TableCell>{new Date(key.createdAt).toLocaleDateString()}</TableCell>
                  <TableCell>
                    {key.status === "active" && (
                      <div className="flex justify-end text-sm">
                        <button className="text-red-500 hover:text-red-700" onClick={() => handleRevoke(key)}>
                          Revoke
                        </button>
                      </div>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </Card>

      <CreateApiKeyDialog open={dialogOpen} onClose={() => setDialogOpen(false)} projectId={projectId!} />
    </div>
  );
}
