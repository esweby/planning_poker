import { createFileRoute, useNavigate } from "@tanstack/react-router";
import UserDetailsForm from "../components/organisms/userDetailsForm/UserDetailsForm";
import { useUser } from "../contexts/UserContext";
import Container from "../components/atoms/containers/Container";
import Title from "../components/atoms/title/Title";
import Text from "../components/atoms/text/Text";
import { useState } from "react";
import { API_URL } from "../config";
import { checkJoin } from "../session";

export const Route = createFileRoute("/join/$roomNumber")({
  component: Join,
});

function Join() {
  const navigate = useNavigate();
  const { name, role } = useUser();
  const { roomNumber } = Route.useParams();
  const [error, setError] = useState<string | null>(null);

  const handleJoinRoom = async () => {
    if (!name || !role) return;

    const joinError = await checkJoin(API_URL, roomNumber, name);
    setError(joinError);
    if (!joinError) navigate({ to: `/room/${roomNumber}` });
  };

  return (
    <>
      <Title level={1} weight="bold" size="xl">
        Planning Poker
      </Title>
      <Container
        display="flex"
        type="main"
        justifyContent="center"
        padding="4rem 0 0"
      >
        {error && <Text>{error}</Text>}
        <UserDetailsForm ctaText="Join Room" handleClick={handleJoinRoom} />
      </Container>
    </>
  );
}
