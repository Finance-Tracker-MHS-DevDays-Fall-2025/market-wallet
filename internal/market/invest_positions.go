package market_impl

import (
	"context"
	"errors"
	"fmt"

	investapi "github.com/russianinvestments/invest-api-go-sdk/investgo"
	cm "market-wallet/internal/generated/api-common"
	m_pb "market-wallet/internal/generated/api-market"

	"market-wallet/internal/utils"
)

func GetInvestmentPositions(ctx context.Context, req *m_pb.GetInvestmentPositionsRequest) (*m_pb.GetInvestmentPositionsResponse, error) {
	// Проверяем, что backend правильного типа
	if req.Backend == nil || req.Backend.Type != "TInvest" {
		return nil, errors.New("invalid backend type")
	}

	// Создаем клиент Tinkoff Invest API
	cfg := utils.DefaultConfig(req.Backend.Token)
	client, err := investapi.NewClient(ctx, cfg, utils.GetGlobalLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to create invest client: %w", err)
	}

	// Получаем сервис операций
	operationsService := client.NewOperationsServiceClient()

	accountID := "123"

	// Получаем позиции по счету
	positionsResp, err := operationsService.GetPositions(accountID)
	if err != nil {
		return nil, err
	}

	// Получаем список позиций по ценным бумагам
	securities := positionsResp.GetSecurities()

	// Создаем слайс для хранения информации о позициях
	positions := make([]*m_pb.InvestmentPosition, 0, len(securities))

	// Получаем сервис инструментов для получения дополнительной информации
	instrumentsService := client.NewInstrumentsServiceClient()

	for _, security := range securities {
		// Получаем информацию об инструменте по его uid
		instrumentResp, err := instrumentsService.InstrumentByUid(security.GetInstrumentUid())
		if err != nil {
			return nil, err
		}

		instrument := instrumentResp.GetInstrument()

		// Создаем позицию
		position := &m_pb.InvestmentPosition{
			Figi:     instrument.GetFigi(),
			Quantity: int32(security.GetBalance()),
			Price: &cm.Money{
				Amount:   1,
				Currency: "RUR",
			},
		}

		positions = append(positions, position)
	}

	// Формируем ответ
	response := &m_pb.GetInvestmentPositionsResponse{
		Positions: positions,
	}

	return response, nil
}
