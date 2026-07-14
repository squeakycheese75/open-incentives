import type { SelectHTMLAttributes } from "react";

import { cn } from "../../lib/cn";

export function Select({ className, ...props }: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "w-full rounded border border-gray-300 bg-white px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500",
        className,
      )}
      {...props}
    />
  );
}
