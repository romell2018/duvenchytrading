import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';

export default function HomeScreen({ navigation }) {
  const [tradeIdeas, setTradeIdeas] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('http://192.168.1.100:8080/trade-ideas?symbol=EURUSD')  // Use your machine's local IP here
      .then(res => res.json())
      .then(data => {
        setTradeIdeas(data.ideas);
        setLoading(false);
      })
      .catch(() => {
        setError('Failed to load trade ideas.');
        setLoading(false);
      });
  }, []);

  const renderItem = ({ item, index }) => (
    <TouchableOpacity
      style={styles.card}
      onPress={() => navigation.navigate('Trade Details', { trade: item })}
    >
      <Text style={styles.symbol}>{item.symbol} - {item.direction}</Text>
      <Text>Entry: {item.entry.toFixed(5)} | SL: {item.stop_loss.toFixed(5)} | TP: {item.take_profit.toFixed(5)}</Text>
      <Text>Pivot: {item.pivot.toFixed(5)}</Text>
      <Text>RSI: {item.rsi.toFixed(2)} | MACD Hist: {item.macd_hist.toFixed(5)}</Text>
      <Text style={styles.valid}>{item.valid ? '🟢 Valid' : '⚠️ Expired'}</Text>
      <Text>{item.comment}</Text>
    </TouchableOpacity>
  );

  if (loading) {
    return <ActivityIndicator size="large" style={{ flex: 1, justifyContent: 'center' }} />;
  }

  if (error) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', padding: 20 }}>
        <Text style={{ color: 'red', fontSize: 16 }}>{error}</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.header}>📋 Trade Setups</Text>
      <FlatList
        data={tradeIdeas}
        keyExtractor={(item, index) => item.symbol + index}
        renderItem={renderItem}
        contentContainerStyle={{ paddingBottom: 20 }}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 15 },
  header: { fontSize: 24, marginBottom: 10 },
  card: {
    padding: 15,
    borderRadius: 10,
    backgroundColor: '#f2f2f2',
    marginBottom: 10,
  },
  symbol: { fontSize: 18, fontWeight: 'bold' },
  valid: { marginTop: 5, fontStyle: 'italic' },
});
