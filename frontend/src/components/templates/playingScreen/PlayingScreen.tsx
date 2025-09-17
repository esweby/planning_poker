import { useUser } from "../../../contexts/UserContext";
import avatar from "animal-avatar-generator";
import styles from "./PlayingScreen.module.css";

const voteValues = ["0", "1", "2", "3", "6", "8", "10", "12", "?"];

const PlayingScreen = ({
  data,
  sendMessage,
}: {
  data: Message;
  sendMessage: (arg: any) => void;
}) => {
  const { name, seed } = useUser();

  const currUser = data.users[name];
  const votes = data.game.voted;
  const myVote = data.game.votes[name];

  const makeVote = (vote: string) =>
    sendMessage({ type: "vote", payload: { vote: vote } });

  console.log(votes);
  console.log(data.users);

  return (
    <main className={styles.container}>
      <section>
        <h2 className={styles.cta}>Hey {currUser.username} make your vote!</h2>
        <div className={styles.votesContainer}>
          {voteValues.map((v) => {
            let classes = `${styles.voteBtn}`;
            if (myVote === v) {
              classes += ` ${styles.chosen}`;
            }

            return (
              <button key={v} onClick={() => makeVote(v)} className={classes}>
                {v}
              </button>
            );
          })}
        </div>
      </section>
      <h2 className={styles.votedTitle}>Voted</h2>
      <section className={styles.votedContainer}>
        <ul className={styles.voteList}>
          {Object.keys(data.users).map((user: string) => (
            <li
              key={`userVote-${user}`}
              className={votes?.[user] ? styles.voted : styles.notVoted}
            >
              <img
                src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(data.users[user]?.seed ?? data.users[user].username, { size: 40, blackout: false }))}`}
                className={styles.image}
                alt="avatar"
              />
              {user}
            </li>
          ))}
        </ul>
        <div>
          <ul className={styles.votedStats}>
            <li>
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Delivery Lead" && votes[curr]) {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              of{" "}
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Delivery Lead") {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              DLs have voted.
            </li>
            <li>
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Product Owner" && votes[curr]) {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              of{" "}
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Product Owner") {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              POs have voted.
            </li>
            <li>
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Developer" && votes[curr]) {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              of{" "}
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Developer") {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              Devs have voted.
            </li>
            <li>
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Tester" && votes[curr]) {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              of{" "}
              {Object.keys(data.users).reduce((prev, curr) => {
                const currRole = data.users[curr].role;
                if (currRole === "Tester") {
                  return prev + 1;
                }

                return prev;
              }, 0)}{" "}
              Testers have voted.
            </li>
          </ul>
        </div>
      </section>
    </main>
  );
};

export default PlayingScreen;
