package harbor

import "time"

type FindImageJson []struct {
	Accessories   interface{} `json:"accessories,omitempty"`
	AdditionLinks struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute,omitempty"`
			Href     string `json:"href,omitempty"`
		} `json:"build_history,omitempty"`
		Vulnerabilities struct {
			Absolute bool   `json:"absolute,omitempty"`
			Href     string `json:"href,omitempty"`
		} `json:"vulnerabilities,omitempty"`
	} `json:"addition_links,omitempty"`
	Digest     string `json:"digest,omitempty"`
	ExtraAttrs struct {
		Architecture string `json:"architecture,omitempty"`
		Author       string `json:"author,omitempty"`
		Config       struct {
			Entrypoint   []string `json:"Entrypoint,omitempty"`
			Env          []string `json:"Env,omitempty"`
			ExposedPorts struct {
				Eight080TCP struct {
				} `json:"8080/tcp,omitempty"`
			} `json:"ExposedPorts,omitempty"`
			Labels struct {
				Maintainer string `json:"maintainer,omitempty"`
			} `json:"Labels,omitempty"`
			User       string `json:"User,omitempty"`
			WorkingDir string `json:"WorkingDir,omitempty"`
		} `json:"config,omitempty"`
		Created time.Time `json:"created,omitempty"`
		Os      string    `json:"os,omitempty"`
	} `json:"extra_attrs,omitempty"`
	Icon              string      `json:"icon,omitempty"`
	ID                int         `json:"id,omitempty"`
	Labels            interface{} `json:"labels,omitempty"`
	ManifestMediaType string      `json:"manifest_media_type,omitempty"`
	MediaType         string      `json:"media_type,omitempty"`
	ProjectID         int         `json:"project_id,omitempty"`
	PullTime          time.Time   `json:"pull_time,omitempty"`
	PushTime          time.Time   `json:"push_time,omitempty"`
	References        interface{} `json:"references,omitempty"`
	RepositoryID      int         `json:"repository_id,omitempty"`
	ScanOverview      struct {
		VulnerabilityReport struct {
			CompletePercent int       `json:"complete_percent,omitempty"`
			Duration        int       `json:"duration,omitempty"`
			EndTime         time.Time `json:"end_time,omitempty"`
			ReportID        string    `json:"report_id,omitempty"`
			ScanStatus      string    `json:"scan_status,omitempty"`
			Scanner         struct {
				Name    string `json:"name,omitempty"`
				Vendor  string `json:"vendor,omitempty"`
				Version string `json:"version,omitempty"`
			} `json:"scanner,omitempty"`
			Severity  string    `json:"severity,omitempty"`
			StartTime time.Time `json:"start_time,omitempty"`
			Summary   struct {
				Fixable int `json:"fixable,omitempty"`
				Summary struct {
				} `json:"summary,omitempty"`
				Total int `json:"total,omitempty"`
			} `json:"summary,omitempty"`
		} `json:"application/vnd.security.vulnerability.report; version=1.1,omitempty"`
	} `json:"scan_overview,omitempty"`
	Size int `json:"size,omitempty"`
	Tags []struct {
		ArtifactID   int       `json:"artifact_id,omitempty"`
		ID           int       `json:"id,omitempty"`
		Immutable    bool      `json:"immutable,omitempty"`
		Name         string    `json:"name,omitempty"`
		PullTime     time.Time `json:"pull_time,omitempty"`
		PushTime     time.Time `json:"push_time,omitempty"`
		RepositoryID int       `json:"repository_id,omitempty"`
		Signed       bool      `json:"signed,omitempty"`
	} `json:"tags,omitempty"`
	Type string `json:"type,omitempty"`
}
