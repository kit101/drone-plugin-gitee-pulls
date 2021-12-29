package provider

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullRequestLabelsList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/demo1/pulls/3/labels").
		MatchParam("page", "1").
		MatchParam("per_page", "100").
		Reply(200).
		Type("application/json").
		File("testdata/pr_labels.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	opts := scm.ListOptions{Page: 1, Size: 100}
	_, got, err := client.PullRequestLabels.List("kit101/demo1", 3, opts)

	if err != nil {
		t.Error(err)
		return
	}

	var want []*Label
	raw, _ := ioutil.ReadFile("testdata/pr_labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestPullRequestLabelsCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/pulls/3/labels").
		Reply(201).
		Type("application/json").
		File("testdata/pr_labels.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	in := []string{"drone-build/running", "drone-build/pass"}
	_, got, err := client.PullRequestLabels.Create("kit101/demo1", 3, in)

	if err != nil {
		t.Error(err)
		return
	}

	var want []*Label
	raw, _ := ioutil.ReadFile("testdata/pr_labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestLabelsUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Put("/repos/kit101/demo1/pulls/3/labels").
		Reply(200).
		Type("application/json").
		File("testdata/pr_labels.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	in := []string{"drone-build/running", "drone-build/pass"}
	_, got, err := client.PullRequestLabels.Update("kit101/demo1", 3, in)

	if err != nil {
		t.Error(err)
		return
	}

	var want []*Label
	raw, _ := ioutil.ReadFile("testdata/pr_labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestPullRequestLabelsDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/demo1/pulls/3/labels/hello/1,hello/2").
		Reply(204).
		Type("application/json").
		File("testdata/pr_labels.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	names := []string{"hello/1", "hello/2"}
	res, got, err := client.PullRequestLabels.Delete("kit101/demo1", 3, names)

	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	var want []*Label
	raw, _ := ioutil.ReadFile("testdata/pr_labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
