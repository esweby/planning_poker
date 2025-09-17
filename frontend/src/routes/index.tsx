import { createFileRoute, useNavigate } from "@tanstack/react-router";
import UserDetailsForm from "../components/organisms/userDetailsForm/UserDetailsForm";
import { useUser } from "../contexts/UserContext";
import Title from "../components/atoms/title/Title";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const { name, role } = useUser();
  const navigate = useNavigate();

  const handleJoinRoom = async () => {
    if (!name || !role) {
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/api/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          owner: {
            username: name,
            role: role,
          },
        }),
      });

      if (!res.ok) throw new Error("Failed to create room");

      const room = await res.json();

      navigate({ to: `/room/${room.roomNumber}` });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      <Title level={1} size="lg" weight="bold">
        Planning Poker
      </Title>
      <main className="index--main">
        <UserDetailsForm ctaText="Start Session" handleClick={handleJoinRoom} />
      </main>
    </>
  );
}
