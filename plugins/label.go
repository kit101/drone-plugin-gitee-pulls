package plugins

import (
	"strings"

	"github.com/kit101/drone-plugin-gitee-pulls/config"
	"github.com/kit101/drone-plugin-gitee-pulls/provider"

	"github.com/sirupsen/logrus"
)

type LabelPlugin interface {
	// UpdatePullRequestStatusLabel 更新pull request的构建状态标签
	UpdatePullRequestStatusLabel(name, color string) error
}

type labelPlugin struct {
	client *provider.Client

	repo              string
	pullRequestNumber int
	labels            []string
}

func NewLabelPlugin(config config.Config, client *provider.Client, labels []string) LabelPlugin {
	return &labelPlugin{
		client: client,

		repo:              config.Repo,
		pullRequestNumber: config.PullRequestNumber,
		labels:            labels,
	}
}

func (l *labelPlugin) UpdatePullRequestStatusLabel(name, color string) error {
	logrus.Info("label:UpdatePullRequestStatusLabel")
	var names []string
	for _, labelName := range l.labels {
		if labelName != name {
			names = append(names, labelName)
		}
	}
	logrus.WithField("delete-labels", names).Debug("delete labels")
	_, currentLabels, err := l.client.PullRequestLabels.Delete(l.repo, l.pullRequestNumber, names)
	if err != nil {
		logrus.WithField("labels", currentLabels).Debug("current labels")
	}

	return l.create(name, color)
}

func (l *labelPlugin) create(name, color string) error {
	res, _, err := l.client.Labels.Get(l.repo, name)
	if err != nil && res != nil && res.Status == 404 {
		_, _, err = l.client.Labels.Create(l.repo, provider.Label{Name: name, Color: color})
		if err != nil {
			logrus.WithField("name", name).Warn("Create label failure", err)
		}
	}
	_, _, err = l.client.PullRequestLabels.Create(l.repo, l.pullRequestNumber, []string{name})
	return err
}

func GetLabelName(labelAndColor string) string {
	return strings.Split(labelAndColor, ",")[0]
}

func GetLabelColor(labelAndColor string) string {
	arrs := strings.Split(labelAndColor, ",")
	if len(arrs) > 1 {
		return arrs[1]
	}
	return ""
}
