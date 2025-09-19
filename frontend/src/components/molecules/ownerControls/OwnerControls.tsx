import Button from "../../atoms/button/Button";
import Container from "../../atoms/containers/Container";

const OwnerControls = ({
  gameStatus,
  sendMessage,
}: {
  gameStatus: string;
  sendMessage: (arg: { type: string }) => void;
}) => {
  const handleStartVoting = () => sendMessage({ type: "start_voting" });
  const handleRevealVotes = () => sendMessage({ type: "reveal_votes" });
  const handleRestartVoting = () => sendMessage({ type: "restart_voting" });

  return (
    <Container type="section" display="flex" gap="5px">
      {gameStatus === "waiting" ||
        (gameStatus === "finished" && (
          <>
            <Button onClick={handleStartVoting} dark>
              Start Voting
            </Button>
          </>
        ))}
      {gameStatus === "playing" && (
        <>
          <Button onClick={handleRevealVotes} dark>
            Reveal Votes
          </Button>
          <Button onClick={handleRestartVoting} dark>
            Restart
          </Button>
        </>
      )}
    </Container>
  );
};

export default OwnerControls;
