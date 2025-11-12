package market_impl

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"

	cm "market-wallet/internal/generated/api-common"
	m_pb "market-wallet/internal/generated/api-market"

	"market-wallet/internal/utils"
)

// PaymentInfo содержит информацию о выплате
type PaymentInfo struct {
	FIGI   string
	Amount int64
	Date   *timestamppb.Timestamp
}

func GetFuturePayments(ctx context.Context, figis []string, start_date, stop_date *timestamppb.Timestamp) ([]*m_pb.SecurityPayment, error) {
	cfg := utils.DefaultConfig("TOKEN")
	client, err := investgo.NewClient(ctx, cfg, utils.GetGlobalLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to create invest client: %w", err)
	}
	defer client.Stop()

	list, err := getFuturePayments(client, figis, start_date, stop_date)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment calendar: %w", err)
	}

	ret := make([]*m_pb.SecurityPayment, 0, len(figis))
	for _, v := range list {
		ret = append(ret, &m_pb.SecurityPayment{
			Figi: v.FIGI,
			Payment: &cm.Money{
				Amount:   v.Amount,
				Currency: "RUR",
			},
			PaymentDate: v.Date,
		})
	}

	return ret, nil
}

// GetFuturePayments возвращает список будущих выплат для указанных инструментов
func getFuturePayments(
	client *investgo.Client,
	figis []string,
	start *timestamppb.Timestamp,
	stop *timestamppb.Timestamp,
) ([]PaymentInfo, error) {
	instrumentsService := client.NewInstrumentsServiceClient()
	var payments []PaymentInfo

	for _, figi := range figis {
		// Получаем информацию об инструменте по figi
		instrumentResp, err := instrumentsService.InstrumentByFigi(figi)
		if err != nil {
			continue
		}

		instrument := instrumentResp.GetInstrument()
		uid := instrument.GetUid()

		switch instrument.GetInstrumentKind() {
		case pb.InstrumentType_INSTRUMENT_TYPE_SHARE:
			// Для акций получаем дивиденды
			dividendsResp, err := instrumentsService.GetDividents(uid, start.AsTime(), stop.AsTime())
			if err != nil {
				return nil, err
			}

			for _, div := range dividendsResp.GetDividends() {
				payments = append(payments, PaymentInfo{
					FIGI:   figi,
					Amount: utils.ToRUR(div.GetDividendNet()),
					Date:   div.GetPaymentDate(),
				})
			}

		case pb.InstrumentType_INSTRUMENT_TYPE_BOND:
			// Для облигаций получаем купоны
			couponsResp, err := instrumentsService.GetBondCoupons(uid, start.AsTime(), stop.AsTime())
			if err != nil {
				return nil, err
			}

			for _, coupon := range couponsResp.GetEvents() {
				payments = append(payments, PaymentInfo{
					FIGI:   figi,
					Amount: utils.ToRUR(coupon.GetPayOneBond()),
					Date:   coupon.GetCouponDate(),
				})
			}
		}
	}

	return payments, nil
}
