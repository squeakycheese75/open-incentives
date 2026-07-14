"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { products, type Product } from "@/data/products";

export type CartItem = {
  productId: string;
  name: string;
  unitPrice: number;
  quantity: number;
};

export type Customer = {
  id: string;
  country: string;
  tier: string;
};

const CART_STORAGE_KEY = "open-incentives-demo-cart";
const CUSTOMER_STORAGE_KEY = "open-incentives-demo-customer";

const DEFAULT_CUSTOMER: Customer = {
  id: "user_123",
  country: "DE",
  tier: "gold",
};

type CartContextValue = {
  items: CartItem[];
  currency: string;
  customer: Customer;
  subtotal: number;
  itemCount: number;
  addItem: (product: Product) => void;
  removeItem: (productId: string) => void;
  setQuantity: (productId: string, quantity: number) => void;
  setCustomer: (customer: Customer) => void;
};

const CartContext = createContext<CartContextValue | undefined>(undefined);

function readStorage<T>(key: string, fallback: T): T {
  if (typeof window === "undefined") {
    return fallback;
  }

  try {
    const raw = window.localStorage.getItem(key);
    if (!raw) {
      return fallback;
    }
    return JSON.parse(raw) as T;
  } catch {
    return fallback;
  }
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>(() => readStorage(CART_STORAGE_KEY, []));
  const [customer, setCustomerState] = useState<Customer>(() =>
    readStorage(CUSTOMER_STORAGE_KEY, DEFAULT_CUSTOMER),
  );
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    setHydrated(true);
  }, []);

  useEffect(() => {
    if (!hydrated) return;
    window.localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(items));
  }, [items, hydrated]);

  useEffect(() => {
    if (!hydrated) return;
    window.localStorage.setItem(CUSTOMER_STORAGE_KEY, JSON.stringify(customer));
  }, [customer, hydrated]);

  const addItem = useCallback((product: Product) => {
    setItems((current) => {
      const existing = current.find((item) => item.productId === product.id);
      if (existing) {
        return current.map((item) =>
          item.productId === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        );
      }

      return [
        ...current,
        {
          productId: product.id,
          name: product.name,
          unitPrice: product.price,
          quantity: 1,
        },
      ];
    });
  }, []);

  const removeItem = useCallback((productId: string) => {
    setItems((current) => current.filter((item) => item.productId !== productId));
  }, []);

  const setQuantity = useCallback((productId: string, quantity: number) => {
    if (quantity <= 0) {
      setItems((current) => current.filter((item) => item.productId !== productId));
      return;
    }

    setItems((current) =>
      current.map((item) => (item.productId === productId ? { ...item, quantity } : item)),
    );
  }, []);

  const setCustomer = useCallback((next: Customer) => {
    setCustomerState(next);
  }, []);

  const currency = products[0]?.currency ?? "EUR";

  const subtotal = useMemo(
    () => items.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0),
    [items],
  );

  const itemCount = useMemo(
    () => items.reduce((count, item) => count + item.quantity, 0),
    [items],
  );

  const value: CartContextValue = {
    items,
    currency,
    customer,
    subtotal,
    itemCount,
    addItem,
    removeItem,
    setQuantity,
    setCustomer,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return ctx;
}
