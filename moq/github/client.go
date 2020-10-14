package github

//go:generate moq -out client_test.moq.go . Client

import "context"

type Client interface {
	ListBranches (ctx context.Context, owner, repo string) ([]string, error)
}
