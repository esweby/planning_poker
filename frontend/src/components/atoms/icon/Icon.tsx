import avatar from "animal-avatar-generator";
import cl from "./Icon.module.css";

type Size = "sm" | "md" | "lg" | "huge";

export interface IconProps {
  size: Size;
  seed?: string;
  name: string;
}

const Icon = (props: IconProps) => {
  const { size, seed, name } = props;

  const seedString = (seed && seed.length > 0 ? seed : name) || "base";

  return (
    <img
      src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(seedString, { size: 150, blackout: false }))}`}
      className={`
        ${cl.icon} ${cl[size]}`}
      alt="avatar"
    />
  );
};

export default Icon;
