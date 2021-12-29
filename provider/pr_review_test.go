package provider

import (
	"testing"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullRequestReviewsPass(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/pulls/3/review").
		Reply(204).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	res, err := client.PullRequestReviews.Pass("kit101/demo1", 3, false)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestReviewsReset(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/kit101/demo1/pulls/3/assignees").
		Reply(200).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	res, err := client.PullRequestReviews.Reset("kit101/demo1", 3, false)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 200); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestReviewsAddAssignees(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/pulls/3/assignees").
		Reply(200).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	assignees := []string{"https://gitee.com/kit101"}
	res, err := client.PullRequestReviews.AddAssignees("kit101/demo1", 3, assignees)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 200); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestReviewsDeleteAssignees(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/demo1/pulls/3/assignees").
		Reply(204).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	assignees := []string{"https://gitee.com/kit101"}
	res, err := client.PullRequestReviews.DeleteAssignees("kit101/demo1", 3, assignees)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
