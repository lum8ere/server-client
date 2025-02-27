import axios from 'axios';

// Здесь можно сконфигурировать базовый URL, интерцепторы и т.д.
const instance = axios.create({
    baseURL: 'http://127.0.0.1:4000' // или просто '/', если страница у нас на том же домене
});

export default instance;
