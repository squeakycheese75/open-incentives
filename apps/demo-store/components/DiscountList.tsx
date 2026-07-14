import type { AppliedDiscount } from "@/lib/incentives/types";
import { formatCurrency } from "@/lib/currency";

export function DiscountList({
  discounts,
  currency,
}: {
  discounts: AppliedDiscount[];
  currency: string;
}) {
  if (discounts.length === 0) {
    return <p className="status-text">No promotions currently apply to this cart.</p>;
  }

  return (
    <div>
      {discounts.map((discount) => (
        <div className="summary-row discount-row" key={`${discount.campaignId}-${discount.ruleId}`}>
          <span>{discount.campaignName}</span>
          <span>-{formatCurrency(discount.amount, currency)}</span>
        </div>
      ))}
      <p className="status-text">
        {discounts.length} {discounts.length === 1 ? "promotion" : "promotions"} applied
      </p>
    </div>
  );
}
