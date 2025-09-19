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

  return (
    <Tag
      className={classes}
      style={{
        margin: `${margin.length > 0 ? margin : "0px"}`,
        padding: `${padding.length > 0 ? padding : "0px"}`,
        gap: `${gap.length > 0 ? gap : ""}`,
      }}
    >
      {children}
    </Tag>
  );
};

export default Container;
