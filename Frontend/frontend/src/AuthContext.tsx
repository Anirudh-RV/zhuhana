import React, { createContext, useContext, useEffect, useState } from "react";

interface User {
  ID: string;
  FirstName: string;
  MiddleName: string;
  LastName: string;
  EmailID: string;
  CreatedAt: string;
  UpdatedAt: string;
}

interface AuthContextType {
  user: User | null;
  accessToken: string | null;
  setAuth: (user: User, token: string) => void;
  clearAuth: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    const storedToken = localStorage.getItem("accessToken");
    if (storedUser && storedToken) {
      setUser(JSON.parse(storedUser));
      setAccessToken(storedToken);
    }
  }, []);

  const setAuth = (user: User, token: string) => {
    localStorage.setItem("user", JSON.stringify(user));
    localStorage.setItem("accessToken", token);
    setUser(user);
    setAccessToken(token);
  };

  const clearAuth = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("accessToken");
    setUser(null);
    setAccessToken(null);
  };

  return (
    <AuthContext.Provider value={{ user, accessToken, setAuth, clearAuth }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};
