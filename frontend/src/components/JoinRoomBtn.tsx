import spade from "../assets/svgs/spade.svg";

interface JoinRoomBtnProps {
  children: string;
  onClick: () => void;
}

const JoinRoomBtn = ({ children, onClick }: JoinRoomBtnProps) => {
  return (
    <button className="button--start-session" onClick={onClick}>
      {children}{" "}
      <img
        src={spade}
        className="image--start-spade"
        alt="A spade is a spade"
      />
    </button>
  );
};

export default JoinRoomBtn;
