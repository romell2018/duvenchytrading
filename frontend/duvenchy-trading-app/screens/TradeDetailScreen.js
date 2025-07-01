import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import ChartView from '../components/ChartView';

export default function TradeDetailScreen({ route }) {
  const { trade } = route.params;

  return (
    
    <View style={{ flex: 1 }}>
      <ChartView symbol={trade.symbol} />
      <View style={styles.details}>
        <Text style={styles.heading}>{trade.symbol} Trade Plan</Text>
        <Text>🧭 Direction: {trade.direction}</Text>
        <Text>🎯 Entry: {trade.entry}</Text>
        <Text>🛡️ Stop Loss: {trade.sl}</Text>
        <Text>🏁 Take Profit: {trade.tp}</Text>
        <Text>📉 RSI: {trade.rsi}</Text>
        <Text>📈 MACD: {trade.macd}</Text>
        <Text style={styles.valid}>{trade.valid ? '✅ Still Valid' : '❌ Invalidated'}</Text>
      </View>
    </View>
  );
  
}

const styles = StyleSheet.create({
  details: {
    padding: 15,
    backgroundColor: '#fff',
  },
  heading: {
    fontSize: 20,
    marginBottom: 8,
    fontWeight: 'bold',
  },
  valid: {
    marginTop: 10,
    fontStyle: 'italic',
  },
});
