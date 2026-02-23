package jiraclient

import "fmt"

const (
	projectSearchPath = "/rest/api/3/project/search"
	projectPageSize   = 50
)

// GetAllProjects fetches all Jira projects using paginated requests.
// It returns the complete list of projects across all pages.
func (c *Client) GetAllProjects() ([]Project, error) {
	var all []Project
	startAt := 0

	for {
		var resp ProjectSearchResponse
		path := fmt.Sprintf("%s?startAt=%d&maxResults=%d&expand=lead",
			projectSearchPath, startAt, projectPageSize)

		if err := c.get(path, &resp); err != nil {
			return nil, fmt.Errorf("get projects (startAt=%d): %w", startAt, err)
		}

		all = append(all, resp.Values...)

		if resp.IsLast || len(resp.Values) == 0 {
			break
		}
		startAt += len(resp.Values)
	}

	return all, nil
}
