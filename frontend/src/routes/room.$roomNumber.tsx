import { createFileRoute } from "@tanstack/react-router";
import avatar from "animal-avatar-generator";
import { useUser } from "../contexts/UserContext";
import { useWebSocket } from "../hooks/useWebsocket";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/room/$roomNumber")({
  component: Room,
});

function Room() {
  const [data, setData] = useState<Message | null>();
  const { roomNumber } = Route.useParams();
  const { name, role } = useUser();

  const { isConnected, messages, connect, closeConnection, sendMessage } =
    useWebSocket();

  useEffect(() => {
    if (!roomNumber || !name || !role) return;

    connect(
      `ws://localhost:8080/api/join/${roomNumber}?username=${encodeURIComponent(name)}&role=${encodeURIComponent(role)}&seed=${encodeURIComponent(role)}`
    );

    return () => {
      closeConnection();
    };
  }, [roomNumber, name, role]);

  useEffect(() => {
    if (!isConnected || messages.length === 0) {
      return;
    }
    console.log(messages);
    setData(messages[0].payload);
  }, [isConnected, messages]);

  console.log(data);

  if (!data) return;

  return (
    <div>
      <header>
        <section className="header--title">
          <h2 className="title--room">
            <img
              src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(name, { size: 25, blackout: false }))}`}
              className="image--avatar-room-title"
              alt="avatar"
            />{" "}
            {data.owner} Planning Poker
          </h2>
          <button className="button--copy-invite">Copy Invite Link</button>
        </section>
        {name === data.owner && (
          <>
            Status:{" "}
            <span className="room--game-status">{data.game.status}</span>
          </>
        )}
      </header>
      <main></main>
    </div>
  );
}
