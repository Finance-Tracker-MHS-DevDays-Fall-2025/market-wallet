package market_impl

import (
	"context"
	"errors"
	"fmt"

	investapi "github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "market-wallet/internal/generated/api-market"

	"market-wallet/internal/utils"
	//"github.com/russianinvestments/invest-api-go-sdk/retry"
)

func GetInvestmentPositions(ctx context.Context, req *pb.GetInvestmentPositionsRequest) (*pb.GetInvestmentPositionsResponse, error) {
	// Проверяем, что backend правильного типа
	if req.Backend == nil || req.Backend.Type != "TInvest" {
		return nil, errors.New("invalid backend type")
	}

	// Создаем клиент Tinkoff Invest API
	cfg := utils.DefaultConfig(req.Backend.Token)
	_, err := investapi.NewClient(ctx, cfg, utils.GetGlobalLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to create invest client: %w", err)
	}

	return nil, nil

	/*

	       // создаем клиента для сервиса операций
	   	operationsService := client.NewOperationsServiceClient()

	   	portfolioResp, err := operationsService.GetPortfolio(cfg.AccountId, pb.PortfolioRequest_RUB)
	   	if err != nil {
	   		global_logger.Errorf(err.Error())
	   	} else {
	   		fmt.Printf("amount of shares = %v\n", portfolioResp.GetTotalAmountShares())
	   	}
	*/

	/*

		// Если account_id не указан, берем первый счет
		accountID := req.AccountId
		if accountID == "" {
			accountID = accounts[0].Id
		}

		// Получаем позиции по счету
		positionsResp, err := retry.Retry(ctx, func() (*investapi.PositionsResponse, error) {
			return client.Operations.GetPositions(ctx, accountID)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get positions: %w", err)
		}

		// Конвертируем позиции в protobuf формат
		var positions []*pb.InvestmentPosition
		for _, security := range positionsResp.Securities {
			// Находим текущую цену инструмента
			lastPrice, err := client.MarketData.GetLastPrices(ctx, []string{security.Figi})
			if err != nil {
				return nil, fmt.Errorf("failed to get last price for figi %s: %w", security.Figi, err)
			}

			if len(lastPrice) == 0 {
				return nil, fmt.Errorf("no last price for figi %s", security.Figi)
			}

			positions = append(positions, &pb.InvestmentPosition{
				Figi:     security.Figi,
				Quantity: int32(security.Balance),
				Price: &pb.Money{
					Amount:   int64(lastPrice[0].Price * 100), // конвертируем в копейки
					Currency: "RUR",
				},
			})
		}

		return &pb.GetInvestmentPositionsResponse{
			Positions: positions,
		}, nil

	*/
}
