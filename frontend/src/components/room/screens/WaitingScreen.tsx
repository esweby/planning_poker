import ListUserItem from "../ListUserItem";
import styles from "./WaitingScreen.module.css";

const WaitingScreen = ({ data }: { data: Message }) => {
  return (
    <main className={styles.container}>
      <h3 className={styles.title}>Waiting to start...</h3>
      {Object.keys(data.users).map((username: string) => (
        <ListUserItem key={username} user={data.users[username]} />
      ))}
    </main>
  );
};

export default WaitingScreen;
