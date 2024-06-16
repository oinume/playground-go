package github

import (
	"context"
)

type MockedClient struct {
	Client
}

func (mc *MockedClient) ListBranches(ctx context.Context, owner string, repo string) ([]string, error) {
	return []string{"main", "develop", "feature/a"}, nil
}
