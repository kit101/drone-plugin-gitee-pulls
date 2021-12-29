package provider

import (
	"testing"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullRequestTestsPass(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/pulls/3/test").
		Reply(204).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	res, err := client.PullRequestTests.Pass("kit101/demo1", 3, false)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestTestsReset(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/demo1/pulls/3/test").
		Reply(200).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	res, err := client.PullRequestTests.Reset("kit101/demo1", 3, false)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 200); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestTestsAddTesters(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/pulls/3/testers").
		Reply(200).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	assignees := []string{"https://gitee.com/kit101"}
	res, err := client.PullRequestTests.AddTesters("kit101/demo1", 3, assignees)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 200); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestTestsDeleteTesters(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/demo1/pulls/3/testers").
		Reply(204).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	assignees := []string{"https://gitee.com/kit101"}
	res, err := client.PullRequestTests.DeleteTesters("kit101/demo1", 3, assignees)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
