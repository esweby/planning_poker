import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { UserProvider } from "../contexts/UserContext";
import Container from "../components/atoms/containers/Container";

import cl from "./root.module.css";

export const Route = createRootRoute({
  component: () => (
    <>
      <Container
        display="block"
        type="div"
        margin="2.1rem auto 0"
        className={cl.container}
      >
        <UserProvider>
          <Outlet />
        </UserProvider>
      </Container>
      <TanStackRouterDevtools />
    </>
  ),
});
