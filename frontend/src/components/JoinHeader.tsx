import { Link } from "@tanstack/react-router";

const JoinHeader = () => {
  return (
    <header className="index--header">
      <Link to="/" className="header--link">
        <h1>Planning Poker</h1>
      </Link>
    </header>
  );
};

export default JoinHeader;
