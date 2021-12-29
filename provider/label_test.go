package provider

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestLabelsList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/demo1/labels").
		Reply(200).
		Type("application/json").
		File("testdata/labels.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	_, got, err := client.Labels.List("kit101/demo1")
	if err != nil {
		t.Error(err)
		return
	}

	var want []*Label
	raw, _ := ioutil.ReadFile("testdata/labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestLabelsGet(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/demo1/labels/drone-build/running").
		Reply(200).
		Type("application/json").
		File("testdata/label.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	_, got, err := client.Labels.Get("kit101/demo1", "drone-build/running")
	if err != nil {
		t.Error(err)
		return
	}

	var want *Label
	raw, _ := ioutil.ReadFile("testdata/label.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestLabelsCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/kit101/demo1/labels").
		Reply(201).
		Type("application/json").
		File("testdata/label.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	in := Label{
		Name:  "drone-build/running",
		Color: "E6A23C",
	}

	_, got, err := client.Labels.Create("kit101/demo1", in)
	if err != nil {
		t.Error(err)
		return
	}

	var want *Label
	raw, _ := ioutil.ReadFile("testdata/label.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestLabelsDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Delete("/repos/kit101/demo1/labels/drone-build/running").
		Reply(204).
		Type("application/json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	res, err := client.Labels.Delete("kit101/demo1", "drone-build/running")
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(res.Status, 204); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestLabelsUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Put("/repos/kit101/demo1/labels/drone-build/running").
		Reply(200).
		Type("application/json").
		File("testdata/label.json")

	cfg := config.NewConfig()
	client, _ := NewClient(cfg)

	in := Label{
		Name:  "drone-build/running",
		Color: "E6A23C",
	}
	_, got, err := client.Labels.Update("kit101/demo1", "drone-build/running", in)
	if err != nil {
		t.Error(err)
		return
	}

	var want *Label
	raw, _ := ioutil.ReadFile("testdata/label.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
