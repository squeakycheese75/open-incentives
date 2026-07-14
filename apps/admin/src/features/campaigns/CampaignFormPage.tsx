import { useEffect, useState, type FormEvent } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { Button } from "../../components/ui/Button";
import { Card } from "../../components/ui/Card";
import { Input } from "../../components/ui/Input";
import { Select } from "../../components/ui/Select";
import { Spinner } from "../../components/ui/Spinner";
import { useToast } from "../../components/ui/ToastProvider";
import { ApiError } from "../../api/client";
import type { CampaignStatus } from "../../api/types";
import { RuleEditor } from "./RuleEditor";
import { useCampaign, useCreateCampaign, useUpdateCampaign } from "./api";

const DEFAULT_RULES = JSON.stringify(
  [
    {
      id: "rule_10_percent_over_50",
      name: "10% off over $50",
      conditions: {
        all: [{ fact: "cart.total", operator: "gte", value: 50 }],
      },
      actions: [{ type: "percentage_discount", params: { value: 10 } }],
    },
  ],
  null,
  2,
);

interface CampaignFormPageProps {
  mode: "create" | "edit";
}

export function CampaignFormPage({ mode }: CampaignFormPageProps) {
  const { projectId, campaignId } = useParams<{ projectId: string; campaignId: string }>();
  const navigate = useNavigate();
  const { showToast } = useToast();

  const isEdit = mode === "edit";
  const existing = useCampaign(projectId!, campaignId!);
  const createCampaign = useCreateCampaign(projectId!);
  const updateCampaign = useUpdateCampaign(projectId!, campaignId!);

  const [name, setName] = useState("");
  const [status, setStatus] = useState<CampaignStatus>("active");
  const [rulesText, setRulesText] = useState(DEFAULT_RULES);
  const [rulesError, setRulesError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    if (isEdit && existing.data) {
      setName(existing.data.name);
      setStatus(existing.data.status);
      setRulesText(JSON.stringify(existing.data.rules, null, 2));
    }
  }, [isEdit, existing.data]);

  function handleRulesChange(next: string) {
    setRulesText(next);
    try {
      JSON.parse(next);
      setRulesError(null);
    } catch {
      setRulesError("Rules must be valid JSON");
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setFormError(null);

    let rules: unknown;
    try {
      rules = JSON.parse(rulesText);
    } catch {
      setRulesError("Rules must be valid JSON");
      return;
    }

    try {
      if (isEdit) {
        await updateCampaign.mutateAsync({ name, status, rules });
        showToast("Campaign updated");
      } else {
        await createCampaign.mutateAsync({ name, status, rules });
        showToast("Campaign created");
      }
      navigate(`/projects/${projectId}/campaigns`);
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong");
    }
  }

  if (isEdit && existing.isLoading) {
    return (
      <div className="flex justify-center py-10">
        <Spinner />
      </div>
    );
  }

  const isPending = createCampaign.isPending || updateCampaign.isPending;

  return (
    <div>
      <h1 className="mb-6 text-lg font-semibold text-gray-900">{isEdit ? "Edit campaign" : "New campaign"}</h1>

      <Card className="p-6">
        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label htmlFor="campaign-name" className="mb-1 block text-sm font-medium text-gray-700">
              Name
            </label>
            <Input id="campaign-name" required value={name} onChange={(e) => setName(e.target.value)} />
          </div>

          <div>
            <label htmlFor="campaign-status" className="mb-1 block text-sm font-medium text-gray-700">
              Status
            </label>
            <Select
              id="campaign-status"
              value={status}
              onChange={(e) => setStatus(e.target.value as CampaignStatus)}
              className="w-40"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </Select>
          </div>

          <RuleEditor value={rulesText} onChange={handleRulesChange} error={rulesError} />

          {formError && <p className="text-sm text-red-600">{formError}</p>}

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="secondary" onClick={() => navigate(`/projects/${projectId}/campaigns`)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending || !!rulesError}>
              {isPending ? "Saving..." : "Save"}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
