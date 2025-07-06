import React, { useEffect, useState } from 'react';
import { ScrollView, View, Text, ActivityIndicator, Dimensions } from 'react-native';
import { LineChart } from 'react-native-chart-kit';
import SymbolPicker from './components/SymbolPicker';
import QuoteCard from './components/QuoteCard';
import SignalCard from './components/SignalCard';
import NewsCard from './components/NewsCard';

const BASE_URL = 'http://localhost:8080/api';
const screenWidth = Dimensions.get('window').width;

export default function App() {
  const [symbol, setSymbol] = useState('6E');
  const [quote, setQuote] = useState(null);
  const [signal, setSignal] = useState(null);
  const [news, setNews] = useState(null);
  const [chartData, setChartData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      setErrorMsg('');
      try {
        // Load candle data for chart
        const candleRes = await fetch(`https://api.twelvedata.com/time_series?symbol=${symbol}&interval=1h&outputsize=100&apikey=${process.env.EXPO_PUBLIC_TWELVE_API_KEY}`);
        const candleData = await candleRes.json();

        const closePrices = Array.isArray(candleData.values)
          ? candleData.values.map(c => parseFloat(c.close)).reverse()
          : [];

        setChartData(closePrices);

        // Load quote, signal, news
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
      } catch (err) {
        console.error("Error loading data:", err);
        setErrorMsg('Network error or backend is down.');
      }
      setLoading(false);
    };

    load();
  }, [symbol]);

  return (
    <ScrollView style={{ backgroundColor: '#121212' }} contentContainerStyle={{ padding: 20 }}>
      <SymbolPicker symbol={symbol} setSymbol={setSymbol} />

      {loading && (
        <ActivityIndicator size="large" color="#4fc3f7" style={{ marginTop: 30 }} />
      )}

      {errorMsg ? (
        <Text style={{ color: 'red', marginTop: 20 }}>{errorMsg}</Text>
      ) : (
        !loading && (
          <>
            {chartData.length > 1 && (
              <View style={{ marginBottom: 20 }}>
                <Text style={{ color: '#4fc3f7', fontWeight: 'bold', marginBottom: 10 }}>
                  PRICE CHART (last {chartData.length})
                </Text>
                <LineChart
                  data={{
                    labels: [],
                    datasets: [{ data: chartData }],
                  }}
                  width={screenWidth - 40}
                  height={220}
                  chartConfig={{
                    backgroundColor: '#121212',
                    backgroundGradientFrom: '#1e1e1e',
                    backgroundGradientTo: '#1e1e1e',
                    color: () => '#4fc3f7',
                    labelColor: () => '#ccc',
                    strokeWidth: 2,
                  }}
                  bezier
                  withDots={false}
                  withShadow={false}
                  style={{ borderRadius: 10 }}
                />
              </View>
            )}

            <QuoteCard data={quote} />
            <SignalCard data={signal} />
            <NewsCard data={news} />
          </>
        )
      )}
    </ScrollView>
  );
}
