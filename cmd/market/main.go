package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "market-wallet/internal/generated/api-market"
	market "market-wallet/internal/market"
	"net"
)

type server struct {
	pb.UnimplementedMarketServiceServer
}

func (s *server) GetInvestmentPositions(c context.Context, req *pb.GetInvestmentPositionsRequest) (*pb.GetInvestmentPositionsResponse, error) {
	return market.GetInvestmentPositions(c, req)
}

func (s *server) GetSecurity(c context.Context, req *pb.GetSecurityRequest) (*pb.GetSecurityResponse, error) {
	infos, err := market.GetInstrumentsInfo(c, []string{req.GetFigi()})
	if err != nil {
		return nil, err
	}
	return &pb.GetSecurityResponse{Security: infos[0]}, nil
}

func (s *server) GetSecuritiesPrices(c context.Context, req *pb.GetSecuritiesPricesRequest) (*pb.GetSecuritiesPricesResponse, error) {
	infos, err := market.GetInstrumentsInfo(c, req.GetFigis())
	if err != nil {
		return nil, err
	}
	return &pb.GetSecuritiesPricesResponse{Securities: infos}, nil
}

func (s *server) GetSecurityPayments(c context.Context, req *pb.GetSecuritiesPaymentsRequest) (*pb.GetSecuritiesPaymentsResponse, error) {
	future_payments, err := market.GetFuturePayments(c, req.GetFigis(), req.GetStartDate(), req.GetEndDate())
	if err != nil {
		return nil, err
	}
	return &pb.GetSecuritiesPaymentsResponse{Payments: future_payments}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMarketServiceServer(s, &server{})
	log.Printf("MarketServiceServe listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
