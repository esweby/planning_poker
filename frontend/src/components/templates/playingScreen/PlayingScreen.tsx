import { useUser } from "../../../contexts/UserContext";
import Title from "../../atoms/title/Title";
import Container from "../../atoms/containers/Container";
import UserList from "../../organisms/userLists/UserList";
import Button from "../../atoms/button/Button";

import cl from "./PlayingScreen.module.css";

const voteValues = ["0", "1", "2", "3", "6", "8", "10", "12", "?"];

const PlayingScreen = ({
  data,
  sendMessage,
}: {
  data: Message;
  sendMessage: (arg: any) => void;
}) => {
  const { name } = useUser();

  const currUser = data.users[name];
  const votes = data.game.voted;
  const myVote = data.game.votes[name];

  const makeVote = (vote: string) =>
    sendMessage({ type: "vote", payload: { vote: vote } });

  const countNumVotes = (name: string) => {
    return Object.keys(data.users).reduce((prev, curr) => {
      const currRole = data.users[curr].role;
      if (currRole === name && votes[curr]) {
        return prev + 1;
      }

      return prev;
    }, 0);
  };

  const countNumRole = (name: string): number => {
    return Object.keys(data.users).reduce((prev, curr) => {
      const currRole = data.users[curr].role;
      if (currRole === name) {
        return prev + 1;
      }

      return prev;
    }, 0);
  };

  return (
    <Container
      display="block"
      type="main"
      className={cl.container}
      padding="1.5rem 0"
    >
      <Container display="block" type="section">
        <Title level={2} size="lg" weight="bold">
          Hey {currUser.username} make your vote!
        </Title>
        <Container type="div" display="flex" gap="0.5rem" padding="0 0 1rem">
          {voteValues.map((v) => {
            let classes = `${cl.voteBtn}`;
            if (myVote === v) {
              classes += ` ${cl.chosen}`;
            }

            return (
              <Button
                key={v}
                onClick={() => makeVote(v)}
                className={classes}
                small
              >
                {v}
              </Button>
            );
          })}
        </Container>
      </Container>
      <Title level={2} size="lg" weight="bold">
        Voted
      </Title>
      <Container display="grid" type="section" className={cl.votedContainer}>
        <UserList data={data} showVoted={true} />
        <Container display="block" type="div">
          <Container display="block" type="ul" className={cl.votedStats}>
            <Container display="block" type="li">
              {countNumVotes("Delivery Lead")} of{" "}
              {countNumRole("Delivery Lead")} DLs have voted.
            </Container>
            <Container display="block" type="li">
              {countNumVotes("Product Owner")} of{" "}
              {countNumRole("Product Owner")} POs have voted.
            </Container>
          </Container>
          <Container display="block" type="li">
            {countNumVotes("Developer")} of {countNumRole("Developer")} Devs
            have voted.
          </Container>
          <Container display="block" type="li">
            {countNumVotes("Tester")} of {countNumRole("Tester")} Testers have
            voted.
          </Container>
        </Container>
      </Container>
    </Container>
  );
};

export default PlayingScreen;
