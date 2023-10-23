"use client";

import { MouseEvent, useCallback, useEffect, useRef, useState } from "react";
import { useGlobalContext } from "./context/store";
import { Button, Table } from "flowbite-react";
import createFile from "@/services/file/create";
import getFiles from "@/services/file/get-many";
import { File } from "@/models/file";

export default function Home() {
  const [loading, setLoading] = useState(false);
  const { data, setData } = useGlobalContext();

  const hiddenFileInput = useRef<HTMLInputElement>(null);

  const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
    console.log(hiddenFileInput);
    hiddenFileInput.current?.click();
  };

  const handleChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    setLoading(true);

    const file = event.target?.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file);
    const { data } = await createFile(formData);
    setData((prev) => [data, ...prev]);

    setLoading(false);

    if (hiddenFileInput.current) {
      hiddenFileInput.current.value = "";
    }
  };

  const handleGetFiles = useCallback(async () => {
    setLoading(true);
    const { data } = await getFiles();
    setData(data.reverse());
    setLoading(false);
  }, [setData]);

  const updateFile = useCallback(
    (file: File) => {
      setData((prev) => {
        const index = prev.findIndex((f) => f.id === file.id);
        console.log(index);
        if (index === -1) return prev;

        const newData = [...prev];
        newData[index] = file;
        return newData;
      });
    },
    [setData]
  );

  useEffect(() => {
    if (data.length <= 0 && !loading) {
      handleGetFiles();
    }
  }, [data, loading, handleGetFiles]);

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/events");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      try {
        const fileEvent = JSON.parse(event.data);
        console.log(fileEvent);
        if (fileEvent.message === "process-finished") {
          console.log("process-finished");
          updateFile(fileEvent.data);
        }
      } catch (error) {}
    };

    // terminating the connection on component unmount
    return () => eventSource.close();
  }, []);

  return (
    <div className="flex flex-col w-screen h-screen p-10 sm:p-20 gap-16 ">
      <Button
        className="w-full"
        gradientDuoTone="purpleToBlue"
        onClick={handleClick}
      >
        Adicionar arquivo
      </Button>
      <input
        type="file"
        onChange={handleChange}
        ref={hiddenFileInput}
        style={{ display: "none" }}
        accept=".csv"
      />

      <h2 className="text-2xl font-bold text-center mb-8">Arquivos</h2>
      <div className="overflow-x-auto">
        <Table>
          <Table.Head>
            <Table.HeadCell>Arquivo</Table.HeadCell>
            <Table.HeadCell>Data</Table.HeadCell>
            <Table.HeadCell>Processado</Table.HeadCell>
            <Table.HeadCell>Tempo</Table.HeadCell>
          </Table.Head>
          <Table.Body className="divide-y overflow-hidden">
            {data.map((file) => (
              <Table.Row
                key={file.id}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {file.name}
                </Table.Cell>
                <Table.Cell>
                  {new Date(file.created_at).toLocaleDateString()}
                </Table.Cell>
                <Table.Cell>
                  {file.processed ? (
                    <span className="text-green-500">Sim</span>
                  ) : (
                    <span className="text-red-500">Não</span>
                  )}
                </Table.Cell>
                <Table.Cell>
                  {file.processed
                    ? `${
                        Math.abs(
                          new Date(file.processed_at).getTime() -
                            new Date(file.created_at).getTime()
                        ) / 1000
                      } segundos`
                    : "-"}
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>
    </div>
  );
}
