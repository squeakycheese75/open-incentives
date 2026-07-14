import type {
  AppliedDiscount,
  EvaluateApiResponse,
  MatchedCampaign,
  NormalizedEvaluateResult,
} from "./types";

export function mapEvaluateResponse(
  response: EvaluateApiResponse,
): NormalizedEvaluateResult {
  const discounts: AppliedDiscount[] = response.discounts.map((discount) => ({
    campaignId: discount.campaignId,
    campaignName: discount.campaignName,
    ruleId: discount.ruleId,
    type: discount.type,
    amount: discount.amount,
  }));

  const matchedCampaigns: MatchedCampaign[] = [];
  const seen = new Set<string>();

  for (const discount of discounts) {
    if (seen.has(discount.campaignId)) {
      continue;
    }
    seen.add(discount.campaignId);
    matchedCampaigns.push({
      campaignId: discount.campaignId,
      campaignName: discount.campaignName,
    });
  }

  return {
    subtotal: response.cart.subtotal,
    discountTotal: response.cart.discountTotal,
    total: response.cart.total,
    currency: response.cart.currency,
    matchedCampaigns,
    discounts,
  };
}
