import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

export default function SignalCard({ data }) {
  if (!data) return null;

  const signalColor = data.signal === 'bullish' ? '#81c784' :
                      data.signal === 'bearish' ? '#e57373' : '#ccc';

  return (
    <View style={styles.card}>
      <Text style={[styles.title, { color: signalColor }]}>SIGNAL</Text>
      <Text style={styles.text}>Signal: {data.signal}</Text>
      <Text style={styles.text}>Reason: {data.reason}</Text>
      <Text style={styles.text}>RSI: {data.rsi.toFixed(2)}</Text>
      <Text style={styles.text}>MACD Histogram: {data.macd_histogram.toFixed(6)}</Text>
      <Text style={styles.text}>Support: {data.support}</Text>
      <Text style={styles.text}>Resistance: {data.resistance}</Text>
      <Text style={styles.text}>Bollinger Squeeze: {data.bollinger_squeeze ? "Yes" : "No"}</Text>
      <Text style={styles.text}>Divergence: {data.divergence}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#1e1e1e',
    padding: 15,
    borderRadius: 10,
    marginBottom: 20,
  },
  title: {
    fontWeight: 'bold',
    marginBottom: 10,
    fontSize: 16,
  },
  text: {
    color: '#ccc',
  },
});
