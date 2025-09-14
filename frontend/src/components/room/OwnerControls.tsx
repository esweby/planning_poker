import styles from "./OwnerControls.module.css";

const OwnerControls = ({
  gameStatus,
  sendMessage,
}: {
  gameStatus: string;
  sendMessage: (arg: { type: string }) => void;
}) => {
  const handleStartVoting = () => sendMessage({ type: "start_voting" });
  const handleRevealVotes = () => sendMessage({ type: "reveal_votes" });
  const handleResetVoting = () => sendMessage({ type: "reset_voting" });

  return (
    <section className={styles.container}>
      <button
        onClick={handleStartVoting}
        className={styles.btn}
        disabled={gameStatus !== "waiting"}
      >
        Start Voting
      </button>
      <button
        onClick={handleRevealVotes}
        className={styles.btn}
        disabled={gameStatus !== "playing"}
      >
        Reveal Votes
      </button>
      <button
        onClick={handleResetVoting}
        className={styles.btn}
        disabled={gameStatus === "waiting"}
      >
        Restart
      </button>
    </section>
  );
};

export default OwnerControls;
