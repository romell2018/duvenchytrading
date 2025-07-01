import React from 'react';
import { View, StyleSheet, Dimensions, Platform } from 'react-native';
import { WebView } from 'react-native-webview';

const ChartView = ({ symbol }) => {
  const symbolMap = {
    EURUSD: 'FX:EURUSD',
  GBPUSD: 'FX:GBPUSD',
  USDJPY: 'FX:USDJPY',
  BTCUSD: 'BINANCE:BTCUSDT',
  ETHUSD: 'COINBASE:ETHUSD',
  LTCUSD: 'BINANCE:LTCUSDT',
    // Add more as needed
  };

  const mappedSymbol = symbolMap[symbol] || 'FX:GBPUSD';

  const chartUrl = `https://www.tradingview.com/widgetembed/?frameElementId=tv_${symbol}&symbol=${mappedSymbol}&interval=30&theme=dark&style=1&locale=en&studies=["RSI@tv-basicstudies","MACD@tv-basicstudies"]&hide_top_toolbar=false&enable_publishing=false`;


  const getHTML = (tradingSymbol) => `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="UTF-8">
      <style>
        html, body, #chart { height: 100%; margin: 0; padding: 0; }
      </style>
      <script src="https://s3.tradingview.com/tv.js"></script>
    </head>
    <body>
      <div id="chart"></div>
      <script>
        new TradingView.widget({
          container_id: "chart",
          width: "100%",
          height: "100%",
          symbol: "${tradingSymbol}",
          interval: "30",
          timezone: "Etc/UTC",
          theme: "dark",
          style: "1",
          locale: "en",
          enable_publishing: false,
          hide_legend: false,
          hide_top_toolbar: false,
          save_image: false,
          studies: ["RSI@tv-basicstudies", "MACD@tv-basicstudies"],
          support_host: "https://www.tradingview.com"
        });
      </script>
    </body>
    </html>
  `;

  return (
    <View style={styles.container}>
      {Platform.OS === 'web' ? (
        <iframe
          title="TradingView Chart"
          src={chartUrl}
          width="100%"
          height="500"
          style={{ border: 'none' }}
        />
      ) : (
        <WebView
          originWhitelist={['*']}
          source={{ html: getHTML(mappedSymbol) }}
          style={styles.webview}
          javaScriptEnabled
          domStorageEnabled
          startInLoadingState
          scalesPageToFit
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  webview: {
    height: Dimensions.get('window').height / 2,
    width: Dimensions.get('window').width,
  },
});

export default ChartView;
