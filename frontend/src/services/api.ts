import axios from 'axios';

export const fetchServerTime = async (): Promise<string> => {
    const response = await axios.get('/api/time');
    return response.data;
};

export const fetchClientMetrics = async (clientId: string) => {
    const response = await axios.get(`/clientmetrics?client=${clientId}`);
    return response.data;
};

export const fetchClients = async () => {
    const response = await axios.get('/clients'); // Если такой эндпоинт есть на сервере
    return response.data;
};

export const sendCommand = async (clientId: string, cmd: string) => {
    const response = await axios.get(`/command?cmd=${cmd}&id=${clientId}`);
    return response.data;
};
