# S Rank for FIIs

This is a json HTTP server, that scrapes "brazilian REITs" (known as "FIIs") metrics, and calculates s-rank (a quantitative ranking), this metric is a ranking that combines the cheaper FIIs that pays the most dividends.

Source of the idea of s-rank: https://clubedovalor.com.br/blog/melhores-fiis-s-rank/

Source of the data about FIIs: https://www.fundsexplorer.com.br/ranking

## Routes
- GET http://localhost:9999/etl: requests the source of metrics for FIIs, calculates s-rank, and saves to DB
- POST http://localhost:9999/stocks: returns stocks sorted by ascending s-rank (better first)
  - you can apply some filters for the returned stocks using this sample body: `{ "positivePBRatio": true, "greaterThanDailyLiquidityInCurrency": 200000 }`

## Run Server

```
make init # run only once

make run-watch
```
