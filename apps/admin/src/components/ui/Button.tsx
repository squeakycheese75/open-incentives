import type { ButtonHTMLAttributes } from "react";

import { cn } from "../../lib/cn";

type Variant = "primary" | "secondary" | "danger" | "ghost";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

const variantClasses: Record<Variant, string> = {
  primary: "bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-50",
  secondary: "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 disabled:opacity-50",
  danger: "bg-red-600 text-white hover:bg-red-700 disabled:opacity-50",
  ghost: "bg-transparent text-gray-700 hover:bg-gray-100 disabled:opacity-50",
};

export function Button({ variant = "primary", className, ...props }: ButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center gap-1.5 rounded px-3.5 py-2 text-sm font-medium transition-colors",
        variantClasses[variant],
        className,
      )}
      {...props}
    />
  );
}
