package github

import "time"

type IssueComments []struct {
	URL      string `json:"url,omitempty"`
	HTMLURL  string `json:"html_url,omitempty"`
	IssueURL string `json:"issue_url,omitempty"`
	ID       int    `json:"id,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	User     struct {
		Login             string `json:"login,omitempty"`
		ID                int    `json:"id,omitempty"`
		NodeID            string `json:"node_id,omitempty"`
		AvatarURL         string `json:"avatar_url,omitempty"`
		GravatarID        string `json:"gravatar_id,omitempty"`
		URL               string `json:"url,omitempty"`
		HTMLURL           string `json:"html_url,omitempty"`
		FollowersURL      string `json:"followers_url,omitempty"`
		FollowingURL      string `json:"following_url,omitempty"`
		GistsURL          string `json:"gists_url,omitempty"`
		StarredURL        string `json:"starred_url,omitempty"`
		SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
		OrganizationsURL  string `json:"organizations_url,omitempty"`
		ReposURL          string `json:"repos_url,omitempty"`
		EventsURL         string `json:"events_url,omitempty"`
		ReceivedEventsURL string `json:"received_events_url,omitempty"`
		Type              string `json:"type,omitempty"`
		SiteAdmin         bool   `json:"site_admin,omitempty"`
	} `json:"user,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
	AuthorAssociation string    `json:"author_association,omitempty"`
	Body              string    `json:"body,omitempty"`
	Reactions         struct {
		URL        string `json:"url,omitempty"`
		TotalCount int    `json:"total_count,omitempty"`
		Num1       int    `json:"+1,omitempty"`
		Num10      int    `json:"-1,omitempty"`
		Laugh      int    `json:"laugh,omitempty"`
		Hooray     int    `json:"hooray,omitempty"`
		Confused   int    `json:"confused,omitempty"`
		Heart      int    `json:"heart,omitempty"`
		Rocket     int    `json:"rocket,omitempty"`
		Eyes       int    `json:"eyes,omitempty"`
	} `json:"reactions,omitempty"`
	PerformedViaGithubApp interface{} `json:"performed_via_github_app,omitempty"`
}
