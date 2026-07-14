import { NavLink, useParams } from "react-router-dom";

import { cn } from "../../lib/cn";
import { useAuth } from "../../features/auth/AuthContext";
import { ProjectSwitcher } from "../../features/projects/ProjectSwitcher";
import { Button } from "../ui/Button";

const linkClasses = ({ isActive }: { isActive: boolean }) =>
  cn("rounded px-3 py-1.5 text-sm font-medium", isActive ? "bg-primary-50 text-primary-700" : "text-gray-600 hover:bg-gray-100");

export function TopNav() {
  const { logout } = useAuth();
  const { projectId } = useParams<{ projectId: string }>();

  return (
    <header className="flex items-center justify-between border-b border-gray-200 bg-white px-6 py-3">
      <div className="flex items-center gap-6">
        <span className="text-sm font-semibold text-gray-900">Open Incentives</span>
        <nav className="flex items-center gap-1">
          <NavLink to="/projects" className={linkClasses}>
            Projects
          </NavLink>
          {projectId && (
            <>
              <NavLink to={`/projects/${projectId}/campaigns`} className={linkClasses}>
                Campaigns
              </NavLink>
              <NavLink to={`/projects/${projectId}/api-keys`} className={linkClasses}>
                API Keys
              </NavLink>
            </>
          )}
        </nav>
      </div>

      <div className="flex items-center gap-3">
        <ProjectSwitcher />
        <Button variant="ghost" onClick={logout}>
          Sign out
        </Button>
      </div>
    </header>
  );
}
