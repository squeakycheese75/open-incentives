import { NextRequest, NextResponse } from "next/server";
import { createIncentivesClient, IncentivesApiError } from "@/lib/incentives/client";
import { mapEvaluateResponse } from "@/lib/incentives/mapper";
import type { EvaluateApiResponse } from "@/lib/incentives/types";

export type BrowserCartItem = {
  productId: string;
  quantity: number;
  unitPrice: number;
};

export type BrowserEvaluateRequest = {
  customer: {
    id: string;
    country: string;
    tier: string;
  };
  cart: {
    items: BrowserCartItem[];
    currency: string;
  };
};

const UNAVAILABLE_MESSAGE =
  "Promotions are temporarily unavailable. Your cart has not been changed.";
const UNAUTHORIZED_MESSAGE = "The promotion service could not be authenticated.";
const INVALID_REQUEST_MESSAGE = "The cart could not be evaluated.";

function isValidRequest(body: unknown): body is BrowserEvaluateRequest {
  if (!body || typeof body !== "object") {
    return false;
  }

  const req = body as Partial<BrowserEvaluateRequest>;

  if (!req.customer || typeof req.customer.id !== "string" || req.customer.id === "") {
    return false;
  }

  if (!req.cart || typeof req.cart.currency !== "string" || req.cart.currency === "") {
    return false;
  }

  if (!Array.isArray(req.cart.items)) {
    return false;
  }

  return req.cart.items.every(
    (item) =>
      item &&
      typeof item.productId === "string" &&
      item.productId !== "" &&
      typeof item.quantity === "number" &&
      item.quantity > 0 &&
      typeof item.unitPrice === "number" &&
      item.unitPrice >= 0,
  );
}

export async function POST(request: NextRequest) {
  let body: unknown;

  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ message: INVALID_REQUEST_MESSAGE }, { status: 400 });
  }

  if (!isValidRequest(body)) {
    return NextResponse.json({ message: INVALID_REQUEST_MESSAGE }, { status: 400 });
  }

  let client;
  try {
    client = createIncentivesClient();
  } catch {
    return NextResponse.json({ message: UNAVAILABLE_MESSAGE }, { status: 503 });
  }

  let apiResponse: EvaluateApiResponse;
  try {
    apiResponse = await client.evaluate({
      customer: {
        id: body.customer.id,
        country: body.customer.country,
        tier: body.customer.tier,
      },
      cart: {
        currency: body.cart.currency,
        items: body.cart.items,
      },
    });
  } catch (err) {
    if (err instanceof IncentivesApiError) {
      if (err.status === "unauthorized") {
        console.error("incentives api authentication failed");
        return NextResponse.json({ message: UNAUTHORIZED_MESSAGE }, { status: 502 });
      }

      if (err.status === "invalid_request") {
        return NextResponse.json({ message: INVALID_REQUEST_MESSAGE }, { status: 400 });
      }

      return NextResponse.json({ message: UNAVAILABLE_MESSAGE }, { status: 503 });
    }

    console.error("unexpected error calling incentives api", err);
    return NextResponse.json({ message: UNAVAILABLE_MESSAGE }, { status: 503 });
  }

  return NextResponse.json(mapEvaluateResponse(apiResponse));
}
