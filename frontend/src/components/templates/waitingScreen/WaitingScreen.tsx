import Container from "../../atoms/containers/Container";
import Title from "../../atoms/title/Title";
import UserList from "../../organisms/userLists/UserList";

const WaitingScreen = ({ data }: { data: Message }) => {
  return (
    <Container type="main" display="block">
      <Title level={3} weight="bold" size="md">
        Waiting to start...
      </Title>
      <UserList data={data} />
    </Container>
  );
};

export default WaitingScreen;
