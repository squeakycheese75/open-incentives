"use client";

import type { Product } from "@/data/products";
import { formatCurrency } from "@/lib/currency";
import { useCart } from "@/context/CartContext";

export function ProductCard({ product }: { product: Product }) {
  const { addItem } = useCart();

  return (
    <div className="product-card">
      <h3>{product.name}</h3>
      <p>{product.description}</p>
      <span className="product-price">
        {formatCurrency(product.price, product.currency)}
      </span>
      <button
        type="button"
        className="button-primary"
        onClick={() => addItem(product)}
      >
        Add to Cart
      </button>
    </div>
  );
}
