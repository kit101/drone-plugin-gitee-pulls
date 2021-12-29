package plugins

import (
	"errors"
	"strings"

	"github.com/kit101/drone-plugin-gitee-pulls/config"
	"github.com/kit101/drone-plugin-gitee-pulls/provider"
)

type Plugin interface {
	Exec() error
}

type plugin struct {
	config config.Config
	client *provider.Client
}

func NewPlugin(config config.Config) Plugin {
	return &plugin{config: config}
}

func (p *plugin) Exec() error {
	var err error
	p.client, err = provider.NewClient(p.config)
	if err != nil {
		return err
	}

	if !p.config.PluginComment.Disabled {
		err = p.commentBuildLink()
		if err != nil {
			return err
		}
	}

	if !p.config.PluginLabel.Disabled {
		err = p.labelBuildStatus()
		if err != nil {
			return err
		}
	}

	if !p.config.PluginTest.Disabled {
		err = p.pushTestStatus()
		if err != nil {
			return err
		}
	}

	return err
}

func (p *plugin) commentBuildLink() error {
	comment := NewCommentPlugin(p.config, p.client)
	return comment.CommentBuildLink()
}

func (p *plugin) labelBuildStatus() error {
	labels := []string{
		GetLabelName(p.config.PluginLabel.Running),
		GetLabelName(p.config.PluginLabel.Success),
		GetLabelName(p.config.PluginLabel.Failure),
	}
	labelPlugin := NewLabelPlugin(p.config, p.client, labels)

	var labelValue string
	switch p.config.BuildStatus {
	case config.BuildStatusRunning:
		labelValue = p.config.PluginLabel.Running
	case config.BuildStatusSuccess:
		labelValue = p.config.PluginLabel.Success
	case config.BuildStatusFailure:
		labelValue = p.config.PluginLabel.Failure
	default:
		return errors.New("unknown build status")
	}

	labelValues := strings.Split(labelValue, ",")
	var name, color string
	if len(labelValues) > 0 {
		name = labelValues[0]
	}
	if len(labelValues) > 1 {
		color = labelValues[1]
	}

	return labelPlugin.UpdatePullRequestStatusLabel(name, color)
}

func (p *plugin) pushTestStatus() error {
	testPlugin := NewTestPlugin(p.config, p.client)
	return testPlugin.PushTestResult()
}
