import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import CartPage from "./page";
import { CartProvider } from "@/context/CartContext";

function seedCart() {
  window.localStorage.setItem(
    "open-incentives-demo-cart",
    JSON.stringify([
      { productId: "prod_coffee", name: "Premium Coffee", unitPrice: 18, quantity: 2 },
      { productId: "prod_mug", name: "Ceramic Mug", unitPrice: 14, quantity: 1 },
    ]),
  );
  window.localStorage.setItem(
    "open-incentives-demo-customer",
    JSON.stringify({ id: "user_123", country: "DE", tier: "gold" }),
  );
}

function renderCartPage() {
  return render(
    <CartProvider>
      <CartPage />
    </CartProvider>,
  );
}

describe("CartPage", () => {
  beforeEach(() => {
    window.localStorage.clear();
    seedCart();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.useRealTimers();
  });

  it("shows a discount when the API returns one", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            subtotal: 50,
            discountTotal: 5,
            total: 45,
            currency: "EUR",
            matchedCampaigns: [{ campaignId: "cmp_summer", campaignName: "10% off" }],
            discounts: [
              {
                campaignId: "cmp_summer",
                campaignName: "10% off",
                ruleId: "rule_1",
                type: "percentage_discount",
                amount: 5,
              },
            ],
          }),
          { status: 200 },
        ),
      ),
    );

    renderCartPage();

    await waitFor(() => expect(screen.getByText("1 promotion applied")).toBeInTheDocument());
    expect(screen.getByText("10% off")).toBeInTheDocument();
  });

  it("shows the no-discount state when nothing matches", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            subtotal: 50,
            discountTotal: 0,
            total: 50,
            currency: "EUR",
            matchedCampaigns: [],
            discounts: [],
          }),
          { status: 200 },
        ),
      ),
    );

    renderCartPage();

    await waitFor(() =>
      expect(
        screen.getByText("No promotions currently apply to this cart."),
      ).toBeInTheDocument(),
    );
  });

  it("shows an error state when the API is unavailable", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(new Response(JSON.stringify({ message: "down" }), { status: 503 })),
    );

    renderCartPage();

    await waitFor(() =>
      expect(screen.getByText("down")).toBeInTheDocument(),
    );
  });

  it("removes an item from the cart", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            subtotal: 14,
            discountTotal: 0,
            total: 14,
            currency: "EUR",
            matchedCampaigns: [],
            discounts: [],
          }),
          { status: 200 },
        ),
      ),
    );

    renderCartPage();
    const user = userEvent.setup();

    expect(screen.getByText("Premium Coffee")).toBeInTheDocument();
    await user.click(screen.getAllByText("Remove")[0]);

    expect(screen.queryByText("Premium Coffee")).not.toBeInTheDocument();
  });

  it("changes item quantity", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            subtotal: 50,
            discountTotal: 0,
            total: 50,
            currency: "EUR",
            matchedCampaigns: [],
            discounts: [],
          }),
          { status: 200 },
        ),
      ),
    );

    renderCartPage();
    const user = userEvent.setup();

    await user.click(screen.getByLabelText("Increase quantity of Premium Coffee"));

    expect(await screen.findByText("3")).toBeInTheDocument();
  });
});
