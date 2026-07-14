import type { HTMLAttributes } from "react";

import { cn } from "../../lib/cn";

type Tone = "green" | "gray" | "red";

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  tone?: Tone;
}

const toneClasses: Record<Tone, string> = {
  green: "bg-green-100 text-green-800",
  gray: "bg-gray-100 text-gray-700",
  red: "bg-red-100 text-red-700",
};

export function Badge({ tone = "gray", className, ...props }: BadgeProps) {
  return (
    <span
      className={cn("inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium", toneClasses[tone], className)}
      {...props}
    />
  );
}
