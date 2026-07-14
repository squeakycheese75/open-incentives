"use client";

import type { Customer } from "@/context/CartContext";

const COUNTRIES = ["DE", "FR", "GB", "US"];
const TIERS = ["standard", "silver", "gold"];

export function CustomerSelector({
  customer,
  onChange,
}: {
  customer: Customer;
  onChange: (customer: Customer) => void;
}) {
  return (
    <div className="customer-selector">
      <label>
        Customer ID
        <input
          type="text"
          value={customer.id}
          onChange={(e) => onChange({ ...customer, id: e.target.value })}
        />
      </label>
      <label>
        Country
        <select
          value={customer.country}
          onChange={(e) => onChange({ ...customer, country: e.target.value })}
        >
          {COUNTRIES.map((country) => (
            <option key={country} value={country}>
              {country}
            </option>
          ))}
        </select>
      </label>
      <label>
        Tier
        <select
          value={customer.tier}
          onChange={(e) => onChange({ ...customer, tier: e.target.value })}
        >
          {TIERS.map((tier) => (
            <option key={tier} value={tier}>
              {tier}
            </option>
          ))}
        </select>
      </label>
    </div>
  );
}
