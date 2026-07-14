import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";
import { POST } from "./route";

function makeRequest(body: unknown) {
  return new NextRequest("http://localhost/api/evaluate", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

const validBody = {
  customer: { id: "user_123", country: "DE", tier: "gold" },
  cart: {
    currency: "EUR",
    items: [{ productId: "prod_coffee", quantity: 2, unitPrice: 18 }],
  },
};

describe("POST /api/evaluate", () => {
  const originalEnv = { ...process.env };

  beforeEach(() => {
    process.env.INCENTIVES_API_URL = "http://localhost:8080";
    process.env.INCENTIVES_API_KEY = "api_test.secret";
  });

  afterEach(() => {
    process.env = { ...originalEnv };
    vi.unstubAllGlobals();
  });

  it("returns a normalized result on success", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            decision: { matched: true, campaignsMatched: 1 },
            cart: { subtotal: 36, discountTotal: 0, total: 36, currency: "EUR" },
            discounts: [],
          }),
          { status: 200 },
        ),
      ),
    );

    const response = await POST(makeRequest(validBody));
    const json = await response.json();

    expect(response.status).toBe(200);
    expect(json.total).toBe(36);
    expect(json.matchedCampaigns).toEqual([]);
  });

  it("returns 503 when environment variables are missing", async () => {
    delete process.env.INCENTIVES_API_URL;
    delete process.env.INCENTIVES_API_KEY;

    const response = await POST(makeRequest(validBody));
    expect(response.status).toBe(503);
  });

  it("returns 502 when the incentives API returns 401", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(new Response("", { status: 401 })),
    );

    const response = await POST(makeRequest(validBody));
    expect(response.status).toBe(502);
  });

  it("returns 503 when the incentives API returns 500", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(new Response("", { status: 500 })),
    );

    const response = await POST(makeRequest(validBody));
    expect(response.status).toBe(503);
  });

  it("returns 503 when the incentives API times out", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockImplementation(() => {
        const err = new Error("aborted");
        err.name = "AbortError";
        return Promise.reject(err);
      }),
    );

    const response = await POST(makeRequest(validBody));
    expect(response.status).toBe(503);
  });

  it("returns 400 for a malformed browser request", async () => {
    const response = await POST(makeRequest({ customer: {}, cart: {} }));
    expect(response.status).toBe(400);
  });
});
