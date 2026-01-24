package domain

import "time"

// Organization は組織階層を表すエンティティ
type Organization struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ParentID  *int64    `db:"parent_id" json:"parent_id"`
	Path      string    `db:"path" json:"path"`
	Level     int       `db:"level" json:"level"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// OrganizationWithChildren は子組織を含む組織
type OrganizationWithChildren struct {
	Organization
	Children []Organization `json:"children,omitempty"`
}
