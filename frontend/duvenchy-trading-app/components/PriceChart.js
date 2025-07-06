import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { LineChart, Grid } from 'react-native-svg-charts';

export default function PriceChart({ data }) {
  if (!data || data.length < 2) return null;

  return (
    <View style={styles.container}>
      <Text style={styles.title}>PRICE CHART (last 100 candles)</Text>
      <LineChart
        style={{ height: 200 }}
        data={data}
        svg={{ stroke: '#4fc3f7' }}
        contentInset={{ top: 20, bottom: 20 }}
      >
        <Grid />
      </LineChart>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
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
});
