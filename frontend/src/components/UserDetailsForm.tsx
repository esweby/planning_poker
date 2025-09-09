import { useUser } from "../contexts/UserContext";
import avatar from "animal-avatar-generator";
import { useEffect, useState } from "react";
import JoinRoomBtn from "./JoinRoomBtn";

interface UserDetailsForm {
  ctaText: string;
  handleClick: () => void;
}

const UserDetailsForm = ({ ctaText, handleClick }: UserDetailsForm) => {
  const { name, role, seed, setName, setRole, setSeed } = useUser();
  const [svg, setSVG] = useState<null | string>(null);

  useEffect(() => {
    setSVG(
      avatar(seed.length > 0 ? seed : name, { size: 150, blackout: false })
    );
  }, [seed, name]);

  return (
    <div className="container--user-details">
      {svg && (
        <img
          src={`data:image/svg+xml;utf8,${encodeURIComponent(svg)}`}
          className="image--avatar-preview"
          alt="avatar"
        />
      )}
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
          <option value="developer">Developer</option>
          <option value="tester">Tester</option>
          <option value="po">Product Owner</option>
        </select>
      </label>
      <div>
        <JoinRoomBtn onClick={handleClick}>{ctaText}</JoinRoomBtn>
      </div>
    </div>
  );
};

export default UserDetailsForm;
