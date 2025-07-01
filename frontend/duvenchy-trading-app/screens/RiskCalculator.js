import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet } from 'react-native';
import RNPickerSelect from 'react-native-picker-select';

const RiskCalculator = () => {
    const [accountSize, setAccountSize] = useState('');
    const [riskPercent, setRiskPercent] = useState('');
    const [slPoints, setSlPoints] = useState('');
    const [pointValue, setPointValue] = useState('');
    const [result, setResult] = useState(null);
    const [selectedSymbol, setSelectedSymbol] = useState('');


    const calculateRisk = () => {
        const riskDollars = parseFloat(accountSize) * (parseFloat(riskPercent) / 100);
        const contractSize = riskDollars / (parseFloat(slPoints) * parseFloat(pointValue));
        setResult({
            riskDollars: riskDollars.toFixed(2),
            contractSize: contractSize.toFixed(2),
        });
    };
    const symbolPointValues = {
        ES: 50,
        NQ: 20,
        YM: 5,
        RTY: 50,
        CL: 10,
        GC: 10,
        _6E: 12.5,
        _6A: 10,
        ZN: 15.625,
        ZB: 31.25,
        // Add more symbols from your list as needed
    };

    return (
        <View style={styles.container}>
            <Text style={styles.title}>🧮 Risk Calculator</Text>

            <TextInput
                style={styles.input}
                placeholder="Account Size ($)"
                keyboardType="numeric"
                onChangeText={setAccountSize}
            />
            <TextInput
                style={styles.input}
                placeholder="Risk %"
                keyboardType="numeric"
                onChangeText={setRiskPercent}
            />
            <TextInput
                style={styles.input}
                placeholder="Stop Loss (Points)"
                keyboardType="numeric"
                onChangeText={setSlPoints}
            />
            <TextInput
                style={styles.input}
                placeholder="Point Value per Contract (e.g. 1)"
                keyboardType="numeric"
                onChangeText={<RNPickerSelect
                    onValueChange={(value) => {
                        setSelectedSymbol(value);
                        setPointValue(symbolPointValues[value]?.toString() || '');
                    }}
                    items={Object.keys(symbolPointValues).map(sym => ({
                        label: sym,
                        value: sym,
                    }))}
                    placeholder={{ label: 'Select Symbol', value: '' }}
                />
                }
            />

            <Button title="Calculate" onPress={calculateRisk} />

            {result && (
                <View style={styles.result}>
                    <Text>Risking: ${result.riskDollars}</Text>
                    <Text>Recommended Contracts: {result.contractSize}</Text>
                </View>
            )}
        </View>
    );
};

const styles = StyleSheet.create({
    container: { flex: 1, padding: 20 },
    title: { fontSize: 24, marginBottom: 10 },
    input: { borderBottomWidth: 1, marginBottom: 10, fontSize: 18 },
    result: { marginTop: 20 },
});

export default RiskCalculator;
