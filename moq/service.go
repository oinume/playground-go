package main

import (
	"context"
	"fmt"
	"io"

	"github.com/oinume/playground-go/moq/github"
)

type Service struct {
	githubClient github.Client
}

func (s *Service) PrintBranches(ctx context.Context, w io.Writer, owner, repo string) error {
	branches, err := s.githubClient.ListBranches(ctx, owner, repo)
	if err != nil {
		return err
	}
	for _, b := range branches {
		fmt.Fprintf(w, "%s\n", b)
	}
	return nil
}
