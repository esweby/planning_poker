import { useUser } from "../../../contexts/UserContext";
import avatar from "animal-avatar-generator";
import styles from "./PlayingScreen.module.css";
import Title from "../../atoms/title/Title";
import Container from "../../atoms/containers/Container";
import UserList from "../../organisms/userLists/UserList";

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

  return (
    <Container
      display="block"
      type="main"
      className={styles.container}
      padding="1.5rem 0"
    >
      <Container display="block" type="section">
        <Title level={2} size="lg" weight="bold">
          Hey {currUser.username} make your vote!
        </Title>
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
      </Container>
      <Title level={2} size="lg" weight="bold">
        Voted
      </Title>
      <section className={styles.votedContainer}>
        <UserList data={data} showVoted={true} />
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
    </Container>
  );
};

export default PlayingScreen;
