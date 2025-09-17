import Container from "../../atoms/containers/Container";
import Title from "../../atoms/title/Title";
import UserList from "../../organisms/userLists/UserList";
import Text from "../../atoms/text/Text";
import cl from "./FinishedScreen.module.css";

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

  const isItATie =
    highestValues.length > 1 && highestValues[0] === highestValues[1];

  return (
    <Container display="block" type="main" padding="1rem 0 0">
      <Title level={2} size="lg" weight="bold" className={cl.title}>
        And the result is ... {isItATie ? "It's a tie!" : highestValues[0]}
      </Title>
      <Container
        display="flex"
        type="div"
        justifyContent="between"
        padding="0 2rem"
      >
        <Container display="block" type="section">
          <Title level={3} size="md" weight="bold">
            Votes
          </Title>
          <UserList data={data} showVote={true} />
        </Container>
        <Container display="block" type="section">
          <Title level={3} size="md" weight="bold">
            Voted &nbsp;&nbsp;Votes
          </Title>
          <Container display="block" type="ul" className={cl.voteList}>
            {Object.keys(results)
              .sort((a, b) => results[b] - results[a])
              .map((result) => (
                <Container
                  key={`result-${result}`}
                  display="flex"
                  type="li"
                  justifyContent="between"
                  padding="0 1.25rem"
                >
                  <Text type="span" weight="bold" size="xl">
                    {result}
                  </Text>
                  <Text type="span" weight="normal" size="xl">
                    {results[result].toString()}
                  </Text>
                </Container>
              ))}
          </Container>
        </Container>
      </Container>
    </Container>
  );
};

export default FinishedStyles;
