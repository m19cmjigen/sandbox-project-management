package usecase

import (
	"github.com/stretchr/testify/mock"
	"context"
	"errors"
	
	"testing"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// MockOrganizationRepositoryForUsecase extends MockOrganizationRepository with additional methods
type MockOrganizationRepositoryForUsecase struct {
	MockOrganizationRepository
}

func (m *MockOrganizationRepositoryForUsecase) FindRoots(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepositoryForUsecase) HasChildren(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func setupOrganizationUsecase(t *testing.T) (*organizationUsecase, *MockOrganizationRepositoryForUsecase) {
	mockRepo := new(MockOrganizationRepositoryForUsecase)
	usecase := &organizationUsecase{
		orgRepo: mockRepo,
	}
	return usecase, mockRepo
}


func createTestOrganizationWithPath(id int64, name string, parentID *int64, path string, level int) *domain.Organization {
	org := createTestOrganization(id, name)
	org.ParentID = parentID
	org.Path = path
	org.Level = level
	return org
}
// GetAll tests
func TestOrganizationUsecase_GetAll_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	expectedOrgs := []domain.Organization{
		*createTestOrganization(1, "Org1"),
		*createTestOrganization(2, "Org2"),
		*createTestOrganization(3, "Org3"),
	}

	mockRepo.On("FindAll", ctx).Return(expectedOrgs, nil)

	orgs, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, orgs, 3)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_GetAll_EmptyResult(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Organization{}, nil)

	orgs, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, orgs, 0)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_GetAll_RepositoryError(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return(nil, errors.New("database error"))

	orgs, err := uc.GetAll(ctx)

	assert.Nil(t, orgs)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByID tests
func TestOrganizationUsecase_GetByID_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	expectedOrg := createTestOrganization(1, "Test Org")
	mockRepo.On("FindByID", ctx, int64(1)).Return(expectedOrg, nil)

	org, err := uc.GetByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), org.ID)
	assert.Equal(t, "Test Org", org.Name)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_GetByID_NotFound(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	org, err := uc.GetByID(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, org)
	mockRepo.AssertExpectations(t)
}

// GetChildren tests
func TestOrganizationUsecase_GetChildren_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	parentID := int64(1)
	children := []domain.Organization{
		*createTestOrganizationWithPath(2, "Child1", &parentID, "/1/2/", 1),
		*createTestOrganizationWithPath(3, "Child2", &parentID, "/1/3/", 1),
	}

	mockRepo.On("FindByParentID", ctx, int64(1)).Return(children, nil)

	orgs, err := uc.GetChildren(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, orgs, 2)
	for _, org := range orgs {
		assert.Equal(t, int64(1), *org.ParentID)
	}
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_GetChildren_NoChildren(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByParentID", ctx, int64(1)).Return([]domain.Organization{}, nil)

	orgs, err := uc.GetChildren(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, orgs, 0)
	mockRepo.AssertExpectations(t)
}

// GetRoots tests
func TestOrganizationUsecase_GetRoots_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	roots := []domain.Organization{
		*createTestOrganizationWithPath(1, "Root1", nil, "/1/", 0),
		*createTestOrganizationWithPath(2, "Root2", nil, "/2/", 0),
	}

	mockRepo.On("FindRoots", ctx).Return(roots, nil)

	orgs, err := uc.GetRoots(ctx)

	require.NoError(t, err)
	assert.Len(t, orgs, 2)
	for _, org := range orgs {
		assert.Nil(t, org.ParentID)
		assert.Equal(t, 0, org.Level)
	}
	mockRepo.AssertExpectations(t)
}

// GetTree tests
func TestOrganizationUsecase_GetTree_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	parentID1 := int64(1)
	parentID2 := int64(2)

	allOrgs := []domain.Organization{
		*createTestOrganizationWithPath(1, "Root1", nil, "/1/", 0),
		*createTestOrganizationWithPath(2, "Root2", nil, "/2/", 0),
		*createTestOrganizationWithPath(3, "Child1-1", &parentID1, "/1/3/", 1),
		*createTestOrganizationWithPath(4, "Child1-2", &parentID1, "/1/4/", 1),
		*createTestOrganizationWithPath(5, "Child2-1", &parentID2, "/2/5/", 1),
	}

	mockRepo.On("FindAll", ctx).Return(allOrgs, nil)

	tree, err := uc.GetTree(ctx)

	require.NoError(t, err)
	assert.Len(t, tree, 2) // 2 roots

	// Check first root
	assert.Equal(t, "Root1", tree[0].Name)
	assert.Len(t, tree[0].Children, 2)

	// Check second root
	assert.Equal(t, "Root2", tree[1].Name)
	assert.Len(t, tree[1].Children, 1)

	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_GetTree_EmptyTree(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Organization{}, nil)

	tree, err := uc.GetTree(ctx)

	require.NoError(t, err)
	assert.Len(t, tree, 0)
	mockRepo.AssertExpectations(t)
}

// Create tests
func TestOrganizationUsecase_Create_RootOrganization(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("Create", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Name == "New Root" && org.ParentID == nil && org.Level == 0
	})).Run(func(args mock.Arguments) {
		org := args.Get(1).(*domain.Organization)
		org.ID = 1
	}).Return(nil)

	mockRepo.On("Update", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Name == "New Root" && org.Path == "/1/"
	})).Return(nil)

	org, err := uc.Create(ctx, "New Root", nil)

	require.NoError(t, err)
	assert.Equal(t, "New Root", org.Name)
	assert.Nil(t, org.ParentID)
	assert.Equal(t, 0, org.Level)
	assert.Equal(t, "/1/", org.Path)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Create_ChildOrganization(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	parentID := int64(1)
	parent := *createTestOrganizationWithPath(1, "Parent", nil, "/1/", 0)

	mockRepo.On("FindByID", ctx, int64(1)).Return(parent, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Name == "Child" && *org.ParentID == 1 && org.Level == 1
	})).Run(func(args mock.Arguments) {
		org := args.Get(1).(*domain.Organization)
		org.ID = 2
	}).Return(nil)

	mockRepo.On("Update", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Name == "Child" && org.Path == "/1/2/"
	})).Return(nil)

	org, err := uc.Create(ctx, "Child", &parentID)

	require.NoError(t, err)
	assert.Equal(t, "Child", org.Name)
	assert.Equal(t, int64(1), *org.ParentID)
	assert.Equal(t, 1, org.Level)
	assert.Equal(t, "/1/2/", org.Path)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Create_ParentNotFound(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	parentID := int64(999)
	mockRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	org, err := uc.Create(ctx, "Child", &parentID)

	assert.Nil(t, org)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parent organization not found")
	mockRepo.AssertExpectations(t)
}

// Update tests
func TestOrganizationUsecase_Update_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	existingOrg := *createTestOrganizationWithPath(1, "Old Name", nil, "/1/", 0)
	mockRepo.On("FindByID", ctx, int64(1)).Return(existingOrg, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Name == "New Name" && org.ID == 1
	})).Return(nil)

	org, err := uc.Update(ctx, 1, "New Name", nil)

	require.NoError(t, err)
	assert.Equal(t, "New Name", org.Name)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Update_ChangeParent(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	existingOrg := *createTestOrganizationWithPath(2, "Child", nil, "/2/", 0)
	newParent := *createTestOrganizationWithPath(1, "Parent", nil, "/1/", 0)
	newParentID := int64(1)

	mockRepo.On("FindByID", ctx, int64(2)).Return(existingOrg, nil)
	mockRepo.On("FindByID", ctx, int64(1)).Return(newParent, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.ID == 2 && org.Path == "/1/2/" && org.Level == 1
	})).Return(nil)

	org, err := uc.Update(ctx, 2, "Child", &newParentID)

	require.NoError(t, err)
	assert.Equal(t, int64(1), *org.ParentID)
	assert.Equal(t, "/1/2/", org.Path)
	assert.Equal(t, 1, org.Level)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Update_SelfAsParent(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	existingOrg := *createTestOrganizationWithPath(1, "Org", nil, "/1/", 0)
	selfID := int64(1)

	mockRepo.On("FindByID", ctx, int64(1)).Return(existingOrg, nil)

	org, err := uc.Update(ctx, 1, "Org", &selfID)

	assert.Nil(t, org)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot set self as parent")
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Update_CircularReference(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	// Org 1 (parent of Org 2) trying to become child of Org 2
	org1 := *createTestOrganizationWithPath(1, "Org1", nil, "/1/", 0)
	org2 := *createTestOrganizationWithPath(2, "Org2", func() *int64 { id := int64(1); return &id }(), "/1/2/", 1)
	org2ParentID := int64(2)

	mockRepo.On("FindByID", ctx, int64(1)).Return(org1, nil)
	mockRepo.On("FindByID", ctx, int64(2)).Return(org2, nil)

	org, err := uc.Update(ctx, 1, "Org1", &org2ParentID)

	assert.Nil(t, org)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular reference detected")
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Update_NotFound(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	org, err := uc.Update(ctx, 999, "New Name", nil)

	assert.Nil(t, org)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "organization not found")
	mockRepo.AssertExpectations(t)
}

// Delete tests
func TestOrganizationUsecase_Delete_Success(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("HasChildren", ctx, int64(1)).Return(false, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := uc.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Delete_HasChildren(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("HasChildren", ctx, int64(1)).Return(true, nil)

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete organization with children")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestOrganizationUsecase_Delete_RepositoryError(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	mockRepo.On("HasChildren", ctx, int64(1)).Return(false, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("database error"))

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// Helper function tests
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"exact match", "/1/", "/1/", true},
		{"beginning", "/1/2/3/", "/1/", true},
		{"middle", "/1/2/3/", "/2/", true},
		{"end", "/1/2/3/", "/3/", true},
		{"not found", "/1/2/3/", "/4/", false},
		{"empty substring", "/1/2/", "", true},
		{"empty string", "", "/1/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Edge case tests
func TestOrganizationUsecase_GetTree_DeepHierarchy(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	parent1 := int64(1)
	parent2 := int64(2)
	parent3 := int64(3)

	allOrgs := []domain.Organization{
		*createTestOrganizationWithPath(1, "Root", nil, "/1/", 0),
		*createTestOrganizationWithPath(2, "Level1", &parent1, "/1/2/", 1),
		*createTestOrganizationWithPath(3, "Level2", &parent2, "/1/2/3/", 2),
		*createTestOrganizationWithPath(4, "Level3", &parent3, "/1/2/3/4/", 3),
	}

	mockRepo.On("FindAll", ctx).Return(allOrgs, nil)

	tree, err := uc.GetTree(ctx)

	require.NoError(t, err)
	assert.Len(t, tree, 1) // 1 root
	assert.Equal(t, "Root", tree[0].Name)
	assert.Len(t, tree[0].Children, 1)
	assert.Equal(t, "Level1", tree[0].Children[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationUsecase_Create_DeepNesting(t *testing.T) {
	uc, mockRepo := setupOrganizationUsecase(t)
	ctx := context.Background()

	// Create at level 3
	parentID := int64(3)
	parent := *createTestOrganizationWithPath(3, "Level2", nil, "/1/2/3/", 2)

	mockRepo.On("FindByID", ctx, int64(3)).Return(parent, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Level == 3
	})).Run(func(args mock.Arguments) {
		org := args.Get(1).(*domain.Organization)
		org.ID = 4
	}).Return(nil)

	mockRepo.On("Update", ctx, mock.MatchedBy(func(org *domain.Organization) bool {
		return org.Path == "/1/2/3/4/"
	})).Return(nil)

	org, err := uc.Create(ctx, "Level3", &parentID)

	require.NoError(t, err)
	assert.Equal(t, 3, org.Level)
	assert.Equal(t, "/1/2/3/4/", org.Path)
	mockRepo.AssertExpectations(t)
}

// Additional methods for MockOrganizationRepositoryForUsecase
func (m *MockOrganizationRepositoryForUsecase) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrganizationRepositoryForUsecase) FindByPath(ctx context.Context, pathPrefix string) ([]domain.Organization, error) {
	args := m.Called(ctx, pathPrefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}
