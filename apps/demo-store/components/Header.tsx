"use client";

import Link from "next/link";
import { useCart } from "@/context/CartContext";

const STORE_NAME = process.env.NEXT_PUBLIC_STORE_NAME ?? "Open Incentives Demo Store";

export function Header() {
  const { itemCount } = useCart();

  return (
    <header className="header">
      <div className="header-inner">
        <div className="header-title-group">
          <h1>{STORE_NAME}</h1>
          <p className="header-subtitle">
            Promotions in this store are evaluated in real time by Open Incentives.
          </p>
        </div>
        <nav className="header-nav">
          <Link href="/">Store</Link>
          <Link href="/cart">
            Cart <span className="cart-count">{itemCount}</span>
          </Link>
        </nav>
      </div>
    </header>
  );
}
