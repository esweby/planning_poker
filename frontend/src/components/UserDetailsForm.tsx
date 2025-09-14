import { useUser } from "../contexts/UserContext";
import JoinRoomBtn from "./JoinRoomBtn";
import Icon from "./atoms/icon/Icon";
import Container from "./atoms/containers/Container";

interface UserDetailsForm {
  ctaText: string;
  handleClick: () => void;
}

const UserDetailsForm = ({ ctaText, handleClick }: UserDetailsForm) => {
  const { name, role, seed, setName, setRole, setSeed } = useUser();

  return (
    <div className="container--user-details">
      <Container
        type="header"
        display="flex"
        justifyContent="center"
        alignItems="center"
        margin="0 0 2rem"
      >
        <Icon name={name} seed={seed} size="huge" />
      </Container>
      <label className="label--user-details">
        <span>Name</span>
        <input
          type="text"
          placeholder="Bob"
          className="input--name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </label>
      <label className="label--user-details">
        <span>Avatar Seed</span>
        <input
          type="text"
          placeholder="Type till you like your avatar"
          className="input--name"
          value={seed}
          onChange={(e) => setSeed(e.target.value)}
        />
      </label>
      <label className="label--user-details">
        <span>Role</span>
        <select
          className="select--role"
          value={role}
          onChange={(e) => setRole(e.target.value as any)}
        >
          <option>Choose your role</option>
          <option value="Delivery Lead">Delivery Lead</option>
          <option value="Product Owner">Product Owner</option>
          <option value="Developer">Developer</option>
          <option value="Tester">Tester</option>
        </select>
      </label>
      <div>
        <JoinRoomBtn onClick={handleClick}>{ctaText}</JoinRoomBtn>
      </div>
    </div>
  );
};

export default UserDetailsForm;
