import { createFileRoute } from "@tanstack/react-router";
import avatar from "animal-avatar-generator";
import { useUser } from "../contexts/UserContext";
import { useWebSocket } from "../hooks/useWebsocket";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/room/$roomNumber")({
  component: Room,
});

function Room() {
  const [roomDetails, setRoomDetails] = useState({ owner: "" });
  const [gameDetails, setGameDetails] = useState({});
  const { roomNumber } = Route.useParams();
  const { name, role } = useUser();
  const { isConnected, messages, sendMessage } = useWebSocket(
    `ws://localhost:8080/api/join/${roomNumber}?username=${encodeURIComponent(name)}&role=${encodeURIComponent(role)}`
  );

  useEffect(() => {
    if (!isConnected || messages.length === 0) {
      return;
    }

    setRoomDetails({
      owner: messages[0].payload.owner,
    });

    setGameDetails({
      ...messages[0].payload.game,
    });
  }, [isConnected, messages]);

  console.log({ messages });
  console.log({ sendMessage });

  return (
    <main>
      <h2 className="title--room">
        <img
          src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(name, { size: 25, blackout: false }))}`}
          className="image--avatar-room-title"
          alt="avatar"
        />{" "}
        {roomDetails.owner} Planning Poker
      </h2>
      {name === roomDetails.owner && (
        <header>
          Status: {gameDetails.status} <button>Play</button>
        </header>
      )}
    </main>
  );
}
