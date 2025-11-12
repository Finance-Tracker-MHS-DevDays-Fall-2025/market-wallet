package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "market-wallet/internal/generated/api-market"
    cm "market-wallet/internal/generated/api-common"

    // for created_at
    "google.golang.org/protobuf/types/known/timestamppb"
    "time"
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

// Мок для получения инфы о количестве бумаг у конкретного пользователя
// in (user_id, backend, account_id)
func (s *server) GetInvestmentPositions(context.Context, *pb.GetInvestmentPositionsRequest) (*pb.GetInvestmentPositionsResponse, error) {
    log.Printf("Received: GetInvestmentPositions")
    positions := []*pb.InvestmentPosition{
        &pb.InvestmentPosition {
            Figi: "figi1",
            Quantity: 15,
            Price: &cm.Money {
                Amount: 100000, // копеек
                Currency: "RUR",
            },
        },
        &pb.InvestmentPosition{
            Figi: "figi2",
            Quantity: 1,
            Price: &cm.Money {
                Amount: 130000, // копеек
                Currency: "RUR", 
            },
        },
    }
    ret := &pb.GetInvestmentPositionsResponse{Positions: positions};
    return ret, nil
}

// Мок для получения информации о бумаге/облигации
// in figi
// out (id, figi, pretty_name, current_price, price_updated_at)
func (s *server) GetSecurity(_ context.Context, req *pb.GetSecurityRequest) (*pb.GetSecurityResponse, error) {
    log.Printf("Received: GetSecurity")
    sec := &pb.Security{
        Id: "хз что тут должно быть",
        Figi: req.GetFigi(),
        Name: "Полное имя бумаги",
        CurrentPrice: &cm.Money {
            Amount: 90000, // копеек
            Currency: "RUR",
        },
        PriceUpdatedAt: timestamppb.New(time.Now()),
    }
    ret := &pb.GetSecurityResponse{Security: sec};
    return ret, nil
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
