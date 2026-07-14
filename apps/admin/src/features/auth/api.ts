import { useMutation } from "@tanstack/react-query";

import { apiFetch } from "../../api/client";
import type { LoginRequest, LoginResponse } from "../../api/types";

export function useLogin() {
  return useMutation({
    mutationFn: (body: LoginRequest) =>
      apiFetch<LoginResponse>("/admin/auth/login", {
        method: "POST",
        body: JSON.stringify(body),
      }),
  });
}
