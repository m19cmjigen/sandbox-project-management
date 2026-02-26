import type { Page, Route } from '@playwright/test'

// ---- Mock data ----

export const mockDashboardSummary = {
  total_projects: 10,
  red_projects: 3,
  yellow_projects: 2,
  green_projects: 5,
  total_issues: 80,
  red_issues: 15,
  yellow_issues: 10,
  green_issues: 55,
  organizations: [
    {
      id: 1,
      name: '開発本部',
      parent_id: null,
      level: 0,
      total_projects: 6,
      red_projects: 2,
      yellow_projects: 1,
      green_projects: 3,
      delay_status: 'RED',
      delay_rate: 0.33,
    },
    {
      id: 2,
      name: '第一開発部',
      parent_id: 1,
      level: 1,
      total_projects: 3,
      red_projects: 1,
      yellow_projects: 1,
      green_projects: 1,
      delay_status: 'RED',
      delay_rate: 0.33,
    },
    {
      id: 3,
      name: '営業本部',
      parent_id: null,
      level: 0,
      total_projects: 4,
      red_projects: 1,
      yellow_projects: 1,
      green_projects: 2,
      delay_status: 'RED',
      delay_rate: 0.25,
    },
  ],
}

export const mockProjectsResponse = {
  data: [
    {
      id: 1,
      jira_project_id: 'JP-001',
      key: 'PROJ-A',
      name: 'プロジェクトA',
      lead_account_id: null,
      lead_email: 'lead@example.com',
      organization_id: 1,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      red_count: 3,
      yellow_count: 2,
      green_count: 5,
      open_count: 10,
      total_count: 10,
      delay_status: 'RED',
    },
    {
      id: 2,
      jira_project_id: 'JP-002',
      key: 'PROJ-B',
      name: 'プロジェクトB',
      lead_account_id: null,
      lead_email: null,
      organization_id: 2,
      created_at: '2024-01-02T00:00:00Z',
      updated_at: '2024-01-02T00:00:00Z',
      red_count: 0,
      yellow_count: 1,
      green_count: 4,
      open_count: 5,
      total_count: 5,
      delay_status: 'YELLOW',
    },
  ],
  pagination: { page: 1, per_page: 25, total: 2, total_pages: 1 },
}

export const mockIssuesResponse = {
  data: [
    {
      id: 1,
      jira_issue_id: 'ISS-001',
      jira_issue_key: 'PROJ-A-1',
      project_id: 1,
      project_key: 'PROJ-A',
      project_name: 'プロジェクトA',
      summary: '遅延チケットのサンプル',
      status: 'In Progress',
      status_category: 'In Progress',
      due_date: '2024-01-01',
      assignee_name: '担当者A',
      assignee_account_id: 'acc-001',
      delay_status: 'RED',
      priority: 'High',
      issue_type: 'Task',
      last_updated_at: '2024-01-15T10:00:00Z',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-15T10:00:00Z',
    },
    {
      id: 2,
      jira_issue_id: 'ISS-002',
      jira_issue_key: 'PROJ-B-1',
      project_id: 2,
      project_key: 'PROJ-B',
      project_name: 'プロジェクトB',
      summary: '注意チケットのサンプル',
      status: 'To Do',
      status_category: 'To Do',
      due_date: '2024-02-01',
      assignee_name: null,
      assignee_account_id: null,
      delay_status: 'YELLOW',
      priority: 'Medium',
      issue_type: 'Story',
      last_updated_at: '2024-01-16T10:00:00Z',
      created_at: '2024-01-02T00:00:00Z',
      updated_at: '2024-01-16T10:00:00Z',
    },
  ],
  pagination: { page: 1, per_page: 25, total: 2, total_pages: 1 },
}

export const mockOrganizationsResponse = [
  {
    id: 1,
    name: '開発本部',
    parent_id: null,
    path: '/1/',
    level: 0,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    total_projects: 6,
    red_projects: 2,
    yellow_projects: 1,
    green_projects: 3,
    delay_status: 'RED',
  },
  {
    id: 2,
    name: '第一開発部',
    parent_id: 1,
    path: '/1/2/',
    level: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    total_projects: 3,
    red_projects: 1,
    yellow_projects: 1,
    green_projects: 1,
    delay_status: 'RED',
  },
]

// ---- Auth helpers ----

/**
 * Inject a mock admin user into localStorage so ProtectedRoute passes without
 * going through the actual login flow.
 * Must be called before page.goto() to take effect on the initial navigation.
 */
export async function setupAuth(page: Page): Promise<void> {
  // Zustand persist stores state as { state: {...}, version: 0 } under the key "auth-storage"
  await page.addInitScript(() => {
    const authState = {
      state: {
        token: 'mock-token-for-e2e',
        user: { id: 1, email: 'admin@example.com', role: 'admin' },
      },
      version: 0,
    }
    localStorage.setItem('auth-storage', JSON.stringify(authState))
  })
}

// ---- Route setup helpers ----

/** Intercept all /api/v1/* requests and return appropriate mock responses. */
export async function setupApiMocks(page: Page): Promise<void> {
  // Mock the /auth/me endpoint so the app doesn't attempt a real token validation
  await page.route('**/api/v1/auth/me', (route: Route) =>
    route.fulfill({ json: { id: 1, email: 'admin@example.com', role: 'admin' } }),
  )
  await page.route('**/api/v1/dashboard/summary', (route: Route) =>
    route.fulfill({ json: mockDashboardSummary }),
  )

  await page.route('**/api/v1/projects**', (route: Route) =>
    route.fulfill({ json: mockProjectsResponse }),
  )

  await page.route('**/api/v1/issues**', (route: Route) =>
    route.fulfill({ json: mockIssuesResponse }),
  )

  await page.route('**/api/v1/organizations**', (route: Route) =>
    route.fulfill({ json: mockOrganizationsResponse }),
  )
}
