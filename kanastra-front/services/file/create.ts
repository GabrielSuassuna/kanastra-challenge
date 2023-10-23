"use client";

import api from "../api";
import { ApiResponse } from "../types";
import { File } from "@/models/file";

const create = async (data: FormData) => {
  const headers = {
    "Content-Type": "multipart/form-data",
    accept: "application/json",
  };

  return api.post<ApiResponse<File>>("/file", data, {
    headers,
  });
};

export default create;
