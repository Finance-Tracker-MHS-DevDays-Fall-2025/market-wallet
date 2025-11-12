package main

import (
	"context"
	"log"
	//"os"
	"time"

	"google.golang.org/grpc"
	cm "market-wallet/internal/generated/api-common"
	pb "market-wallet/internal/generated/api-market"
	//"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:50051"
)

/*
const (
    YDEX_figi     = "TCS00A107T19"
    LUKOIL_figi   = "BBG004731032"
    SBER_figi     = "BBG0047315Y7"
    TINKOF_figi   = "TCS80A107UL4"
    GASPROM_figi  = "BBG004730RP0"
    TAT_figi      = "BBG004RVFFC0"
    BASHNEFT_figi = "BBG004S68758"
    ROSNEFT_figi  = "BBG004731354"
    Aeroflot_figi = "BBG004S683W7"
    MTS_figi      = "BBG004S681W1"
    OZON_farm_figi= "TCS00A109B25"
    Samolet_figi  = "BBG00F6NKQX3"
    LSR_figi      = "BBG004S68C39"
)
*/

func main() {
	// установка соединения с сервером
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketServiceClient(conn)

	// ус тановим контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// запрс на получение пользователя
	/*
		    start := timestamppb.New(time.Now().AddDate(0, -6, 0))
			stop := timestamppb.New(time.Now().AddDate(0, 6, 0)) // +6 месяцев
			r, err := c.GetSecurityPayments(ctx, &pb.GetSecuritiesPaymentsRequest{Figis: []string{os.Args[1]}, StartDate: start, EndDate: stop})
	*/
	r, err := c.GetInvestmentPositions(ctx, &pb.GetInvestmentPositionsRequest{UserId: "uid", Backend: &cm.AccountBackend{Type: "TInvest", AccountId: "aid", Token: "tok"}, AccountId: "aid"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("GetInvestmentPositions: %v", r)
}
