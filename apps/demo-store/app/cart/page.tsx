"use client";

import Link from "next/link";
import { useCart } from "@/context/CartContext";
import { useDebouncedEvaluation } from "@/lib/useDebouncedEvaluation";
import { CartItem } from "@/components/CartItem";
import { CartSummary } from "@/components/CartSummary";
import { CustomerSelector } from "@/components/CustomerSelector";

export default function CartPage() {
  const { items, currency, customer, subtotal, removeItem, setQuantity, setCustomer } =
    useCart();
  const { result, loading, error } = useDebouncedEvaluation(items, customer, currency);

  if (items.length === 0) {
    return (
      <div>
        <h2>Your cart is empty</h2>
        <p>
          <Link href="/">Browse products</Link> to get started.
        </p>
      </div>
    );
  }

  return (
    <div>
      <h2>Your Cart</h2>

      <div>
        {items.map((item) => (
          <CartItem
            key={item.productId}
            item={item}
            currency={currency}
            onQuantityChange={(quantity) => setQuantity(item.productId, quantity)}
            onRemove={() => removeItem(item.productId)}
          />
        ))}
      </div>

      <CustomerSelector customer={customer} onChange={setCustomer} />

      <CartSummary
        localSubtotal={subtotal}
        currency={currency}
        evaluation={result}
        loading={loading}
        error={error}
      />

      <p>
        <Link href="/checkout" className="button-primary">
          Continue to Checkout
        </Link>
      </p>
    </div>
  );
}
