"use client";

import { useEffect, useRef, useState } from "react";
import type { CartItem, Customer } from "@/context/CartContext";
import type { NormalizedEvaluateResult } from "@/lib/incentives/types";
import { evaluateCart } from "@/lib/evaluate";

const DEBOUNCE_MS = 300;

export function useDebouncedEvaluation(
  items: CartItem[],
  customer: Customer,
  currency: string,
) {
  const [result, setResult] = useState<NormalizedEvaluateResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const requestId = useRef(0);

  const itemsKey = JSON.stringify(items);
  const customerKey = JSON.stringify(customer);

  useEffect(() => {
    if (items.length === 0) {
      setResult(null);
      setError(null);
      setLoading(false);
      return;
    }

    const thisRequestId = ++requestId.current;
    const controller = new AbortController();
    setLoading(true);
    setError(null);

    const timer = setTimeout(async () => {
      try {
        const evaluation = await evaluateCart(items, customer, currency, controller.signal);
        if (thisRequestId === requestId.current) {
          setResult(evaluation);
          setError(null);
        }
      } catch (err) {
        if (thisRequestId === requestId.current && !controller.signal.aborted) {
          setError(err instanceof Error ? err.message : "The cart could not be evaluated.");
        }
      } finally {
        if (thisRequestId === requestId.current) {
          setLoading(false);
        }
      }
    }, DEBOUNCE_MS);

    return () => {
      clearTimeout(timer);
      controller.abort();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [itemsKey, customerKey, currency]);

  return { result, loading, error };
}
