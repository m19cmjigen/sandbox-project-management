package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrgDelayStatus_Red(t *testing.T) {
	org := &OrganizationRow{RedProjects: 3, YellowProjects: 1, GreenProjects: 0}
	assert.Equal(t, "RED", orgDelayStatus(org))
}

func TestOrgDelayStatus_Yellow(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 2, GreenProjects: 1}
	assert.Equal(t, "YELLOW", orgDelayStatus(org))
}

func TestOrgDelayStatus_Green(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 0, GreenProjects: 5}
	assert.Equal(t, "GREEN", orgDelayStatus(org))
}

func TestOrgDelayStatus_AllZero(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 0, GreenProjects: 0}
	assert.Equal(t, "GREEN", orgDelayStatus(org))
}
