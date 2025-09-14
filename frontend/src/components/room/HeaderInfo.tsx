import styles from "./HeaderInfo.module.css";

const HeaderInfo = ({ owner, roomId }: { owner: string; roomId: string }) => {
  const handleCopyInvite = () => {
    const inviteUrl = `${window.location.origin}/join/${roomId}`;

    navigator.clipboard
      .writeText(inviteUrl)
      .then(() => {})
      .catch((err) => {
        console.error("Failed to copy invite link:", err);
      });
  };

  return (
    <section className={styles.container}>
      <h2>{owner}'s room</h2>
      <button className="button--copy-invite" onClick={handleCopyInvite}>
        Copy Invite Link
      </button>
    </section>
  );
};

export default HeaderInfo;
