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
    <Container type="section" display="flex" gap="5px" margin="1rem 0 1rem">
      {(gameStatus === "waiting" || gameStatus === "finished") && (
        <>
          <Button onClick={handleStartVoting} small>
            Start Voting
          </Button>
        </>
      )}
      {gameStatus === "playing" && (
        <>
          <Button onClick={handleRevealVotes} small>
            Reveal Votes
          </Button>
          <Button onClick={handleRestartVoting} small>
            Restart
          </Button>
        </>
      )}
    </Container>
  );
};

export default OwnerControls;
