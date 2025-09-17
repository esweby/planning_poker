import UserListItem from "../../molecules/userListItem/UserListItem";
import Container from "../../atoms/containers/Container";

import cl from "./UserList.module.css";

const UserList = ({
  data,
  showVoted = false,
  showVote = false,
}: {
  data: Message;
  showVoted?: boolean;
  hasVoted?: boolean;
  showVote?: boolean;
}) => {
  return (
    <Container type="ul" display="block" className={cl.noList}>
      {Object.keys(data.users).map((username: string) => {
        let hasVoted = showVoted ? data.game.voted[username] : false;
        let vote = data.game.votes[username] || "0";

        return (
          <UserListItem
            key={`id-UserListItem-${username}`}
            user={data.users[username]}
            showVoted={showVoted}
            hasVoted={hasVoted}
            showVote={showVote}
            vote={vote}
          />
        );
      })}
    </Container>
  );
};

export default UserList;
