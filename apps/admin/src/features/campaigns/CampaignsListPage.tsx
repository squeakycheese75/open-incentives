import { Link, useParams } from "react-router-dom";

import { Badge } from "../../components/ui/Badge";
import { Card } from "../../components/ui/Card";
import { Spinner } from "../../components/ui/Spinner";
import { Table, TableBody, TableCell, TableHead, TableHeaderCell, TableRow } from "../../components/ui/Table";
import { useToast } from "../../components/ui/ToastProvider";
import { ApiError } from "../../api/client";
import type { Campaign } from "../../api/types";
import { isRuleList } from "./ruleTypes";
import { RuleSummary } from "./RuleSummary";
import { useCampaigns, useDeleteCampaign } from "./api";

export function CampaignsListPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const { data, isLoading, isError } = useCampaigns(projectId!);
  const deleteCampaign = useDeleteCampaign(projectId!);
  const { showToast } = useToast();

  async function handleDelete(campaign: Campaign) {
    if (!window.confirm(`Delete campaign "${campaign.name}"? This cannot be undone.`)) return;

    try {
      await deleteCampaign.mutateAsync(campaign.publicId);
      showToast("Campaign deleted");
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : "Failed to delete campaign", "error");
    }
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-gray-900">Campaigns</h1>
        <Link
          to={`/projects/${projectId}/campaigns/new`}
          className="inline-flex items-center justify-center gap-1.5 rounded bg-primary-600 px-3.5 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
        >
          New campaign
        </Link>
      </div>

      <Card>
        {isLoading && (
          <div className="flex justify-center py-10">
            <Spinner />
          </div>
        )}

        {isError && <p className="p-6 text-sm text-red-600">Failed to load campaigns.</p>}

        {data && data.campaigns.length === 0 && (
          <p className="p-6 text-sm text-gray-500">No campaigns yet. Create one to start applying discounts.</p>
        )}

        {data && data.campaigns.length > 0 && (
          <Table>
            <TableHead>
              <TableRow>
                <TableHeaderCell>Name</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
                <TableHeaderCell>Rules</TableHeaderCell>
                <TableHeaderCell></TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data.campaigns.map((campaign) => (
                <TableRow key={campaign.publicId}>
                  <TableCell className="font-medium text-gray-900">{campaign.name}</TableCell>
                  <TableCell>
                    <Badge tone={campaign.status === "active" ? "green" : "gray"}>{campaign.status}</Badge>
                  </TableCell>
                  <TableCell className="max-w-xs">
                    {isRuleList(campaign.rules) && <RuleSummary rules={campaign.rules} />}
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-3 text-sm">
                      <Link
                        to={`/projects/${projectId}/campaigns/${campaign.publicId}/edit`}
                        className="text-gray-500 hover:text-gray-900"
                      >
                        Edit
                      </Link>
                      <button className="text-red-500 hover:text-red-700" onClick={() => handleDelete(campaign)}>
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
    </div>
  );
}
