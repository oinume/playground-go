package main

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/oinume/playground-go/moq/github"
)

func TestService_PrintBranches(t *testing.T) {
	tests := map[string]struct {
		githubClient github.Client
		want         string
		wantErr      error
	}{
		"ok": {
			githubClient: &github.ClientMock{
				ListBranchesFunc: func(ctx context.Context, owner string, repo string) ([]string, error) {
					return []string{"main", "feature/xyz"}, nil
				},
			},
			want:    "main\nfeature/xyz\n",
			wantErr: nil,
		},
		"error": {
			githubClient: &github.ClientMock{
				ListBranchesFunc: func(ctx context.Context, owner string, repo string) ([]string, error) {
					return nil, errors.New("error")
				},
			},
			want:    "",
			wantErr: errors.New("error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := Service{githubClient: tt.githubClient}
			b := new(bytes.Buffer)
			if err := s.PrintBranches(context.Background(), b, "a", "b"); !reflect.DeepEqual(tt.wantErr, err) {
				t.Fatalf("unexpected error: err=%v, wantErr=%v", tt.wantErr, err)
			}
			if tt.wantErr == nil {
				if got := b.String(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unexpected result: got=%q, want=%q", got, tt.want)
				}
			}
		})
	}
}
