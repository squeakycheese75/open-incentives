import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";
import StorefrontPage from "./page";
import { CartProvider, useCart } from "@/context/CartContext";

function CartCountProbe() {
  const { itemCount } = useCart();
  return <span data-testid="item-count">{itemCount}</span>;
}

function renderPage() {
  return render(
    <CartProvider>
      <CartCountProbe />
      <StorefrontPage />
    </CartProvider>,
  );
}

describe("StorefrontPage", () => {
  it("renders at least four products", () => {
    renderPage();

    expect(screen.getByText("Premium Coffee")).toBeInTheDocument();
    expect(screen.getByText("Ceramic Mug")).toBeInTheDocument();
    expect(screen.getByText("Coffee Grinder")).toBeInTheDocument();
    expect(screen.getByText("Coffee Filters")).toBeInTheDocument();
  });

  it("adds a product to the cart", async () => {
    renderPage();
    const user = userEvent.setup();

    await user.click(screen.getAllByText("Add to Cart")[0]);

    expect(screen.getByTestId("item-count").textContent).toBe("1");
  });
});
