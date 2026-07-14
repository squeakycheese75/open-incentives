"use client";

import { useEffect, useRef, useState } from "react";
import type { Product } from "@/data/products";
import { formatCurrency } from "@/lib/currency";
import { useCart } from "@/context/CartContext";

export function ProductCard({ product }: { product: Product }) {
  const { addItem } = useCart();
  const [added, setAdded] = useState(false);
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => () => clearTimeout(timeoutRef.current), []);

  const handleClick = () => {
    addItem(product);
    setAdded(true);
    clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => setAdded(false), 1200);
  };

  return (
    <div className={`product-card${added ? " product-card-added" : ""}`}>
      <h3>{product.name}</h3>
      <p>{product.description}</p>
      <span className="product-price">
        {formatCurrency(product.price, product.currency)}
      </span>
      <button
        type="button"
        className={`button-primary${added ? " button-added" : ""}`}
        onClick={handleClick}
      >
        {added ? "Added ✓" : "Add to Cart"}
      </button>
    </div>
  );
}
