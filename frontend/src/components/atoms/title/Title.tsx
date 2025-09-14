import { forwardRef, type ForwardedRef } from "react";
import cl from "./Title.module.css";

type Level = 1 | 2 | 3 | 4;
type Weight = "light" | "normal" | "bold";
type Size = "sm" | "md" | "lg" | "xl";

export interface TitleProps {
  level: Level;
  children: string;
  weight: Weight;
  size: Size;
}

const Title = forwardRef(
  (props: TitleProps, ref: ForwardedRef<HTMLElement>) => {
    const { level, children, weight, size } = props;

    const Tag: any = `h${level}`;

    let classes = constructClasses(size, weight);

    return (
      <Tag ref={ref} className={classes}>
        {children}
      </Tag>
    );
  }
);

function constructClasses(size: Size, weight: Weight) {
  return `${cl.title} ${cl[size]} ${cl[weight]}`;
}

export default Title;
