package service

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Service interface {
	PostServiceInfo(ctx context.Context, h ServiceInfo) error
	GetServiceInfo(ctx context.Context, id string) (ServiceInfo, error)
	PutServiceInfo(ctx context.Context, id string, h ServiceInfo) error
	DeleteServiceInfo(ctx context.Context, id string) error
}

type ServiceInfo struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	HostID     string     `json:"hostid"`
	CreatedAt  time.Time  `json:"createtime"`
	UpdatedAt  time.Time  `json:"updatetime"`
	Remark     string     `json:"remark"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]ServiceInfo
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]ServiceInfo{},
	}
}

func (s *inmemService) PostServiceInfo(ctx context.Context, h ServiceInfo) error {
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

func (s *inmemService) GetServiceInfo(ctx context.Context, id string) (ServiceInfo, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	h, ok := s.m[id]
	if !ok {
		return ServiceInfo{}, ErrNotFound
	}
	return h, nil
}

func (s *inmemService) PutServiceInfo(ctx context.Context, id string, h ServiceInfo) error {
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

func (s *inmemService) DeleteServiceInfo(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

