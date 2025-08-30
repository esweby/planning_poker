import { createFileRoute } from "@tanstack/react-router";
import UserDetailsForm from "../components/UserDetailsForm";
import JoinHeader from "../components/JoinHeader";
import { useUser } from "../contexts/UserContext";

export const Route = createFileRoute("/join/$roomNumber")({
  component: Join,
});

function Join() {
  const { name, role } = useUser();
  const { roomNumber } = Route.useParams();

  const handleJoinRoom = () => {
    if (!name || !role) {
      return;
    }
  };

  return (
    <>
      <JoinHeader />
      <main className="index--main">
        <UserDetailsForm ctaText="Join Room" handleClick={handleJoinRoom} />
      </main>
    </>
  );
}
