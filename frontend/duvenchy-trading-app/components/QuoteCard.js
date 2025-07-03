import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

export default function QuoteCard({ data }) {
  if (!data) return null;

  return (
    <View style={styles.card}>
      <Text style={styles.title}>QUOTE</Text>
      <Text style={styles.text}>Price: {data.price}</Text>
      <Text style={styles.text}>Change: {data.change} ({data.percent_change}%)</Text>
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
    color: '#4fc3f7',
    fontWeight: 'bold',
    marginBottom: 10,
  },
  text: {
    color: '#ddd',
  },
});
