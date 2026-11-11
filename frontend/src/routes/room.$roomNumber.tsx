import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useUser } from "../contexts/UserContext";
import { useWebSocket } from "../hooks/useWebsocket";
import { useEffect, useState, useRef } from "react";
import HeaderInfo from "../components/molecules/headerInfo/HeaderInfo";
import OwnerControls from "../components/molecules/ownerControls/OwnerControls";
import WaitingScreen from "../components/templates/waitingScreen/WaitingScreen";
import PlayingScreen from "../components/templates/playingScreen/PlayingScreen";
import FinishedStyles from "../components/templates/finishedScreen/FinishedScreen";
import Container from "../components/atoms/containers/Container";
import Text from "../components/atoms/text/Text";

export const Route = createFileRoute("/room/$roomNumber")({
  component: Room,
});

function Room() {
  const navigate = useNavigate();

  const [forceRefresh, setforceRefresh] = useState(0);
  const [data, setData] = useState<Message | null>();
  const { roomNumber } = Route.useParams();
  const { name, role, seed } = useUser();

  const hasConnected = useRef(false);

  const { isConnected, messages, error, connect, closeConnection, sendMessage } =
    useWebSocket();

  useEffect(() => {
    if (!roomNumber || !name || !role) {
      navigate({ to: `${window.location.origin}/join/${roomNumber}` });
      return;
    }

    if (hasConnected.current) return;

    hasConnected.current = true;
    connect(roomNumber, name, role, seed);

    return () => {
      closeConnection();
      hasConnected.current = false;
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

  if (error && (!data || !isConnected)) return <Text>{error}</Text>;
  if (!data) return;
  const status = data.game.status;

  return (
    <>
      <Container display="block" type="header">
        {error && <Text>{error}</Text>}
        <HeaderInfo owner={data.owner} roomId={roomNumber} />
        {name === data.owner && (
          <OwnerControls
            gameStatus={data.game.status}
            sendMessage={sendMessage}
          />
        )}
      </Container>
      {status === "waiting" && <WaitingScreen data={data} />}
      {status === "playing" && (
        <PlayingScreen data={data} sendMessage={sendMessage} />
      )}
      {status === "finished" && <FinishedStyles data={data} />}
    </>
  );
}
