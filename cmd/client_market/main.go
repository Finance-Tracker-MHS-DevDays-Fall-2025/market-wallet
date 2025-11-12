package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "market-wallet/internal/generated/api-market"
    cm "market-wallet/internal/generated/api-common"
)

const (
    address     = "localhost:8888"
)

func main() {
    // установка соединения с сервером
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewMarketServiceClient(conn)

    // ус тановим контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // запрс на получение пользователя
    ab := cm.AccountBackend {
        Type: "Tinek",
        AccountId: "Test-acc",
        Token: "Вот это безопасность, класс!",
    }
    r, err := c.GetInvestmentPositions(ctx, &pb.GetInvestmentPositionsRequest{UserId: "test", Backend: &ab, AccountId: "test-ac"})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("GetInvestmentPositions: %v", r)
}
