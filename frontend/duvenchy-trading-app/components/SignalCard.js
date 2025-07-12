import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

export default function SignalCard({ data }) {
  if (!data) return null;

  const signalColor =
    data.signal === 'bullish' ? '#81c784' :
    data.signal === 'bearish' ? '#e57373' : '#ccc';

  const confidenceColor =
    data.confidence === 'high' ? '#81c784' :
    data.confidence === 'medium' ? 'orange' : '#f44336';

  return (
    <View style={styles.card}>
      <Text style={[styles.title, { color: signalColor }]}>SIGNAL</Text>
      <Text style={styles.text}>Signal: {data.signal}</Text>
      <Text style={styles.text}>Reason: {data.reason}</Text>

      <Text style={[styles.sectionTitle]}>📈 Trade Setup</Text>
      <Text style={[styles.text, { color: '#81c784' }]}>Entry: {data.entry?.toFixed(5)}</Text>
      <Text style={[styles.text, { color: '#e57373' }]}>Stop Loss: {data.stop_loss?.toFixed(5)}</Text>
      <Text style={[styles.text, { color: '#64b5f6' }]}>Take Profit: {data.take_profit?.toFixed(5)}</Text>
      <Text style={styles.text}>Reward/Risk Ratio: {data.reward_risk_ratio?.toFixed(2)}</Text>

      <Text style={styles.text}>
        Confidence: <Text style={{ fontWeight: 'bold', color: confidenceColor }}>{data.confidence}</Text>
      </Text>

      <Text style={[styles.sectionTitle]}>🧠 GPT Summary</Text>
      <Text style={styles.text}>{data.final_trade_idea}</Text>

      <Text style={[styles.sectionTitle]}>🎯 Scenarios</Text>
      <Text style={styles.text}>Preferred: {data.preferred_scenario}</Text>
      <Text style={styles.text}>Alternative: {data.alternative_scenario}</Text>

      <Text style={[styles.sectionTitle]}>📊 Indicators</Text>
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
    fontSize: 18,
    marginBottom: 6,
  },
  sectionTitle: {
    marginTop: 12,
    marginBottom: 4,
    fontWeight: 'bold',
    color: '#eee',
    fontSize: 15,
  },
  text: {
    color: '#ccc',
    marginBottom: 2,
  },
});
ss