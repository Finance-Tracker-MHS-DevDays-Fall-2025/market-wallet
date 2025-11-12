package market_impl


import (
	"context"
	"fmt"
	"log"
	"sync"

    "google.golang.org/protobuf/types/known/timestamppb"
    "time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"

    m_pb "market-wallet/internal/generated/api-market"
    cm "market-wallet/internal/generated/api-common"

)

// instrumentInfo содержит информацию о финансовом инструменте
type instrumentInfo struct {
	FIGI  string
	Name  string
	Price int64
	Error error
}


func GetInstrumentsInfo(ctx context.Context, figis []string) ([]*m_pb.Security, error) {
    cfg := DefaultConfig("TOKEN");
    client, err := investgo.NewClient(ctx, cfg, GetGlobalLogger())
    if err != nil {
        return nil, fmt.Errorf("failed to create invest client: %w", err)
    }
    defer client.Stop()
    
    list, err := getInstrumentsInfoImpl(client, figis)
    if err != nil {
        return nil, fmt.Errorf("failed to get figi infos: %w", err)
    }
    
    results := make([]*m_pb.Security, 0, len(figis))
    for _, v := range list {
        results = append(results, &m_pb.Security{
            Id: "GetInstrumentsInfo: хз что тут должно быть",
            Figi: v.FIGI,
            Name: v.Name,
            CurrentPrice: &cm.Money {
                Amount: v.Price,
                Currency: "RUR",
            },
            PriceUpdatedAt: timestamppb.New(time.Now()),
        }) 
    }
    return results, nil
}

func getInstrumentsInfoImpl(client *investgo.Client, figis []string) ([]instrumentInfo, error) {
	instrumentsService := client.NewInstrumentsServiceClient()
	marketDataService := client.NewMarketDataServiceClient()

	var wg sync.WaitGroup
	results := make([]instrumentInfo, len(figis))
	errors := make(chan error, len(figis))

	for i, figi := range figis {
		wg.Add(1)
		go func(index int, f string) {
			defer wg.Done()
			
			var info instrumentInfo
			info.FIGI = f
			
			// Получаем базовую информацию об инструменте
			resp, err := instrumentsService.InstrumentByFigi(f)
			if err != nil {
				info.Error = fmt.Errorf("failed to get instrument by FIGI: %v", err)
				results[index] = info
				errors <- info.Error
				return
			}
			
			// Сохраняем имя инструмента
			info.Name = resp.GetInstrument().GetName()
			
			// Получаем последнюю цену
			lastPriceResp, err := marketDataService.GetLastPrices([]string{f})
			if err != nil {
				info.Error = fmt.Errorf("failed to get last price: %v", err)
				results[index] = info
				errors <- info.Error
				return
			}
			
			if len(lastPriceResp.GetLastPrices()) == 0 {
				info.Error = fmt.Errorf("no price data available")
				results[index] = info
				errors <- info.Error
				return
			}
			
			// Сохраняем цену
			lastPrice := lastPriceResp.GetLastPrices()[0]
            pr := lastPrice.GetPrice()
			info.Price = int64(pr.Units * 10) + int64(pr.Nano / 1e7)
			
			results[index] = info
		}(i, figi)
	}

	wg.Wait()
	close(errors)
	
	for err := range errors {
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}	
	return results, nil
}
