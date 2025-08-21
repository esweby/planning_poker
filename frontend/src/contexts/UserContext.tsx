import React, { createContext, useContext, useEffect, useState } from "react";

export type Role = "developer" | "tester" | "po" | "";

type UserContextType = {
  name: string;
  role: Role;
  setName: (name: string) => void;
  setRole: (role: Role) => void;
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

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ name, role }));
  }, [name, role]);

  const setName = (n: string) => setNameState(n);
  const setRole = (r: Role) => setRoleState(r);

  const clear = () => {
    setNameState("");
    setRoleState("");
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch {
      /* ignore */
    }
  };

  return (
    <UserContext.Provider value={{ name, role, setName, setRole, clear }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error("useUser must be used within a UserProvider");
  return ctx;
}
