"use client";

import {
  createContext,
  useContext,
  Dispatch,
  SetStateAction,
  useState,
} from "react";
import { File } from "@/models/file";

interface ContextProps {
  data: File[];
  setData: Dispatch<SetStateAction<File[]>>;
}

const GlobalContext = createContext<ContextProps>({
  data: [],
  setData: () => {},
});

export const GlobalContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [data, setData] = useState<File[]>([]);

  return (
    <GlobalContext.Provider value={{ data, setData }}>
      {children}
    </GlobalContext.Provider>
  );
};

export const useGlobalContext = () => useContext(GlobalContext);
