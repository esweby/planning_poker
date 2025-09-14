import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useUser } from "../contexts/UserContext";
import { useWebSocket } from "../hooks/useWebsocket";
import { useEffect, useState } from "react";
import ListUserItem from "../components/room/ListUserItem";
import HeaderInfo from "../components/room/HeaderInfo";
import OwnerControls from "../components/room/OwnerControls";
import WaitingScreen from "../components/room/screens/WaitingScreen";
import PlayingScreen from "../components/room/screens/PlayingScreen";
import FinishedStyles from "../components/room/screens/FinishedScreen";

export const Route = createFileRoute("/room/$roomNumber")({
  component: Room,
});

function Room() {
  const navigate = useNavigate();

  const [forceRefresh, setforceRefresh] = useState(0);
  const [data, setData] = useState<Message | null>();
  const { roomNumber } = Route.useParams();
  const { name, role, seed } = useUser();

  const { isConnected, messages, connect, closeConnection, sendMessage } =
    useWebSocket();

  useEffect(() => {
    if (!roomNumber || !name || !role) {
      navigate({ to: `${window.location.origin}/join/${roomNumber}` });
      return;
    }

    connect(
      `ws://localhost:8080/api/join/${roomNumber}?username=${encodeURIComponent(name)}&role=${encodeURIComponent(role)}&seed=${encodeURIComponent(seed)}`
    );

    return () => {
      closeConnection();
    };
  }, [roomNumber, name, role]);

  useEffect(() => {
    if (!isConnected || messages.length === 0) {
      if (forceRefresh < 5) {
        setTimeout(() => setforceRefresh((curr) => curr + 1));
      }
      return;
    }

    setData(messages[0].payload);
  }, [isConnected, messages]);

  if (!data) return;

  const status = data.game.status;

  return (
    <div>
      <header>
        <HeaderInfo owner={data.owner} roomId={roomNumber} />
        {name === data.owner && (
          <OwnerControls
            gameStatus={data.game.status}
            sendMessage={sendMessage}
          />
        )}
      </header>
      {status === "waiting" && <WaitingScreen data={data} />}
      {status === "playing" && (
        <PlayingScreen data={data} sendMessage={sendMessage} />
      )}
      {status === "finished" && <FinishedStyles data={data} />}
    </div>
  );
}
