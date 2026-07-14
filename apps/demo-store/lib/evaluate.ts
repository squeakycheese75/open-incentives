import type { CartItem, Customer } from "@/context/CartContext";
import type { NormalizedEvaluateResult } from "@/lib/incentives/types";

export class EvaluationRequestError extends Error {}

export async function evaluateCart(
  items: CartItem[],
  customer: Customer,
  currency: string,
  signal?: AbortSignal,
): Promise<NormalizedEvaluateResult> {
  const response = await fetch("/api/evaluate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      customer,
      cart: {
        currency,
        items: items.map((item) => ({
          productId: item.productId,
          quantity: item.quantity,
          unitPrice: item.unitPrice,
        })),
      },
    }),
    signal,
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as { message?: string };
    throw new EvaluationRequestError(body.message ?? "The cart could not be evaluated.");
  }

  return (await response.json()) as NormalizedEvaluateResult;
}
