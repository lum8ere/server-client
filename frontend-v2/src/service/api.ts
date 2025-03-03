import axios from 'axios';

// Здесь можно сконфигурировать базовый URL, интерцепторы и т.д.
const instance = axios.create({
    baseURL: 'http://localhost:9000' // или просто '/', если страница у нас на том же домене
});

export default instance;
