import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it } from "vitest";
import { CartProvider, useCart } from "./CartContext";
import { products } from "@/data/products";

function setup() {
  return renderHook(() => useCart(), { wrapper: CartProvider });
}

describe("CartContext", () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  it("computes subtotal as sum(unitPrice * quantity)", () => {
    const { result } = setup();

    act(() => {
      result.current.addItem(products[0]); // coffee 18.00
      result.current.addItem(products[0]); // quantity -> 2
      result.current.addItem(products[1]); // mug 14.00
    });

    expect(result.current.subtotal).toBe(50);
    expect(result.current.itemCount).toBe(3);
  });

  it("updates quantity for an existing item", () => {
    const { result } = setup();

    act(() => {
      result.current.addItem(products[0]);
    });

    act(() => {
      result.current.setQuantity(products[0].id, 5);
    });

    expect(result.current.items[0].quantity).toBe(5);
    expect(result.current.subtotal).toBe(90);
  });

  it("removes an item when quantity is set to zero", () => {
    const { result } = setup();

    act(() => {
      result.current.addItem(products[0]);
      result.current.setQuantity(products[0].id, 0);
    });

    expect(result.current.items).toHaveLength(0);
  });

  it("removes an item explicitly", () => {
    const { result } = setup();

    act(() => {
      result.current.addItem(products[0]);
      result.current.removeItem(products[0].id);
    });

    expect(result.current.items).toHaveLength(0);
  });
});
