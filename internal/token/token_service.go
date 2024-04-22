package token

import (
	"context"
	"fmt"
	"github.com/instinctG/Test-task/internal/db"
)

type Store interface {
	PostRefreshToken(ctx context.Context, guid, refresh string) error
	ReadRefreshToken(ctx context.Context, guid string) (*db.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, guid, refresh string) error
}

type Service struct {
	Store Store
}

func NewTokenService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) PostRefreshToken(ctx context.Context, guid, refresh string) error {
	fmt.Println("posting a refresh token to database")
	err := s.Store.PostRefreshToken(ctx, guid, refresh)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) ReadRefreshToken(ctx context.Context, guid string) (*db.RefreshToken, error) {
	fmt.Println("Retrieving a refresh token")
	token, err := s.Store.ReadRefreshToken(ctx, guid)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Service) UpdateRefreshToken(ctx context.Context, guid, refresh string) error {
	err := s.Store.UpdateRefreshToken(ctx, guid, refresh)
	if err != nil {
		return err
	}
	return nil
}
