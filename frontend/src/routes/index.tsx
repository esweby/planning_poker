import { createFileRoute, useNavigate } from "@tanstack/react-router";
import UserDetailsForm from "../components/organisms/userDetailsForm/UserDetailsForm";
import { useUser } from "../contexts/UserContext";
import Title from "../components/atoms/title/Title";
import Text from "../components/atoms/text/Text";
import { useState } from "react";
import { API_URL } from "../config";
import { setSessionToken } from "../session";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const { name, role } = useUser();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  const handleJoinRoom = async () => {
    if (!name || !role) return;
    setError(null);

    try {
      const res = await fetch(`${API_URL}/api/`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          owner: {
            username: name,
            role: role,
          },
        }),
      });

      const room = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(room.error ?? "Failed to create room");

      setSessionToken(room.roomNumber, name, room.token);
      navigate({ to: `/room/${room.roomNumber}` });
    } catch (error) {
      console.error(error);
      setError(
        error instanceof Error ? error.message : "Failed to create room"
      );
    }
  };

  return (
    <>
      <Title level={1} size="lg" weight="bold">
        Planning Poker
      </Title>
      <main className="index--main">
        {error && <Text>{error}</Text>}
        <UserDetailsForm ctaText="Start Session" handleClick={handleJoinRoom} />
      </main>
    </>
  );
}
