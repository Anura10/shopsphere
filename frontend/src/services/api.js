import axios from "axios";

export const API_URLS = {
  AUTH: "http://localhost:8081",
  PRODUCT: "http://localhost:8082",
  INVENTORY: "http://localhost:8083",
  CART: "http://localhost:8084",
  ORDER: "http://localhost:8085",
};

export const authAPI = axios.create({
  baseURL: `${API_URLS.AUTH}/api/v1/auth`,
});

export const productAPI = axios.create({
  baseURL: `${API_URLS.PRODUCT}/api/v1`,
});

export const inventoryAPI = axios.create({
  baseURL: `${API_URLS.INVENTORY}/api/v1`,
});

export const cartAPI = axios.create({
  baseURL: `${API_URLS.CART}/api/v1`,
});

export const orderAPI = axios.create({
  baseURL: `${API_URLS.ORDER}/api/v1`,
});