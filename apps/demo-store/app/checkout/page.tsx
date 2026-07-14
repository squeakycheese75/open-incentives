"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useCart } from "@/context/CartContext";
import { evaluateCart, EvaluationRequestError } from "@/lib/evaluate";
import type { NormalizedEvaluateResult } from "@/lib/incentives/types";
import { CartSummary } from "@/components/CartSummary";

export default function CheckoutPage() {
  const { items, currency, customer, subtotal } = useCart();
  const [evaluation, setEvaluation] = useState<NormalizedEvaluateResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [orderCompleted, setOrderCompleted] = useState(false);

  useEffect(() => {
    if (items.length === 0) return;

    let cancelled = false;
    setLoading(true);
    setError(null);

    evaluateCart(items, customer, currency)
      .then((result) => {
        if (!cancelled) setEvaluation(result);
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof EvaluationRequestError
              ? err.message
              : "The cart could not be evaluated.",
          );
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handlePlaceOrder() {
    setLoading(true);
    setError(null);

    try {
      const result = await evaluateCart(items, customer, currency);
      setEvaluation(result);
      setOrderCompleted(true);
    } catch (err) {
      setError(
        err instanceof EvaluationRequestError ? err.message : "The cart could not be evaluated.",
      );
    } finally {
      setLoading(false);
    }
  }

  if (items.length === 0 && !orderCompleted) {
    return (
      <div>
        <h2>Your cart is empty</h2>
        <p>
          <Link href="/">Browse products</Link> to get started.
        </p>
      </div>
    );
  }

  if (orderCompleted) {
    return (
      <div>
        <h2>Demo order completed</h2>
        <p>
          <Link href="/">Return to store</Link>
        </p>
      </div>
    );
  }

  return (
    <div>
      <h2>Checkout Summary</h2>

      <div>
        {items.map((item) => (
          <div className="cart-item-row" key={item.productId}>
            <span className="cart-item-name">
              {item.name} × {item.quantity}
            </span>
          </div>
        ))}
      </div>

      <p className="status-text">
        Customer: {customer.id} · {customer.country} · {customer.tier}
      </p>

      <CartSummary
        localSubtotal={subtotal}
        currency={currency}
        evaluation={evaluation}
        loading={loading}
        error={error}
      />

      <p>
        <button
          type="button"
          className="button-primary"
          disabled={loading}
          onClick={handlePlaceOrder}
        >
          Place Demo Order
        </button>
      </p>
    </div>
  );
}
