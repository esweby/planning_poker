import { createFileRoute, useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/helloWorld")({
  component: RouteComponent,
  beforeLoad: () => {
    const navigate = useNavigate();

    if (true) {
      navigate({ to: "/" });
    }
  },
});

function RouteComponent() {
  return <div>Hello "/helloWorld"!</div>;
}
