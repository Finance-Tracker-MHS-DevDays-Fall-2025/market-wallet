package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "market-wallet/internal/generated/api-wallet"
    cm "market-wallet/internal/generated/api-common"

    // for created_at
    "google.golang.org/protobuf/types/known/timestamppb"
    "time"
)


/*
type WalletServiceServer interface {
    GetAccounts(context.Context, *GetAccountsRequest) (*GetAccountsResponse, error)
    GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error)
    mustEmbedUnimplementedWalletServiceServer()
}
*/


type server struct {
    pb.UnimplementedWalletServiceServer
}

// мок для получения аккаунтов по пользователю
// in (user_id, [backends]+)
// out (account_id, user_id, name, AccountType, total Money, created_at)+
func (s *server) GetAccounts(_ context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
    accounts := []*pb.Account{
        &pb.Account{
            AccountId: "хз че тут должно быть",
            UserId: req.GetUserId(),
            Name: "Т-Банк",
            Type: cm.AccountType_ACCOUNT_TYPE_INVESTMENT,
            Balance: &cm.Money {
                Amount: 10,
                Currency: "RUR",
            },
            CreatedAt: timestamppb.New(time.Now()), // bruh^3 
        },
    }
    return &pb.GetAccountsResponse{Accounts: accounts}, nil
}

func (s *server) GetTransactions(context.Context, *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
    
    return nil, nil
}


func main() {
    lis, err := net.Listen("tcp", ":9999")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterWalletServiceServer(s, &server{})
    log.Printf("WalletServiceServer listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
