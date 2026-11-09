// Each name in a room is protected by a secret token the server hands out on
// first join. We keep it per room so a refresh or reconnect can reclaim the
// same name, and nobody else can join as us.

const keyFor = (roomNumber: string) => `planning-poker:session:${roomNumber}`;

export function getSessionToken(roomNumber: string, name: string): string {
  try {
    const raw = localStorage.getItem(keyFor(roomNumber));
    if (!raw) return "";
    const session = JSON.parse(raw);
    return session.name === name ? (session.token ?? "") : "";
  } catch {
    return "";
  }
}

export function setSessionToken(
  roomNumber: string,
  name: string,
  token: string
) {
  try {
    localStorage.setItem(keyFor(roomNumber), JSON.stringify({ name, token }));
  } catch {
    /* ignore: we just won't be able to reclaim the name after a refresh */
  }
}

// Asks the server whether this name (and token) may join the room. Returns
// an error message, or null if the join should succeed.
export async function checkJoin(
  apiUrl: string,
  roomNumber: string,
  name: string
): Promise<string | null> {
  const params = new URLSearchParams({
    username: name,
    token: getSessionToken(roomNumber, name),
  });

  try {
    const res = await fetch(`${apiUrl}/api/check/${roomNumber}?${params}`);
    if (res.ok) return null;
    const body = await res.json().catch(() => ({}));
    return body.error ?? "Could not join this room";
  } catch {
    return "Could not reach the server";
  }
}
