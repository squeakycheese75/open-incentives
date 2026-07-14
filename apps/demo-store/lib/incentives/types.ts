// Request/response shapes for the Open Incentives /v1/evaluate API.
// Note: the real API is fully camelCase and returns fields nested under
// `decision`/`cart`/`discounts` — it does not match the flatter shape
// described in early integration drafts.

export type EvaluateCustomer = {
  id: string;
  country: string;
  tier: string;
};

export type EvaluateCartItem = {
  productId: string;
  quantity: number;
  unitPrice: number;
};

export type EvaluateRequest = {
  customer: EvaluateCustomer;
  cart: {
    currency: string;
    items: EvaluateCartItem[];
  };
};

export type EvaluateApiDecision = {
  matched: boolean;
  campaignsMatched: number;
};

export type EvaluateApiCart = {
  subtotal: number;
  discountTotal: number;
  total: number;
  currency: string;
};

export type EvaluateApiDiscount = {
  campaignId: string;
  campaignName: string;
  ruleId: string;
  type: string;
  amount: number;
};

export type EvaluateApiResponse = {
  decision: EvaluateApiDecision;
  cart: EvaluateApiCart;
  discounts: EvaluateApiDiscount[];
};

// Frontend-facing, normalized shape produced by lib/incentives/mapper.ts.
export type AppliedDiscount = {
  campaignId: string;
  campaignName: string;
  ruleId: string;
  type: string;
  amount: number;
};

export type MatchedCampaign = {
  campaignId: string;
  campaignName: string;
};

export type NormalizedEvaluateResult = {
  subtotal: number;
  discountTotal: number;
  total: number;
  currency: string;
  matchedCampaigns: MatchedCampaign[];
  discounts: AppliedDiscount[];
};
