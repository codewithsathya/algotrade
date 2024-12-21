package algo

import (
	"algotrade/pkg/binance"
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"
)

type ArbitrageDetector struct {
	marketDetails *binance.MarketDetails
	ctx           context.Context
	startAsset    string
	asks          map[string][][2]string
	bids          map[string][][2]string
    currencyIndexMap map[string]int
    edges           []Edge
}

func NewArbitrageDetector(ctx context.Context, marketDetails *binance.MarketDetails) *ArbitrageDetector {
	return &ArbitrageDetector{
		ctx:           ctx,
		marketDetails: marketDetails,
		startAsset:    "USDT",
        currencyIndexMap: make(map[string]int),
	}
}

func (a *ArbitrageDetector) Start() {
	t := time.NewTicker(1 * time.Second)
	time.Sleep(1 * time.Second)
	for {
		select {
		case <-t.C:
			err := a.SetEdges()
			if err != nil {
				continue
			}
            if(len(a.edges) < 2 * len(a.marketDetails.GetProducts())) {
                continue
            }
			negativeCycles, err := a.DetectArbitrage()
			if err != nil {
				continue
			}
            if len(negativeCycles) == 0 {
                fmt.Println("No negative cycles detected")
                continue
            }
            for _, cycle := range negativeCycles {
                profit, err := a.CalculateProfit(cycle)
                if err != nil {
                    continue
                }
                fmt.Println("Negative cycle detected")
                fmt.Printf("profit: %v\n", profit)
                for _, edge := range cycle {
                    edgeType := "SELL " + edge.from + " for " + edge.to
                    if edge.isBase {
                        edgeType = "BUY " + edge.to + " with " + edge.from
                    }
                    fmt.Printf("%v -> %v %v %v\n", edge.from, edge.to, edgeType, edge.rate)
                }
                fmt.Println()
            }
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *ArbitrageDetector) SetEdges() error {
	a.asks, a.bids = a.marketDetails.GetAsksAndBids()
	products := a.marketDetails.GetProducts()
	for _, product := range products {
		ask := a.asks[product.Symbol]
		if ask == nil {
			continue
		}
		askValue, err := strconv.ParseFloat(ask[0][0], 64)
		if err != nil {
			return err
		}
		bid := a.bids[product.Symbol]
		bidValue, err := strconv.ParseFloat(bid[0][0], 64)
		if err != nil {
			return err
		}
        if a.currencyIndexMap[product.QuoteAsset] == 0 {
            a.currencyIndexMap[product.QuoteAsset] = len(a.currencyIndexMap)
        }
        if a.currencyIndexMap[product.BaseAsset] == 0 {
            a.currencyIndexMap[product.BaseAsset] = len(a.currencyIndexMap)
        }
		a.edges = append(a.edges, Edge{
			from:   product.QuoteAsset,
			to:     product.BaseAsset,
			weight: -math.Log10(askValue),
            rate:  askValue,
            isBase: true,
		})
		a.edges = append(a.edges, Edge{
			from:   product.BaseAsset,
			to:     product.QuoteAsset,
			weight: -math.Log10((1 / bidValue)),
            rate: 1 / bidValue,
            isBase: false,
		})
	}
    return nil
}

func (a *ArbitrageDetector) GetCurrencyIndex(currency string) int {
    return a.currencyIndexMap[currency]
}

func (a *ArbitrageDetector) DetectArbitrage() ([][]Edge, error) {
    numberOfCurrencies := len(a.currencyIndexMap)
    distances := make([]float64, numberOfCurrencies)
    predecessors := make([]Edge, numberOfCurrencies)
    visitedEdges := make(map[string]bool)
    
    for i := range distances {
        distances[i] = math.Inf(1)
    }
    
    startIndex := a.GetCurrencyIndex(a.startAsset)
    distances[startIndex] = 0.0

    for i := 0; i < numberOfCurrencies; i++ {
        for _, edge := range a.edges {
            u := a.GetCurrencyIndex(edge.from)
            v := a.GetCurrencyIndex(edge.to)
            if !math.IsInf(distances[u], 1) && distances[u]+edge.weight < distances[v] {
                distances[v] = distances[u] + edge.weight
                predecessors[v] = edge
            }
        }
    }

    cycles := [][]Edge{}
    for _, edge := range a.edges {
        u := a.GetCurrencyIndex(edge.from)
        v := a.GetCurrencyIndex(edge.to)
        
        if visitedEdges[edge.from+edge.to] {
            continue
        }
        
        if !math.IsInf(distances[u], 1) && distances[u]+edge.weight < distances[v] {
            cycle := a.ConstructCycle(v, predecessors)
            if len(cycle) > 0 {
                // Mark all edges in the cycle as visited
                for _, cycleEdge := range cycle {
                    visitedEdges[cycleEdge.from+cycleEdge.to] = true
                }
                cycles = append(cycles, cycle)
            }
        }
    }
    return cycles, nil
}

func (a *ArbitrageDetector) ConstructCycle(start int, predecessors []Edge) []Edge {
    cycle := []Edge{}
    visited := make(map[int]bool)
    curr := start
    var cycleStart *int
    for {
        if visited[curr] {
            cycleStart = &curr
            break
        }
        visited[curr] = true
        if curr < 0 || curr >= len(predecessors) || predecessors[curr].from == "" {
            break
        }
        curr = a.GetCurrencyIndex(predecessors[curr].from)
    }
    if cycleStart != nil {
        cycle = []Edge{}
        visited = make(map[int]bool)
        curr = *cycleStart
        for {
            if visited[curr] {
                break
            }
            visited[curr] = true
            if curr < 0 || curr >= len(predecessors) || predecessors[curr].from == "" {
				break
			}
            cycle = append(cycle, predecessors[curr])
            curr = a.GetCurrencyIndex(predecessors[curr].from)
        }
    }
    if len(cycle) <= 2 {
        cycle = []Edge{}
    }
    slices.Reverse(cycle)
    return cycle
}

func (a *ArbitrageDetector) CalculateProfit(cycle []Edge) (float64, error) {
    totalRate := 1.0
    
    for _, edge := range cycle {
        rate := 1 / edge.rate
        totalRate *= rate
    }
    profitPercentage := (totalRate - 1.0) * 100
    
    return profitPercentage, nil
}