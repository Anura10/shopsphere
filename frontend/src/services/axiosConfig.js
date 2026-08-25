import {
  authAPI,
  productAPI,
  inventoryAPI,
  cartAPI,
  orderAPI,
} from "./api";

const clients = [
  authAPI,
  productAPI,
  inventoryAPI,
  cartAPI,
  orderAPI,
];

clients.forEach((client) => {
  client.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem("token");

      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
    },
    (error) => Promise.reject(error)
  );
});