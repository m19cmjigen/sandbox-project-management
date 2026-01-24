package usecase

import (
	"context"
	"fmt"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// OrganizationUsecase は組織のユースケース
type OrganizationUsecase interface {
	GetAll(ctx context.Context) ([]domain.Organization, error)
	GetByID(ctx context.Context, id int64) (*domain.Organization, error)
	GetChildren(ctx context.Context, parentID int64) ([]domain.Organization, error)
	GetRoots(ctx context.Context) ([]domain.Organization, error)
	GetTree(ctx context.Context) ([]domain.OrganizationWithChildren, error)
	Create(ctx context.Context, name string, parentID *int64) (*domain.Organization, error)
	Update(ctx context.Context, id int64, name string, parentID *int64) (*domain.Organization, error)
	Delete(ctx context.Context, id int64) error
}

type organizationUsecase struct {
	orgRepo repository.OrganizationRepository
}

// NewOrganizationUsecase は新しいOrganizationUsecaseを作成
func NewOrganizationUsecase(orgRepo repository.OrganizationRepository) OrganizationUsecase {
	return &organizationUsecase{
		orgRepo: orgRepo,
	}
}

// GetAll は全ての組織を取得
func (u *organizationUsecase) GetAll(ctx context.Context) ([]domain.Organization, error) {
	return u.orgRepo.FindAll(ctx)
}

// GetByID はIDで組織を取得
func (u *organizationUsecase) GetByID(ctx context.Context, id int64) (*domain.Organization, error) {
	return u.orgRepo.FindByID(ctx, id)
}

// GetChildren は子組織を取得
func (u *organizationUsecase) GetChildren(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	return u.orgRepo.FindByParentID(ctx, parentID)
}

// GetRoots はルート組織を取得
func (u *organizationUsecase) GetRoots(ctx context.Context) ([]domain.Organization, error) {
	return u.orgRepo.FindRoots(ctx)
}

// GetTree は組織ツリーを取得（階層構造）
func (u *organizationUsecase) GetTree(ctx context.Context) ([]domain.OrganizationWithChildren, error) {
	// 全組織を取得
	allOrgs, err := u.orgRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 組織マップを作成
	orgMap := make(map[int64]*domain.OrganizationWithChildren)
	for _, org := range allOrgs {
		orgWithChildren := domain.OrganizationWithChildren{
			Organization: org,
			Children:     []domain.Organization{},
		}
		orgMap[org.ID] = &orgWithChildren
	}

	// ツリー構造を構築
	var roots []domain.OrganizationWithChildren
	for _, org := range allOrgs {
		if org.ParentID == nil {
			// ルート組織
			if orgWithChildren, ok := orgMap[org.ID]; ok {
				roots = append(roots, *orgWithChildren)
			}
		} else {
			// 子組織
			if parent, ok := orgMap[*org.ParentID]; ok {
				parent.Children = append(parent.Children, org)
			}
		}
	}

	return roots, nil
}

// Create は新しい組織を作成
func (u *organizationUsecase) Create(ctx context.Context, name string, parentID *int64) (*domain.Organization, error) {
	// 親組織の存在確認
	var path string
	var level int

	if parentID != nil {
		parent, err := u.orgRepo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("parent organization not found: id=%d", *parentID)
		}
		path = parent.Path
		level = parent.Level + 1
	} else {
		path = "/"
		level = 0
	}

	// 組織を作成
	org := &domain.Organization{
		Name:     name,
		ParentID: parentID,
		Path:     path, // 仮のパス、作成後に更新
		Level:    level,
	}

	err := u.orgRepo.Create(ctx, org)
	if err != nil {
		return nil, err
	}

	// パスを更新（自身のIDを含める）
	org.Path = fmt.Sprintf("%s%d/", path, org.ID)
	err = u.orgRepo.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Update は組織を更新
func (u *organizationUsecase) Update(ctx context.Context, id int64, name string, parentID *int64) (*domain.Organization, error) {
	// 既存組織の取得
	org, err := u.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found: id=%d", id)
	}

	// 親組織の変更チェック
	if parentID != nil {
		// 自分自身を親にできない
		if *parentID == id {
			return nil, fmt.Errorf("cannot set self as parent")
		}

		// 親組織の存在確認
		parent, err := u.orgRepo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("parent organization not found: id=%d", *parentID)
		}

		// 循環参照チェック（親のパスに自分のIDが含まれていないか）
		if contains(parent.Path, fmt.Sprintf("/%d/", id)) {
			return nil, fmt.Errorf("circular reference detected")
		}

		org.Path = fmt.Sprintf("%s%d/", parent.Path, org.ID)
		org.Level = parent.Level + 1
	} else {
		org.Path = fmt.Sprintf("/%d/", org.ID)
		org.Level = 0
	}

	// 更新
	org.Name = name
	org.ParentID = parentID

	err = u.orgRepo.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Delete は組織を削除
func (u *organizationUsecase) Delete(ctx context.Context, id int64) error {
	// 子組織の存在チェック
	hasChildren, err := u.orgRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return fmt.Errorf("cannot delete organization with children")
	}

	return u.orgRepo.Delete(ctx, id)
}

// contains は文字列に部分文字列が含まれるかチェック
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
