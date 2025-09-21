import { useUser } from "../../../contexts/UserContext";
import JoinRoomBtn from "../../JoinRoomBtn";
import Icon from "../../atoms/icon/Icon";
import Container from "../../atoms/containers/Container";

import cl from "./UserDetailsForm.module.css";
import LabelWithChild from "../../atoms/label/LabelWithChild";
import Input from "../../atoms/input/Input";
import Select from "../../atoms/select/Select";
import JoinButton from "../../molecules/joinButton/JoinButton";

interface UserDetailsForm {
  ctaText: string;
  handleClick: () => void;
}

const UserDetailsForm = ({ ctaText, handleClick }: UserDetailsForm) => {
  const { name, role, seed, setName, setRole, setSeed } = useUser();

  return (
    <Container display="block" type="div" className={cl.container}>
      <Container
        type="header"
        display="flex"
        justifyContent="center"
        alignItems="center"
        margin="0 0 2rem"
      >
        <Icon name={name} seed={seed} size="huge" />
      </Container>
      <LabelWithChild name="Name">
        <Input
          name="Name"
          placeholder="Tobias"
          value={name}
          onChange={setName}
        />
      </LabelWithChild>
      <LabelWithChild name="Avatar Seed">
        <Input
          name="avatar"
          placeholder="Random string?"
          value={seed}
          onChange={setSeed}
        />
      </LabelWithChild>
      <LabelWithChild name="Role">
        <Select
          value={role}
          onChange={setRole}
          options={[
            { display: "Choose your role" },
            { value: "Delivery Lead", display: "Delivery Lead" },
            { value: "Product Owner", display: "Product Owner" },
            { value: "Developer", display: "Developer" },
            { value: "Tester", display: "Tester" },
          ]}
        />
      </LabelWithChild>
      <Container display="block" type="div">
        <JoinButton onClick={handleClick} />
      </Container>
    </Container>
  );
};

export default UserDetailsForm;
