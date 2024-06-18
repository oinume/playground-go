package github

import "context"

//go:generate go run github.com/matryer/moq -out=client.moq.go . Client
//go:generate go run go.uber.org/mock/mockgen -destination=client.gomock.go -package=github . Client

type Client interface {
	ListBranches(ctx context.Context, owner, repo string) ([]string, error)
	Foo() error
	Bar() error
}

type MockedClient struct {
	Client
}

func (m *MockedClient) ListBranches(
	ctx context.Context,
	owner, repo string,
) ([]string, error) {
	return []string{"main", "develop", "feature/a"}, nil
}

// Foo, Barは実装しなくてOK
