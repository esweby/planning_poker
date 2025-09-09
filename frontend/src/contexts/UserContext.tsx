import React, { createContext, useContext, useEffect, useState } from "react";

export type Role = "developer" | "tester" | "po" | "";

type UserContextType = {
  name: string;
  role: Role;
  seed: string;
  setName: (name: string) => void;
  setRole: (role: Role) => void;
  setSeed: (see: string) => void;
  clear: () => void;
};

const STORAGE_KEY = "planning-poker:user";

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const [name, setNameState] = useState<string>(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      return raw ? (JSON.parse(raw).name ?? "") : "";
    } catch {
      return "";
    }
  });

  const [role, setRoleState] = useState<Role>(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      return raw ? (JSON.parse(raw).role ?? "") : "";
    } catch {
      return "";
    }
  });

  const [seed, setSeedState] = useState<string>(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      return raw ? (JSON.parse(raw).seed ?? "") : "";
    } catch {
      return "";
    }
  });

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ name, role, seed }));
  }, [name, role]);

  const setName = (n: string) => setNameState(n);
  const setRole = (r: Role) => setRoleState(r);
  const setSeed = (s: string) => setSeedState(s);

  const clear = () => {
    setNameState("");
    setRoleState("");
    setSeedState("");
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch {
      /* ignore */
    }
  };

  return (
    <UserContext.Provider
      value={{ name, role, seed, setName, setRole, setSeed, clear }}
    >
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error("useUser must be used within a UserProvider");
  return ctx;
}
