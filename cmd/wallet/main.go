package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	cm "market-wallet/internal/generated/api-common"
	pb "market-wallet/internal/generated/api-wallet"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedWalletServiceServer
}

func (s *server) GetAccounts(_ context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts := []*pb.Account{
		&pb.Account{
			AccountId: "aid",
			UserId:    req.GetUserId(),
			Name:      "Т-Банк",
			Type:      cm.AccountType_ACCOUNT_TYPE_INVESTMENT,
			Balance: &cm.Money{
				Amount:   10,
				Currency: "RUR",
			},
			CreatedAt: timestamppb.New(time.Now()),
		},
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}

func (s *server) GetTransactions(_ context.Context, req *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	transactions := []*pb.Transaction{
		&pb.Transaction{
			AccountId: "aid",
			UserId:    req.GetUserId(),
			Type:      cm.TransactionType_TRANSACTION_TYPE_EXPENSE,
			Amount: &cm.Money{
				Amount:   50,
				Currency: "RUR",
			},
			Category:      "eda",
			FromAccountId: "??",
			ToAccountId:   "??",
			Date:          timestamppb.New(time.Now()),
			Description:   "Заказал питсу",
		},
		&pb.Transaction{
			AccountId: "aid",
			UserId:    req.GetUserId(),
			Type:      cm.TransactionType_TRANSACTION_TYPE_EXPENSE,
			Amount: &cm.Money{
				Amount:   99999,
				Currency: "RUR",
			},
			Category:      "it",
			FromAccountId: "??",
			ToAccountId:   "??",
			Date:          timestamppb.New(time.Now()),
			Description:   "Арендовал сервак для хакатона",
		},
	}
	return &pb.GetTransactionsResponse{Transactions: transactions}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
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
