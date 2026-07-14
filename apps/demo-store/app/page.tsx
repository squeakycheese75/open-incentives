import { products } from "@/data/products";
import { ProductGrid } from "@/components/ProductGrid";

export default function StorefrontPage() {
  return (
    <div>
      <p>
        Browse the catalogue below. Discounts are evaluated live by Open Incentives
        as you build your cart.
      </p>
      <ProductGrid products={products} />
    </div>
  );
}
