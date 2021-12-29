package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/kit101/drone-plugin-gitee-pulls/config"
	"github.com/kit101/drone-plugin-gitee-pulls/provider"

	"github.com/drone/go-scm/scm"
	"github.com/sirupsen/logrus"
)

var (
	commentHeader            = "<!-- drone-plugin:gitee-pull-request -->"
	commentBodyTemplateHttps = "|  drone-plugin/gitee-pulls ||\n|:---:|:---:|\n| **continuous-integration/drone/pr** — Build is **%s** [Details](%s) | [![Build status](%s)](%s) |"
	commentBodyTemplateHttp  = "|  drone-plugin/gitee-pulls |\n|:---:|\n| :%s: **continuous-integration/drone/pr** — Build is **%s** [Details](%s) |"
	runningIcon              = "clock1"
	successIcon              = "white_check_mark"
	failureIcon              = "x"
)

type CommentPlugin interface {
	// CommentBuildLink 评论drone build链接
	CommentBuildLink() error
}

type commentPlugin struct {
	ctx    context.Context
	client *provider.Client

	droneProto        string
	droneHost         string
	repo              string
	ref               string
	pullRequestNumber int
	buildStatus       config.BuildStatus
	buildLink         string
	isRunning         bool
}

func NewCommentPlugin(config config.Config, client *provider.Client) CommentPlugin {
	return &commentPlugin{
		client: client,
		ctx:    client.Client.Ctx,

		droneProto:        config.DroneProto,
		droneHost:         config.DroneHost,
		repo:              config.Repo,
		ref:               config.CommitRef,
		pullRequestNumber: config.PullRequestNumber,
		buildStatus:       config.BuildStatus,
		buildLink:         config.BuildLink,
		isRunning:         config.IsRunning,
	}
}

// CommentBuildLink 评论drone build链接
func (p *commentPlugin) CommentBuildLink() error {
	comments, _, _ := p.client.PullRequests.ListComments(p.ctx, p.repo, p.pullRequestNumber, scm.ListOptions{})
	var existing *scm.Comment
	for _, c := range comments {
		if strings.ContainsAny(c.Body, commentHeader) {
			existing = c
			break
		}
	}

	commentBody := p.buildCommentBody()
	commentInput := &scm.CommentInput{
		Body: commentBody,
	}

	var err error
	if existing != nil {
		_, err = p.client.PullRequests.UpdateComment(p.repo, existing.ID, commentInput.Body)
		if err != nil {
			logrus.Warn("update comment failure, will Create new comment, cause: ", err)
			_, _, err = p.client.PullRequests.CreateComment(p.ctx, p.repo, p.pullRequestNumber, commentInput)
		}
	} else {
		_, _, err = p.client.PullRequests.CreateComment(p.ctx, p.repo, p.pullRequestNumber, commentInput)
	}
	return err
}

func (p *commentPlugin) buildCommentBody() string {
	var body string
	if p.droneProto == "https" {
		badgeLink := p.buildBadgeLink()
		body = fmt.Sprintf(commentBodyTemplateHttps, p.buildStatus, p.buildLink, badgeLink, p.buildLink)
	} else {
		var icon string
		switch p.buildStatus {
		case config.BuildStatusRunning:
			icon = runningIcon
		case config.BuildStatusSuccess:
			icon = successIcon
		case config.BuildStatusFailure:
			icon = failureIcon
		}
		body = fmt.Sprintf(commentBodyTemplateHttp, icon, p.buildStatus, p.buildLink)
	}
	return commentHeader + "\n" + body
}

func (p *commentPlugin) buildBadgeLink() string {
	return fmt.Sprintf("%s://%s/api/badges/%s/status.svg?ref=%s", p.droneProto, p.droneHost, p.repo, p.ref)
}
