import cl from "./Button.module.css";

interface ButtonProps {
  onClick: () => void;
  children: string | string[];
  border?: boolean;
  small?: boolean;
  medium?: boolean;
  large?: boolean;
  rounded?: boolean;
  thick?: boolean;
  dark?: boolean;
  disabled?: boolean;
  className?: string;
}

const Button = ({
  onClick,
  children,
  border = false,
  small = false,
  medium = true,
  large = false,
  rounded = false,
  thick = false,
  dark = false,
  disabled = false,
  className,
}: ButtonProps) => {
  let classes = cl.btn;
  if (border) classes += ` ${cl.border}`;
  if (small) classes += ` ${cl.small}`;
  if (medium) classes += ` ${cl.medium}`;
  if (large) classes += ` ${cl.large}`;
  if (rounded) classes += ` ${cl.rounded}`;
  if (thick) classes += ` ${cl.thick}`;
  if (dark) classes += ` ${cl.dark}`;
  if (className) classes += ` ${className}`;

  return (
    <button onClick={onClick} className={classes} disabled={disabled}>
      {children}
    </button>
  );
};

export default Button;
