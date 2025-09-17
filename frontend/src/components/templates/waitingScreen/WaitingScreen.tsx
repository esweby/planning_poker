import UserList from "../../organisms/userLists/UserList";
import styles from "./WaitingScreen.module.css";

const WaitingScreen = ({ data }: { data: Message }) => {
  return (
    <main className={styles.container}>
      <h3 className={styles.title}>Waiting to start...</h3>
      <UserList data={data} />
    </main>
  );
};

export default WaitingScreen;
