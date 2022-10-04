package harbor

import "time"

type ScanResultsJson struct {
	VulnerabilityReport struct {
		GeneratedAt time.Time `json:"generated_at,omitempty"`
		Scanner     struct {
			Name    string `json:"name,omitempty"`
			Vendor  string `json:"vendor,omitempty"`
			Version string `json:"version,omitempty"`
		} `json:"scanner,omitempty"`
		Severity        string `json:"severity,omitempty"`
		Vulnerabilities []struct {
			ID              string   `json:"id,omitempty"`
			Package         string   `json:"package,omitempty"`
			Version         string   `json:"version,omitempty"`
			FixVersion      string   `json:"fix_version,omitempty"`
			Severity        string   `json:"severity,omitempty"`
			Description     string   `json:"description,omitempty"`
			Links           []string `json:"links,omitempty"`
			ArtifactDigests []string `json:"artifact_digests,omitempty"`
			PreferredCvss   struct {
				ScoreV3  string      `json:"score_v3,omitempty"`
				ScoreV2  interface{} `json:"score_v2,omitempty"`
				VectorV3 string      `json:"vector_v3,omitempty"`
				VectorV2 string      `json:"vector_v2,omitempty"`
			} `json:"preferred_cvss,omitempty"`
			CweIds           []string `json:"cwe_ids,omitempty"`
			VendorAttributes struct {
				Cvss struct {
					Nvd struct {
						V3Score  float64 `json:"V3Score,omitempty"`
						V3Vector string  `json:"V3Vector,omitempty"`
					} `json:"nvd,omitempty"`
				} `json:"CVSS,omitempty"`
			} `json:"vendor_attributes,omitempty"`
		} `json:"vulnerabilities,omitempty"`
	} `json:"application/vnd.security.vulnerability.report; version=1.1,omitempty"`
}
