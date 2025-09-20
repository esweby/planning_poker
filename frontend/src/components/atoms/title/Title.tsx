import { forwardRef, type ForwardedRef } from "react";
import type { ReactNode } from "react";
import cl from "./Title.module.css";

type Level = 1 | 2 | 3 | 4;
type Weight = "light" | "normal" | "bold";
type Size = "sm" | "md" | "lg" | "xl";

export interface TitleProps {
  level: Level;
  children: string[] | string | ReactNode | ReactNode[];
  weight: Weight;
  size: Size;
  marginOff?: boolean;
  className?: string;
}

const Title = forwardRef(
  (props: TitleProps, ref: ForwardedRef<HTMLElement>) => {
    const {
      level,
      children,
      weight,
      size,
      marginOff = false,
      className,
    } = props;

    const Tag: any = `h${level}`;

    let classes = `${cl.title} ${cl[size]} ${cl[weight]}`;
    if (marginOff) classes += ` ${cl.noMargin}`;
    if (className) classes += ` ${className}`;

    return (
      <Tag ref={ref} className={classes}>
        {children}
      </Tag>
    );
  }
);

export default Title;
