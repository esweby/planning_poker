type User = {
  name: string;
  role: string;
  seed: string;
};

type Game = {
  title: string;
  revealed: boolean;
  status: "waiting" | "playing" | "finished";
  guesses: { [string]: { guessingOn: string; prediction: string } };
  voted: { [string]: boolean };
  votes: { [string]: string };
};

type Message = {
  id: string;
  owner: string;
  users: { string: User };
  game: Game;
};
