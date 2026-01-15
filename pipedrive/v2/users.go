package v2

import (
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type UsersService struct {
	client *Client
}

type ListUserFollowersOption interface {
	applyListUserFollowers(*listUserFollowersOptions)
}

type UserFollowersRequestOption interface {
	ListUserFollowersOption
}

type listUserFollowersOptions struct {
	params         genv2.GetUserFollowersParams
	requestOptions []pipedrive.RequestOption
}

type userFollowersRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o userFollowersRequestOptions) applyListUserFollowers(cfg *listUserFollowersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listUserFollowersOptionFunc func(*listUserFollowersOptions)

func (f listUserFollowersOptionFunc) applyListUserFollowers(cfg *listUserFollowersOptions) {
	f(cfg)
}

func WithUserFollowersRequestOptions(opts ...pipedrive.RequestOption) UserFollowersRequestOption {
	return userFollowersRequestOptions{requestOptions: opts}
}

func WithUserFollowersPageSize(limit int) ListUserFollowersOption {
	return listUserFollowersOptionFunc(func(cfg *listUserFollowersOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithUserFollowersCursor(cursor string) ListUserFollowersOption {
	return listUserFollowersOptionFunc(func(cfg *listUserFollowersOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newListUserFollowersOptions(opts []ListUserFollowersOption) listUserFollowersOptions {
	var cfg listUserFollowersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListUserFollowers(&cfg)
	}
	return cfg
}

func (s *UsersService) ListFollowers(ctx context.Context, id UserID, opts ...ListUserFollowersOption) ([]Follower, *string, error) {
	cfg := newListUserFollowersOptions(opts)
	return s.listFollowers(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *UsersService) ListFollowersPager(id UserID, opts ...ListUserFollowersOption) *pipedrive.CursorPager[Follower] {
	cfg := newListUserFollowersOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Follower, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listFollowers(ctx, id, params, cfg.requestOptions)
	})
}

func (s *UsersService) ForEachFollowers(ctx context.Context, id UserID, fn func(Follower) error, opts ...ListUserFollowersOption) error {
	return s.ListFollowersPager(id, opts...).ForEach(ctx, fn)
}

func (s *UsersService) listFollowers(ctx context.Context, id UserID, params genv2.GetUserFollowersParams, requestOptions []pipedrive.RequestOption) ([]Follower, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetUserFollowersWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Follower `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}
