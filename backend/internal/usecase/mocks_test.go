package usecase

import (
	"context"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockOrganizationRepository is a shared mock for OrganizationRepository
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) FindByID(ctx context.Context, id int64) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindAll(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindByParentID(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindRoots(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrganizationRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrganizationRepository) FindByPath(ctx context.Context, path string) ([]domain.Organization, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

// createTestOrganization is a shared helper for creating test organizations
func createTestOrganization(id int64, name string) *domain.Organization {
	now := time.Now()
	return &domain.Organization{
		ID:        id,
		Name:      name,
		ParentID:  nil,
		Path:      "/" + name + "/",
		Level:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
