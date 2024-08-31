import axios from 'axios';

const instance = axios.create({
  baseUrl: 'http://localhost:3000',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const analisis = async (value) => {
  const { data } = await instance.post('/analisis', { value });
  return data;
}