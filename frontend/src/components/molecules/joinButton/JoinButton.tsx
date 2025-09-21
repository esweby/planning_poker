import Button from "../../atoms/button/Button";
import spade from "../../../assets/svgs/spade.svg";

import cl from "./JoinButton.module.css";

interface JoinButtonProps {
  onClick: () => void;
}

const JoinButton = ({ onClick }: JoinButtonProps) => {
  return (
    <Button onClick={onClick} className={cl.btn}>
      <>
        Join <img src={spade} className={cl.img} alt="a spade icon" />
      </>
    </Button>
  );
};

export default JoinButton;
