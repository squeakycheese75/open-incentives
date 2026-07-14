import type { NormalizedEvaluateResult } from "@/lib/incentives/types";
import { formatCurrency } from "@/lib/currency";
import { DiscountList } from "./DiscountList";

export function CartSummary({
  localSubtotal,
  currency,
  evaluation,
  loading,
  error,
}: {
  localSubtotal: number;
  currency: string;
  evaluation: NormalizedEvaluateResult | null;
  loading: boolean;
  error: string | null;
}) {
  const subtotal = evaluation?.subtotal ?? localSubtotal;
  const total = evaluation?.total ?? localSubtotal;

  return (
    <div className="summary-card">
      <div className="summary-row">
        <span>Subtotal</span>
        <span>{formatCurrency(subtotal, currency)}</span>
      </div>

      {evaluation && <DiscountList discounts={evaluation.discounts} currency={currency} />}

      <div className="summary-row total">
        <span>Total</span>
        <span>{formatCurrency(total, currency)}</span>
      </div>

      {loading && <p className="status-text">Checking available promotions…</p>}
      {error && <p className="error-text">{error}</p>}
    </div>
  );
}
