import cl from "./Text.module.css";

type Type = "p" | "span";
type Size = "sm" | "nm" | "lg" | "xl";
type Weight = "light" | "normal" | "bold";

export interface TitleProps {
  type?: Type;
  size?: Size;
  weight?: Weight;
  children: string | string[];
  className?: string;
}

const Text = (props: TitleProps) => {
  const {
    type = "p",
    size = "nm",
    weight = "normal",
    children,
    className,
  } = props;

  const Tag: any = `${type}`;

  let classes = `${cl[size]} ${cl[weight]}`;
  if (className) classes += ` ${className}`;

  return <Tag className={classes}>{children}</Tag>;
};

export default Text;
