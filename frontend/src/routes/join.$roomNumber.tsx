import { createFileRoute, useNavigate } from "@tanstack/react-router";
import UserDetailsForm from "../components/organisms/userDetailsForm/UserDetailsForm";
import { useUser } from "../contexts/UserContext";
import Container from "../components/atoms/containers/Container";
import Title from "../components/atoms/title/Title";

export const Route = createFileRoute("/join/$roomNumber")({
  component: Join,
});

function Join() {
  const navigate = useNavigate();
  const { name, role } = useUser();
  const { roomNumber } = Route.useParams();

  const handleJoinRoom = async () => {
    if (!name || !role) return;

    try {
      const res = await fetch(`http://localhost:8080/api/check/${roomNumber}`, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      });

      console.log(res);
      if (!res.ok) throw new Error("Failed to find room");

      navigate({ to: `/room/${roomNumber}` });
    } catch (error) {
      console.error(error);
    }
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
        <UserDetailsForm ctaText="Join Room" handleClick={handleJoinRoom} />
      </Container>
    </>
  );
}
