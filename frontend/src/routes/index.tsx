import { createFileRoute } from "@tanstack/react-router";
import UserDetailsForm from "../components/UserDetailsForm";
import JoinRoomBtn from "../components/JoinRoomBtn";
import JoinHeader from "../components/JoinHeader";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <>
      <JoinHeader />
      <main className="index--main">
        <UserDetailsForm>
          <JoinRoomBtn onClick={() => {}}>Start Session</JoinRoomBtn>
        </UserDetailsForm>
      </main>
    </>
  );
}
