import React, { useState, useEffect } from 'react';
import { View, ScrollView, ActivityIndicator } from 'react-native';
import SymbolPicker from '../components/SymbolPicker';
import QuoteCard from '../components/QuoteCard';
import SignalCard from '../components/SignalCard';
import NewsCard from '../components/NewsCard';
import { fetchQuote, fetchSignal, fetchNews } from '../utils/api';

export default function HomeScreen() {
  const [symbol, setSymbol] = useState('6E');
  const [quote, setQuote] = useState(null);
  const [signal, setSignal] = useState(null);
  const [news, setNews] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const [q, s, n] = await Promise.all([
          fetchQuote(symbol),
          fetchSignal(symbol),
          fetchNews(symbol),
        ]);
        setQuote(q.data);
        setSignal(s.data);
        setNews(n.data);
      } catch (err) {
        console.error(err);
      }
      setLoading(false);
    };
    load();
  }, [symbol]);

  return (
    <ScrollView contentContainerStyle={{ padding: 20 }}>
      <SymbolPicker symbol={symbol} setSymbol={setSymbol} />
      {loading ? (
        <ActivityIndicator size="large" />
      ) : (
        <>
          <QuoteCard data={quote} />
          <SignalCard data={signal} />
          <NewsCard data={news} />
        </>
      )}
    </ScrollView>
  );
}
