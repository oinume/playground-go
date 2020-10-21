// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package github

import (
	"context"
	"sync"
)

// Ensure, that ClientMock does implement Client.
// If this is not the case, regenerate this file with moq.
var _ Client = &ClientMock{}

// ClientMock is a mock implementation of Client.
//
//     func TestSomethingThatUsesClient(t *testing.T) {
//
//         // make and configure a mocked Client
//         mockedClient := &ClientMock{
//             ListBranchesFunc: func(ctx context.Context, owner string, repo string) ([]string, error) {
// 	               panic("mock out the ListBranches method")
//             },
//         }
//
//         // use mockedClient in code that requires Client
//         // and then make assertions.
//
//     }
type ClientMock struct {
	// ListBranchesFunc mocks the ListBranches method.
	ListBranchesFunc func(ctx context.Context, owner string, repo string) ([]string, error)

	// calls tracks calls to the methods.
	calls struct {
		// ListBranches holds details about calls to the ListBranches method.
		ListBranches []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Owner is the owner argument value.
			Owner string
			// Repo is the repo argument value.
			Repo string
		}
	}
	lockListBranches sync.RWMutex
}

// ListBranches calls ListBranchesFunc.
func (mock *ClientMock) ListBranches(ctx context.Context, owner string, repo string) ([]string, error) {
	if mock.ListBranchesFunc == nil {
		panic("ClientMock.ListBranchesFunc: method is nil but Client.ListBranches was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Owner string
		Repo  string
	}{
		Ctx:   ctx,
		Owner: owner,
		Repo:  repo,
	}
	mock.lockListBranches.Lock()
	mock.calls.ListBranches = append(mock.calls.ListBranches, callInfo)
	mock.lockListBranches.Unlock()
	return mock.ListBranchesFunc(ctx, owner, repo)
}

// ListBranchesCalls gets all the calls that were made to ListBranches.
// Check the length with:
//     len(mockedClient.ListBranchesCalls())
func (mock *ClientMock) ListBranchesCalls() []struct {
	Ctx   context.Context
	Owner string
	Repo  string
} {
	var calls []struct {
		Ctx   context.Context
		Owner string
		Repo  string
	}
	mock.lockListBranches.RLock()
	calls = mock.calls.ListBranches
	mock.lockListBranches.RUnlock()
	return calls
}