import styles from "./FinishedScreen.module.css";
import avatar from "animal-avatar-generator";

const FinishedStyles = ({ data }: { data: Message }) => {
  const votes = data.game.votes;

  const results: { [key: string]: number } = {};

  Object.values(votes).forEach((vote: string) => {
    if (results[vote]) {
      results[vote] += 1;
    } else {
      results[vote] = 1;
    }
  });

  const highestValues = Object.keys(results).sort(
    (a, b) => results[b] - results[a]
  );

  console.log(votes);

  const isItATie =
    highestValues.length > 1 && highestValues[0] === highestValues[1];

  console.log(highestValues.length === 1);
  console.log(highestValues[0] === highestValues[1]);
  console.log(highestValues);
  console.log(results);

  return (
    <main className={styles.container}>
      <h2 className={styles.title}>
        And the result is ... {isItATie ? "It's a tie!" : highestValues[0]}
      </h2>
      <div className={styles.results}>
        <section>
          {Object.keys(data.users).map((user: string) => (
            <li
              key={`userVote-${user}`}
              className={votes?.[user] ? styles.voted : styles.notVoted}
            >
              <div>
                <img
                  src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(data.users[user]?.seed ?? data.users[user].username, { size: 40, blackout: false }))}`}
                  className={styles.image}
                  alt="avatar"
                />
              </div>
              <div>{user}</div>
              <div>{votes[user]}</div>
            </li>
          ))}
        </section>
        <section>
          <ul className={styles.voteList}>
            {Object.keys(results)
              .sort((a, b) => results[b] - results[a])
              .map((result) => (
                <li key={`result-${result}`}>
                  {result}: {results[result]}
                </li>
              ))}
          </ul>
        </section>
        <section></section>
      </div>
    </main>
  );
};

export default FinishedStyles;
