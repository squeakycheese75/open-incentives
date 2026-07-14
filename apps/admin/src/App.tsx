import { RouterProvider } from "react-router-dom";

import { AuthProvider } from "./features/auth/AuthContext";
import { ToastProvider } from "./components/ui/ToastProvider";
import { router } from "./routes/router";

export default function App() {
  return (
    <AuthProvider>
      <ToastProvider>
        <RouterProvider router={router} />
      </ToastProvider>
    </AuthProvider>
  );
}
