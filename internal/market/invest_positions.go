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
	if req.Backend == nil || req.Backend.GetType() != "TInvest" {
		return nil, errors.New("invalid backend type")
	}

	// Создаем клиент Tinkoff Invest API
	cfg := utils.DefaultConfig(req.Backend.Token)
	client, err := investapi.NewClient(ctx, cfg, utils.GetGlobalLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to create invest client: %w", err)
	}

	usersClient := client.NewUsersServiceClient()
	accounts, err := usersClient.GetAccounts(nil) // получаем все счета

	accountID := ""
	for _, v := range accounts.GetAccounts() {
		accountID = v.GetId()
	}

	// Получаем сервис операций
	operationsService := client.NewOperationsServiceClient()

	// Получаем позиции по счету
	positionsResp, err := operationsService.GetPositions(accountID)
	if err != nil {
		return nil, err
	}

	// Получаем список позиций по ценным бумагам
	securities := positionsResp.GetSecurities()

	// Создаем слайс для хранения информации о позициях
	positions := make([]*m_pb.InvestmentPosition, 0, len(securities))

	figis := make([]string, 0, len(securities))
	for _, security := range securities {
		if err != nil {
			return nil, err
		}

		// Создаем позицию
		position := &m_pb.InvestmentPosition{
			Figi:     security.GetFigi(),
			Quantity: int32(security.GetBalance()),
			Price: &cm.Money{
				Amount:   0,
				Currency: "RUR",
			},
		}

		figis = append(figis, security.GetFigi())
		positions = append(positions, position)
	}

	infos, err := getInstrumentsInfoImpl(client, figis)
	if err != nil {
		return nil, err
	}
	for i, v := range infos {
		positions[i].Price.Amount = v.Price
	}

	// Формируем ответ
	response := &m_pb.GetInvestmentPositionsResponse{
		Positions: positions,
	}

	return response, nil
}
