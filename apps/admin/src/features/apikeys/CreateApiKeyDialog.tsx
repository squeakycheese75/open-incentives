import { useEffect, useState, type FormEvent } from "react";

import { Dialog } from "../../components/ui/Dialog";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { ApiError } from "../../api/client";
import { useCreateApiKey } from "./api";

interface CreateApiKeyDialogProps {
  open: boolean;
  onClose: () => void;
  projectId: string;
}

export function CreateApiKeyDialog({ open, onClose, projectId }: CreateApiKeyDialogProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [createdKey, setCreatedKey] = useState<string | null>(null);

  const createApiKey = useCreateApiKey(projectId);

  useEffect(() => {
    if (open) {
      setName("");
      setDescription("");
      setError(null);
      setCreatedKey(null);
    }
  }, [open]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    try {
      const res = await createApiKey.mutateAsync({ name, description });
      setCreatedKey(res.apiKey);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to create API key");
    }
  }

  function handleClose() {
    setCreatedKey(null);
    onClose();
  }

  return (
    <Dialog open={open} onClose={handleClose} title="Create API key">
      {createdKey ? (
        <div className="space-y-4">
          <p className="text-sm text-gray-700">
            Copy this key now — for security, it will not be shown again.
          </p>
          <div className="flex items-center gap-2">
            <code className="flex-1 truncate rounded bg-gray-100 px-3 py-2 text-xs">{createdKey}</code>
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigator.clipboard.writeText(createdKey)}
            >
              Copy
            </Button>
          </div>
          <div className="flex justify-end pt-2">
            <Button type="button" onClick={handleClose}>
              Done
            </Button>
          </div>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="key-name" className="mb-1 block text-sm font-medium text-gray-700">
              Name
            </label>
            <Input id="key-name" required value={name} onChange={(e) => setName(e.target.value)} autoFocus />
          </div>

          <div>
            <label htmlFor="key-description" className="mb-1 block text-sm font-medium text-gray-700">
              Description
            </label>
            <Input id="key-description" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>

          {error && <p className="text-sm text-red-600">{error}</p>}

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="secondary" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={createApiKey.isPending}>
              {createApiKey.isPending ? "Creating..." : "Create"}
            </Button>
          </div>
        </form>
      )}
    </Dialog>
  );
}
