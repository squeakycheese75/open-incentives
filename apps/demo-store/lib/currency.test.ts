import { describe, expect, it } from "vitest";
import { formatCurrency } from "./currency";

describe("formatCurrency", () => {
  it("formats EUR amounts", () => {
    expect(formatCurrency(50, "EUR")).toBe("€50.00");
  });

  it("formats zero", () => {
    expect(formatCurrency(0, "EUR")).toBe("€0.00");
  });
});
