package fetcher

/*
 * API OBJECTS
 *
 * These contain the actual API things one cares about and are what are
 * returned by each of the fetch functions.
 *
 * NOTE: These currently implement what the Rollbar provider uses and should
 *       not be considered comprehensive as to what the API itself can provide.
 */
type AccessToken struct {
	AccessToken          string   `json:"access_token"`
	Name                 string   `json:"name"`
	ProjectID            int      `json:"project_id"`
	RateLimitWindowCount int      `json:"rate_limit_window_count",omitempty`
	RateLimitWindowSize  int      `json:"rate_limit_window_size",omitempty`
	Scopes               []string `json:"scopes"`
	Token                string   `json:"token"`
}

type Project struct {
	ID           int    `json:"id"`
	AccountID    int    `json:"account_id"`
	Name         string `json:name`
	AccessTokens []AccessToken
}
type Team struct {
	ID          int    `json:"id"`
	AccountID   int    `json:"account_id"`
	AccessLevel string `json:"access_level"`
	Name        string `json:"name"`
	Users       []int
	Projects    []int
}

type TeamProjects struct {
	TeamID    int `json:"team_id"`
	ProjectID int `json:"project_id"`
}

type TeamUsers struct {
	TeamID int `json:"team_id"`
	UserID int `json:"user_id"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Teams    []Team
}

/*
 * RESPONSES
 *
 * All API responses contain an error code value and a result object that may
 * or may not contain a nested JSON object.
 */
type accessTokenResponse struct {
	Err    int           `json:"err"`
	Result []AccessToken `json:"result"`
}

type projectResponse struct {
	Err    int       `json:err`
	Result []Project `json:"result"`
}

type projectTeamResponse struct {
	Err    int           `json:"err"`
	Result []interface{} `json:"result"`
}

type teamResponse struct {
	Err    int    `json:"err"`
	Result []Team `json:"result"`
}

type teamProjectsResponse struct {
	Err    int            `json:"err"`
	Result []TeamProjects `json:"result"`
}

type teamUsersResponse struct {
	Err    int         `json:"err"`
	Result []TeamUsers `json:"result"`
}

type userResponse struct {
	Err    int `json:"err"`
	Result struct {
		Users []User `json:"users"`
	} `json:"result"`
}

type userTeamResponse struct {
	Err    int `json:"err"`
	Result struct {
		Teams []Team `json:"teams"`
	} `json:"result"`
}
