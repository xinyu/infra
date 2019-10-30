package host

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Host interface {
	PostHostInfo(ctx context.Context, h HostInfo) error
	GetHostInfo(ctx context.Context, id string) (HostInfo, error)
	PutHostInfo(ctx context.Context, id string, h HostInfo) error
	DeleteHostInfo(ctx context.Context, id string) error
}

type HostInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	IP          string     `json:"ip"`
	Port        string     `json:"port"`
	Rack        string     `json:"rack"`
	DataCenter  string     `json:"datacenter"`
	CreatedAt   time.Time  `json:"createtime"`
	UpdatedAt   time.Time  `json:"updatetime"`
	Remark      string     `json:"remark"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrNotFoundID      = errors.New("not found host ID")
)

type inmemHost struct {
	mtx sync.RWMutex
	m   map[string]HostInfo
}

func NewInmemHost() Host {
	return &inmemHost{
		m: map[string]HostInfo{},
	}
}

func (s *inmemHost) PostHostInfo(ctx context.Context, h HostInfo) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.m[h.ID]; ok {
		return ErrAlreadyExists
	}

	currentTime := time.Now()
	h.CreatedAt = currentTime
	h.UpdatedAt = currentTime

	s.m[h.ID] = h

	return nil
}

func (s *inmemHost) GetHostInfo(ctx context.Context, id string) (HostInfo, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	h, ok := s.m[id]
	if !ok {
		return HostInfo{}, ErrNotFound
	}
	return h, nil
}

func (s *inmemHost) PutHostInfo(ctx context.Context, id string, h HostInfo) error {
	if id != h.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()

	currentTime := time.Now()
	h.UpdatedAt = currentTime

	hLast, ok := s.m[id]
	if ok {
		h.CreatedAt = hLast.CreatedAt
	}

	s.m[id] = h
	return nil
}

func (s *inmemHost) DeleteHostInfo(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}