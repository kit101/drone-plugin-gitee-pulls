package provider

import (
	"testing"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullRequestUpdateComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/demo1/pulls/comments/8006883").
		Reply(200).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	body := "<!-- drone-plugin:gitee-pull-request -->  [![Build Status](https://www.example.com/api/badges/kit101/demo1/status.svg?ref=refs/pulls/6/head)](https://www.example.com/build/link)"
	res, err := client.PullRequests.UpdateComment("kit101/demo1", 8006883, body)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 200); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
