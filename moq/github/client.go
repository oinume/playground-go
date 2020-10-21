package github

//go:generate moq -out=client_test.moq.go . Client
//go:generate mockgen -destination=client_test.gomock.go -package=github . Client

import "context"

type Client interface {
	ListBranches(ctx context.Context, owner, repo string) ([]string, error)
}
