import type { ReactNode } from "react";
import cl from "./Container.module.css";

type Type = "div" | "header" | "main" | "section" | "ul" | "li";
type Display = "block" | "flex" | "grid";
type JustifyContent = "start" | "end" | "center" | "between";
type AlignItems = "center";

export interface ContainerProps {
  type: Type;
  display: Display;
  children: ReactNode | ReactNode[];
  justifyContent?: JustifyContent;
  alignItems?: AlignItems;
  margin?: string;
  padding?: string;
  className?: string;
  gap?: string;
}

const Container = (props: ContainerProps) => {
  const {
    type,
    display,
    children,
    justifyContent,
    alignItems,
    margin = "",
    padding = "",
    gap = "",
    className,
  } = props;

  let classes = `${cl[display]}`;
  if (display === "flex") {
    if (justifyContent) classes += ` ${cl[`justify-${justifyContent}`]}`;
    if (alignItems) classes += ` ${cl[`align-${alignItems}`]}`;
  }

  const Tag: any = `${type}`;

  if (className) classes += ` ${className}`;

  const styles: { [K: string]: string } = {};
  if (margin) styles.margin = margin;
  if (padding) styles.padding = padding;
  if (gap) styles.gap = gap;

  return (
    <Tag className={classes} style={styles}>
      {children}
    </Tag>
  );
};

export default Container;
