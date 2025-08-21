import { useUser } from "../contexts/UserContext";
import avatar from "animal-avatar-generator";
import { useEffect, useState, type ReactNode } from "react";

interface UserDetailsForm {
  children: ReactNode;
}

const UserDetailsForm = ({ children }: UserDetailsForm) => {
  const { name, role, setName, setRole } = useUser();
  const [svg, setSVG] = useState<null | string>(null);

  useEffect(() => {
    setSVG(avatar(name, { size: 150, blackout: false }));
  }, [name]);

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
      <div>{children}</div>
    </div>
  );
};

export default UserDetailsForm;
