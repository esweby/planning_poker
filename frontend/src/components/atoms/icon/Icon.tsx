import avatar from "animal-avatar-generator";
import cl from "./Icon.module.css";

type Size = "sm" | "md" | "lg" | "huge";

export interface IconProps {
  size: Size;
  seed?: string;
  name: string;
  centered?: boolean;
  className?: string;
}

const Icon = (props: IconProps) => {
  const { size, seed, name, centered = true, className } = props;

  const seedString = (seed && seed.length > 0 ? seed : name) || "base";

  let classes = `${cl.icon} ${cl[size]}`;
  if (!centered) classes += ` ${cl.noMargin}`;
  if (className) classes += ` ${className}`;

  return (
    <img
      src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(seedString, { size: 150, blackout: false }))}`}
      className={classes}
      alt="avatar"
    />
  );
};

export default Icon;
