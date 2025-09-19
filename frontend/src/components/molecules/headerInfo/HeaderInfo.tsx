import Button from "../../atoms/button/Button";
import Container from "../../atoms/containers/Container";
import Title from "../../atoms/title/Title";

const HeaderInfo = ({ owner, roomId }: { owner: string; roomId: string }) => {
  const handleCopyInvite = () => {
    const inviteUrl = `${window.location.origin}/join/${roomId}`;

    navigator.clipboard
      .writeText(inviteUrl)
      .then(() => {})
      .catch((err) => {
        console.error("Failed to copy invite link:", err);
      });
  };

  return (
    <Container type="section" display="flex" justifyContent="between">
      <Title level={2} weight="bold" size="lg">
        {owner}'s room
      </Title>
      <Button onClick={handleCopyInvite}>Copy Invite Link</Button>
    </Container>
  );
};

export default HeaderInfo;
