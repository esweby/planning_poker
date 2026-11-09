import { useRef, useState } from "react";
import { API_URL, WS_URL } from "../config";
import { checkJoin, getSessionToken, setSessionToken } from "../session";

export function useWebSocket() {
  const socketRef = useRef<WebSocket | null>(null);
  // Bumped by closeConnection so a connect still waiting on its pre-check
  // knows it was cancelled (e.g. React StrictMode's mount/unmount/mount).
  const attemptRef = useRef(0);
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<any[]>([]);
  const [error, setError] = useState<string | null>(null);

  const connect = async (
    roomNumber: string,
    name: string,
    role: string,
    seed: string
  ) => {
    const attempt = ++attemptRef.current;
    setError(null);

    // A rejected websocket handshake gives the browser no reason, so ask
    // first and show the server's explanation if the join would fail.
    const joinError = await checkJoin(API_URL, roomNumber, name);
    if (attempt !== attemptRef.current) return;
    if (joinError) {
      setError(joinError);
      return;
    }

    const params = new URLSearchParams({
      username: name,
      role,
      seed,
      token: getSessionToken(roomNumber, name),
    });
    const socket = new WebSocket(`${WS_URL}/api/join/${roomNumber}?${params}`);
    socketRef.current = socket;

    socket.onopen = () => setIsConnected(true);
    socket.onclose = (event) => {
      if (socketRef.current !== socket) return;
      setIsConnected(false);
      if (event.reason) setError(event.reason);
    };

    socket.onmessage = (event) => {
      let data;
      try {
        data = JSON.parse(event.data);
      } catch {
        console.error("Unreadable message from server", event.data);
        return;
      }

      if (data.type === "session") {
        setSessionToken(roomNumber, name, data.token);
      } else if (data.error) {
        setError(data.error);
      } else if (data.payload) {
        setError(null);
        setMessages((prev) => [data, ...prev]);
      }
    };
  };

  const closeConnection = () => {
    attemptRef.current++;
    socketRef.current?.close();
    socketRef.current = null;
  };

  const sendMessage = (msg: any) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(msg));
    }
  };

  return { isConnected, messages, error, connect, closeConnection, sendMessage };
}
