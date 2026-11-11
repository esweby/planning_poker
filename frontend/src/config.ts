// Backend address. Set VITE_API_URL (e.g. in .env.local) to point elsewhere.
export const API_URL: string =
  import.meta.env.VITE_API_URL ?? "http://localhost:8080";

// The websocket address is the same host with ws:// or wss://.
export const WS_URL: string = API_URL.replace(/^http/, "ws");
