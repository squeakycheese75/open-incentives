export type Product = {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
};

export const products: Product[] = [
  {
    id: "prod_coffee",
    name: "Premium Coffee",
    description: "Single-origin roasted coffee beans.",
    price: 18.0,
    currency: "EUR",
  },
  {
    id: "prod_mug",
    name: "Ceramic Mug",
    description: "Reusable ceramic coffee mug.",
    price: 14.0,
    currency: "EUR",
  },
  {
    id: "prod_grinder",
    name: "Coffee Grinder",
    description: "Compact manual coffee grinder.",
    price: 34.0,
    currency: "EUR",
  },
  {
    id: "prod_filter",
    name: "Coffee Filters",
    description: "Pack of reusable coffee filters.",
    price: 9.0,
    currency: "EUR",
  },
];
