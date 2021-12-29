package plugins

import (
	"github.com/kit101/drone-plugin-gitee-pulls/config"
	"github.com/kit101/drone-plugin-gitee-pulls/provider"
)

type TestPlugin interface {
	// PushTestResult 推送测试结果
	PushTestResult() error
}

type testPlugin struct {
	client *provider.Client

	buildStatus       config.BuildStatus
	repo              string
	pullRequestNumber int
}

func NewTestPlugin(config config.Config, client *provider.Client) TestPlugin {
	return &testPlugin{
		client: client,

		buildStatus:       config.BuildStatus,
		repo:              config.Repo,
		pullRequestNumber: config.PullRequestNumber,
	}
}

func (t *testPlugin) PushTestResult() error {
	if t.buildStatus == config.BuildStatusSuccess {
		return t.passTest()
	} else {
		return t.resetTest()
	}
}

func (t *testPlugin) passTest() error {
	_, err := t.client.PullRequestTests.Pass(t.repo, t.pullRequestNumber, false)
	return err
}

func (t *testPlugin) resetTest() error {
	_, err := t.client.PullRequestTests.Reset(t.repo, t.pullRequestNumber, false)
	return err
}
