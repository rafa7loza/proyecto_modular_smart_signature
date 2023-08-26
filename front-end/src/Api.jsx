import { REST_API } from "./endpoint_urls";

import axios from 'axios';

const instance = axios.create({
  baseURL: `${REST_API.base_url}:${REST_API.port}/`
});

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default instance;