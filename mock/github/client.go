package github

//go:generate go run github.com/matryer/moq -out=client_test.moq.go . Client
//go:generate go run github.com/golang/mock/mockgen -destination=client_test.gomock.go -package=github . Client

import "context"

type Client interface {
	ListBranches(ctx context.Context, owner, repo string) ([]string, error)
}
