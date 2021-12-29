package config

import "encoding/json"

type (
	BuildStatus int
	Config      struct {
		GiteeApiUrl string `json:"gitee_api_url" default:"https://gitee.com/api/v5"`
		AccessToken string `json:"access_token"`

		DroneProto        string      `json:"drone_proto"`
		DroneHost         string      `json:"drone_host"`
		Repo              string      `json:"repo"`
		PullRequestNumber int         `json:"pull_request_number"`
		BuildLink         string      `json:"build_link"`
		BuildStatus       BuildStatus `json:"build_status"`
		CommitRef         string      `json:"commit_ref"`

		IsRunning     bool    `json:"is_running" default:"false"`
		PluginComment Comment `json:"plugin_comment"`
		PluginLabel   Label   `json:"plugin_label"`
		PluginTest    Test    `json:"plugin_test"`
	}
	Comment struct {
		Disabled bool `json:"disabled" default:"false"`
	}
	Label struct {
		Disabled bool   `json:"disabled" default:"false"`
		Running  string `json:"running" default:"drone-build/running,E6A23C"`
		Success  string `json:"success" default:"drone-build/success,67C23A"`
		Failure  string `json:"failure" default:"drone-build/failure,DB2828"`
	}
	Test struct {
		Disabled bool `json:"disabled" default:"false"`
	}
)

func NewConfig() Config {
	return Config{
		GiteeApiUrl: "https://gitee.com/api/v5",
	}
}

func (c *Config) IsPullRequest() bool {
	return c.PullRequestNumber != 0
}

var statuses = []string{"unknown", "running", "success", "failure"}

const (
	BuildStatusUnknown BuildStatus = iota
	BuildStatusRunning
	BuildStatusSuccess
	BuildStatusFailure
)

func (b BuildStatus) String() (s string) {
	return statuses[b]
}

func BuildStatusOfValue(value string) BuildStatus {
	for i, status := range statuses {
		if status == value {
			return BuildStatus(i)
		}
	}
	return BuildStatusUnknown
}

// MarshalJSON returns the JSON-encoded BuildStatus.
func (b *BuildStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON unmarshal the JSON-encoded BuildStatus.
func (b *BuildStatus) UnmarshalJSON(data []byte) error {
	*b = BuildStatusOfValue(string(data))
	return nil
}
