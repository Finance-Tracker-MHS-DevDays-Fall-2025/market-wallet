package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	cm "market-wallet/internal/generated/api-common"
	pb "market-wallet/internal/generated/api-market"

	// for created_at
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	market "market-wallet/internal/market"
)

type server struct {
	pb.UnimplementedMarketServiceServer
}

// Мок для получения инфы о количестве бумаг у конкретного пользователя
// in (user_id, backend, account_id)
func (s *server) GetInvestmentPositions(c context.Context, req *pb.GetInvestmentPositionsRequest) (*pb.GetInvestmentPositionsResponse, error) {
	return market.GetInvestmentPositions(c, req)
}

// Мок для получения информации о бумаге/облигации
// in figi
// out (id, figi, pretty_name, current_price, price_updated_at)
func (s *server) GetSecurity(c context.Context, req *pb.GetSecurityRequest) (*pb.GetSecurityResponse, error) {
	infos, err := market.GetInstrumentsInfo(c, []string{req.GetFigi()})
	if err != nil {
		return nil, err
	}
	return &pb.GetSecurityResponse{Security: infos[0]}, nil
}

// Мок для получения информации о бумаге/облигации (теперь получаем массив
// in [figi]+
// out [(id, figi, pretty_name, current_price, price_updated_at)]+
func (s *server) GetSecuritiesPrices(c context.Context, req *pb.GetSecuritiesPricesRequest) (*pb.GetSecuritiesPricesResponse, error) {
	infos, err := market.GetInstrumentsInfo(c, req.GetFigis())
	if err != nil {
		return nil, err
	}
	return &pb.GetSecuritiesPricesResponse{Securities: infos}, nil
}

// Мок для получения выплат по бумагам
// in [figi]+, start_date, stop_date
// out [(figi,payment, payment_date), ...]+
func (s *server) GetSecurityPayments(_ context.Context, req *pb.GetSecuritiesPaymentsRequest) (*pb.GetSecuritiesPaymentsResponse, error) {
	log.Printf("Received: GetSecurityPayments")
	figis := req.GetFigis() // []string
	start_date := req.GetStartDate()
	end_date := req.GetEndDate()
	log.Printf("Received: GetSecurityPayments %v, %v, %v", figis, start_date, end_date)

	ret := []*pb.SecurityPayment{
		&pb.SecurityPayment{
			Figi: "figi1",
			Payment: &cm.Money{
				Amount:   34500, // копеек
				Currency: "RUR",
			},
			PaymentDate: timestamppb.New(time.Now()), // bruh
		},
		&pb.SecurityPayment{
			Figi: "figi2",
			Payment: &cm.Money{
				Amount:   100, // копеек
				Currency: "RUR",
			},
			PaymentDate: timestamppb.New(time.Now()), // bruh^2
		},
	}

	return &pb.GetSecuritiesPaymentsResponse{Payments: ret}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMarketServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
