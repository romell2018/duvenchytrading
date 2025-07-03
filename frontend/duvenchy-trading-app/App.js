import React, { useEffect, useState } from 'react';
import { ScrollView, View, Text, ActivityIndicator } from 'react-native';
import SymbolPicker from './components/SymbolPicker';
import QuoteCard from './components/QuoteCard';
import SignalCard from './components/SignalCard';
import NewsCard from './components/NewsCard';

const BASE_URL = 'http://localhost:8080/api';

export default function App() {
  const [symbol, setSymbol] = useState('6E');
  const [quote, setQuote] = useState(null);
  const [signal, setSignal] = useState(null);
  const [news, setNews] = useState(null);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      setErrorMsg('');
      try {
        const q = await fetch(`${BASE_URL}/quote/${symbol}`).then(res => res.json());
        const s = await fetch(`${BASE_URL}/signal/${symbol}`).then(res => res.json());
        const n = await fetch(`${BASE_URL}/news/${symbol}`).then(res => res.json());

        if (q.error || s.error || n.error) {
          setErrorMsg(q.error || s.error || n.error);
        } else {
          setQuote(q);
          setSignal(s);
          setNews(n);
        }

        console.log("Quote:", q);
        console.log("Signal:", s);
        console.log("News:", n);
      } catch (err) {
        console.error("Error fetching data:", err);
        setErrorMsg('Network error or backend is down.');
      }
      setLoading(false);
    };

    load();
  }, [symbol]);

  return (
    <ScrollView style={{ backgroundColor: '#121212' }} contentContainerStyle={{ padding: 20 }}>
      <SymbolPicker symbol={symbol} setSymbol={setSymbol} />

      {loading && <ActivityIndicator size="large" color="#4fc3f7" style={{ marginTop: 30 }} />}

      {errorMsg ? (
        <Text style={{ color: 'red', marginTop: 20 }}>{errorMsg}</Text>
      ) : (
        !loading && (
          <>
            <QuoteCard data={quote} />
            <SignalCard data={signal} />
            <NewsCard data={news} />
          </>
        )
      )}
    </ScrollView>
  );
}
