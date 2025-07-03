import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api';

export const fetchQuote = symbol =>
  fetch(`${BASE_URL}/quote/${symbol}`).then(res => res.json());

export const fetchSignal = symbol =>
  fetch(`${BASE_URL}/signal/${symbol}`).then(res => res.json());

export const fetchNews = symbol =>
  fetch(`${BASE_URL}/news/${symbol}`).then(res => res.json());
