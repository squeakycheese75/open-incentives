import { createBrowserRouter, Navigate } from "react-router-dom";

import { AppShell } from "../components/layout/AppShell";
import { LoginPage } from "../features/auth/LoginPage";
import { RequireAuth } from "../features/auth/RequireAuth";
import { ProjectsListPage } from "../features/projects/ProjectsListPage";
import { ApiKeysPage } from "../features/apikeys/ApiKeysPage";
import { CampaignsListPage } from "../features/campaigns/CampaignsListPage";
import { CampaignFormPage } from "../features/campaigns/CampaignFormPage";

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    element: <RequireAuth />,
    children: [
      {
        element: <AppShell />,
        children: [
          { path: "/", element: <Navigate to="/projects" replace /> },
          { path: "/projects", element: <ProjectsListPage /> },
          { path: "/projects/:projectId/campaigns", element: <CampaignsListPage /> },
          { path: "/projects/:projectId/campaigns/new", element: <CampaignFormPage mode="create" /> },
          { path: "/projects/:projectId/campaigns/:campaignId/edit", element: <CampaignFormPage mode="edit" /> },
          { path: "/projects/:projectId/api-keys", element: <ApiKeysPage /> },
        ],
      },
    ],
  },
]);
