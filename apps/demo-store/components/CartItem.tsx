"use client";

import type { CartItem as CartItemType } from "@/context/CartContext";
import { formatCurrency } from "@/lib/currency";

export function CartItem({
  item,
  currency,
  onQuantityChange,
  onRemove,
}: {
  item: CartItemType;
  currency: string;
  onQuantityChange: (quantity: number) => void;
  onRemove: () => void;
}) {
  return (
    <div className="cart-item-row">
      <span className="cart-item-name">{item.name}</span>
      <span>{formatCurrency(item.unitPrice, currency)}</span>
      <div className="quantity-controls">
        <button
          type="button"
          className="button-secondary"
          aria-label={`Decrease quantity of ${item.name}`}
          onClick={() => onQuantityChange(item.quantity - 1)}
        >
          -
        </button>
        <span>{item.quantity}</span>
        <button
          type="button"
          className="button-secondary"
          aria-label={`Increase quantity of ${item.name}`}
          onClick={() => onQuantityChange(item.quantity + 1)}
        >
          +
        </button>
      </div>
      <span>{formatCurrency(item.unitPrice * item.quantity, currency)}</span>
      <button type="button" className="button-secondary" onClick={onRemove}>
        Remove
      </button>
    </div>
  );
}
