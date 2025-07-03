import React from 'react';
import { View, Text, ScrollView, Pressable, StyleSheet } from 'react-native';

const supportedSymbols = [
  "6A", "6B", "6C", "6E", "6J", "6N", "6S", "CL", "E7", "ES", "GC",
  "HE", "HG", "HO", "LE", "M2K", "M6A", "M6E", "MCL", "MES", "MGC",
  "MNQ", "MYM", "NG", "NQ", "QI", "QM", "QO", "RB", "RTY", "SI", "UB",
  "YM", "ZB", "ZC", "ZF", "ZL", "ZN", "ZQ", "ZS", "ZT", "ZW"
];

export default function SymbolPicker({ symbol, setSymbol }) {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Select Symbol:</Text>
      <ScrollView horizontal showsHorizontalScrollIndicator={false}>
        {supportedSymbols.map(sym => (
          <Pressable
            key={sym}
            onPress={() => setSymbol(sym)}
            style={[
              styles.button,
              { backgroundColor: symbol === sym ? '#4fc3f7' : '#333' }
            ]}
          >
            <Text style={{ color: symbol === sym ? '#000' : '#ccc' }}>{sym}</Text>
          </Pressable>
        ))}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: 20
  },
  title: {
    color: '#fff',
    marginBottom: 10,
  },
  button: {
    paddingVertical: 8,
    paddingHorizontal: 12,
    borderRadius: 8,
    marginRight: 10,
  },
});
