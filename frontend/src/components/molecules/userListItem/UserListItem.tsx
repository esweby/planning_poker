import Container from "../../atoms/containers/Container";
import Icon from "../../atoms/icon/Icon";
import Text from "../../atoms/text/Text";
import Title from "../../atoms/title/Title";

import cl from "./UserListItem.module.css";

const UserListItem = ({
  user,
  showVoted,
  hasVoted,
  showVote,
  vote,
}: {
  user: User;
  showVoted: boolean;
  hasVoted: boolean;
  showVote: boolean;
  vote: string;
}) => {
  let classes = "";
  if (showVoted && !hasVoted) classes = cl.notVoted;
  if (showVoted && hasVoted) classes = cl.voted;

  return (
    <Container
      type="li"
      display="flex"
      justifyContent="start"
      alignItems="center"
      margin="0 0 .5rem"
    >
      <Icon
        size="md"
        name={user.username}
        seed={user.seed || ""}
        centered={false}
        className={cl.marginRight}
      />
      <Title
        level={3}
        weight="bold"
        size="md"
        marginOff={true}
        className={cl.marginRight}
      >
        {user.username}
      </Title>
      {showVote && vote && (
        <Text size="xl" weight="bold" className={cl.vote}>
          {vote}
        </Text>
      )}
    </Container>
  );
};

export default UserListItem;
