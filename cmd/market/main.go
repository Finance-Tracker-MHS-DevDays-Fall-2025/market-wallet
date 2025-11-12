package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "market-wallet/internal/generated/api-market"
)

/*
type MarketServiceServer interface {
    GetInvestmentPositions(context.Context, *GetInvestmentPositionsRequest) (*GetInvestmentPositionsResponse, error)
    GetSecurity(context.Context, *GetSecurityRequest) (*GetSecurityResponse, error)
    GetSecuritiesPrices(context.Context, *GetSecuritiesPricesRequest) (*GetSecuritiesPricesResponse, error)
    GetSecurityPayments(context.Context, *GetSecuritiesPaymentsRequest) (*GetSecuritiesPaymentsResponse, error)
    mustEmbedUnimplementedMarketServiceServer()
}
*/

type server struct {
    pb.UnimplementedMarketServiceServer
}


func (s *server) GetInvestmentPositions(context.Context, *pb.GetInvestmentPositionsRequest) (*pb.GetInvestmentPositionsResponse, error) {
    log.Printf("Received: GetInvestmentPositions")
    return nil, nil
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
