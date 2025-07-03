import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

export default function NewsCard({ data }) {
  if (!data) return null;

  return (
    <View style={styles.card}>
      <Text style={styles.title}>NEWS BIAS</Text>
      <Text style={styles.text}>Sentiment: {data.sentiment}</Text>
      <Text style={styles.text}>Bias: {data.bias}</Text>
      <Text style={styles.text}>Summary: {data.summary}</Text>
      <Text style={styles.text}>Source: {data.source}</Text>
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
    color: '#f48fb1',
    fontWeight: 'bold',
    marginBottom: 10,
  },
  text: {
    color: '#ccc',
  },
});
