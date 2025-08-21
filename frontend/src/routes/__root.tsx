import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { UserProvider } from "../contexts/UserContext";

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="container">
        <UserProvider>
          <Outlet />
        </UserProvider>
      </div>
      <TanStackRouterDevtools />
    </>
  ),
});
