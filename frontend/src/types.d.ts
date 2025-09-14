type User = {
  username: string;
  role: string;
  seed: string;
};

type Game = {
  title: string;
  revealed: boolean;
  status: "waiting" | "playing" | "finished";
  guesses: { [k: string]: { guessingOn: string; prediction: string } };
  voted: { [k: string]: boolean };
  votes: { [k: string]: string };
};

type Message = {
  id: string;
  owner: string;
  users: { [k: string]: User };
  game: Game;
};
