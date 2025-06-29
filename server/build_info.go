package server

import (
	"encoding/json"
	"os"
)

type BuildInfo struct {
	BuildTime      string `json:"buildTime"`
	GitCommit      string `json:"gitCommit"`
	GitShortCommit string `json:"gitShortCommit"`
	GitBranch      string `json:"gitBranch"`
	GitTag         string `json:"gitTag"`
	GithubUrl      string `json:"githubUrl"`
}

func ReadBuildInfo() (*BuildInfo, error) {
	data, err := os.ReadFile("/app/build-info.json")
	if err != nil {
		// Return local development build info when file doesn't exist
		if os.IsNotExist(err) {
			return &BuildInfo{
				BuildTime:      "localdev",
				GitCommit:      "localdev",
				GitShortCommit: "localdev",
				GitBranch:      "localdev",
				GitTag:         "localdev",
				GithubUrl:      "https://github.com/tokuhirom/blog4",
			}, nil
		}
		return nil, err
	}

	var buildInfo BuildInfo
	err = json.Unmarshal(data, &buildInfo)
	if err != nil {
		return nil, err
	}

	return &buildInfo, nil
}
