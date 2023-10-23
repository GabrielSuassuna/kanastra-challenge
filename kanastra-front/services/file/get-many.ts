import api from "../api";
import { ApiResponse } from "../types";
import { File } from "@/models/file";

const getMany = () => api.get<ApiResponse<File[]>>("/file");

export default getMany;
