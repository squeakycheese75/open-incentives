import { describe, expect, it } from "vitest";
import { mapEvaluateResponse } from "./mapper";
import type { EvaluateApiResponse } from "./types";

describe("mapEvaluateResponse", () => {
  it("maps a response with no discounts", () => {
    const response: EvaluateApiResponse = {
      decision: { matched: false, campaignsMatched: 0 },
      cart: { subtotal: 30, discountTotal: 0, total: 30, currency: "EUR" },
      discounts: [],
    };

    const result = mapEvaluateResponse(response);

    expect(result).toEqual({
      subtotal: 30,
      discountTotal: 0,
      total: 30,
      currency: "EUR",
      matchedCampaigns: [],
      discounts: [],
    });
  });

  it("maps a response with multiple discounts and dedupes matched campaigns", () => {
    const response: EvaluateApiResponse = {
      decision: { matched: true, campaignsMatched: 1 },
      cart: { subtotal: 50, discountTotal: 8, total: 42, currency: "EUR" },
      discounts: [
        {
          campaignId: "cmp_summer",
          campaignName: "10% off orders over €50",
          ruleId: "rule_1",
          type: "percentage_discount",
          amount: 5,
        },
        {
          campaignId: "cmp_summer",
          campaignName: "10% off orders over €50",
          ruleId: "rule_2",
          type: "fixed_discount",
          amount: 3,
        },
      ],
    };

    const result = mapEvaluateResponse(response);

    expect(result.total).toBe(42);
    expect(result.discounts).toHaveLength(2);
    expect(result.matchedCampaigns).toEqual([
      { campaignId: "cmp_summer", campaignName: "10% off orders over €50" },
    ]);
  });
});
