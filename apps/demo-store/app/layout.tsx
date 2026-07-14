import type { Metadata } from "next";
import type { ReactNode } from "react";
import { CartProvider } from "@/context/CartContext";
import { Header } from "@/components/Header";
import "./globals.css";

export const metadata: Metadata = {
  title: "Open Incentives Demo Store",
  description: "A demo storefront showing integration with the Open Incentives API.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <CartProvider>
          <Header />
          <main className="page">{children}</main>
        </CartProvider>
      </body>
    </html>
  );
}
