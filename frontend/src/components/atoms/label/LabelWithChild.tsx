import type { ReactNode } from "react";
import cl from "./Label.module.css";

interface LabelProps {
  name: string;
  children: ReactNode;
  className?: string;
}

const LabelWithChild = ({ name, children, className }: LabelProps) => {
  let classes = cl.labelWithChild;
  if (className) classes += ` ${className}`;

  return (
    <label className={classes}>
      <span>{name}</span>
      {children}
    </label>
  );
};

export default LabelWithChild;
