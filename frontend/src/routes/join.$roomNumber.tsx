import { createFileRoute } from "@tanstack/react-router";
import UserDetailsForm from "../components/UserDetailsForm";
import JoinRoomBtn from "../components/JoinRoomBtn";
import JoinHeader from "../components/JoinHeader";

export const Route = createFileRoute("/join/$roomNumber")({
  component: Join,
});

function Join() {
  const { roomNumber } = Route.useParams();

  return (
    <>
      <JoinHeader />
      <main className="index--main">
        <UserDetailsForm>
          <JoinRoomBtn onClick={() => {}}>Join Room</JoinRoomBtn>
        </UserDetailsForm>
      </main>
    </>
  );
}
