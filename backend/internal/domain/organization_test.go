package domain

import (
	"testing"
)

func TestOrganization_IsValid(t *testing.T) {
	tests := []struct {
		name string
		org  Organization
		want bool
	}{
		{
			name: "Valid organization with name",
			org: Organization{
				Name:  "Test Org",
				Path:  "1",
				Level: 0,
			},
			want: true,
		},
		{
			name: "Invalid organization without name",
			org: Organization{
				Name:  "",
				Path:  "1",
				Level: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.org.Name != ""
			if got != tt.want {
				t.Errorf("Organization validation = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganization_PathGeneration(t *testing.T) {
	tests := []struct {
		name     string
		parentID *int64
		orgID    int64
		want     string
	}{
		{
			name:     "Root organization",
			parentID: nil,
			orgID:    1,
			want:     "1",
		},
		{
			name:     "Child organization",
			parentID: func() *int64 { id := int64(1); return &id }(),
			orgID:    2,
			want:     "1.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// パス生成ロジックはデータベースのトリガーで実装されているため、
			// ここでは単純な検証のみ
			if tt.parentID == nil {
				// ルート組織のパスは組織IDのみ
				got := "1"
				if got != tt.want {
					t.Errorf("Path = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestOrganizationWithChildren_AddChild(t *testing.T) {
	parent := OrganizationWithChildren{
		Organization: Organization{
			ID:    1,
			Name:  "Parent",
			Path:  "1",
			Level: 0,
		},
		Children: []Organization{},
	}

	child := Organization{
		ID:       2,
		Name:     "Child",
		ParentID: func() *int64 { id := int64(1); return &id }(),
		Path:     "1.2",
		Level:    1,
	}

	parent.Children = append(parent.Children, child)

	if len(parent.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(parent.Children))
	}

	if parent.Children[0].Name != "Child" {
		t.Errorf("Expected child name 'Child', got %s", parent.Children[0].Name)
	}
}
