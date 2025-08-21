import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/room/$roomNumber")({
  component: Room,
});

function Room() {
  const { roomNumber } = Route.useParams();

  return (
    <main>
      <h2>Room: {roomNumber}</h2>
    </main>
  );
}
