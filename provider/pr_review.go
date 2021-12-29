package provider

import (
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
)

// PullRequestReviewService Pull Request 审查
type PullRequestReviewService interface {
	// Pass 处理 Pull Request 审查
	Pass(repo string, number int, force bool) (*scm.Response, error)
	// Reset 重置 Pull Request 审查 的状态
	Reset(repo string, number int, resetAll bool) (*scm.Response, error)
	// AddAssignees 指派用户审查 Pull Request
	AddAssignees(repo string, number int, assignees []string) (*scm.Response, error)
	// DeleteAssignees 取消用户审查 Pull Request
	DeleteAssignees(repo string, number int, assignees []string) (*scm.Response, error)
}

type pullRequestReviewService struct {
	client *wrapper
}

func (p pullRequestReviewService) Pass(repo string, number int, force bool) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/review", repo, number)
	params := map[string]bool{
		"name": force,
	}
	res, err := p.client.do("POST", path, params, nil)
	return res, err
}

func (p pullRequestReviewService) Reset(repo string, number int, resetAll bool) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/assignees", repo, number)
	params := map[string]bool{
		"reset_all": resetAll,
	}
	res, err := p.client.do("PATCH", path, params, nil)
	return res, err
}

func (p pullRequestReviewService) AddAssignees(repo string, number int, assignees []string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/assignees", repo, number)
	data := map[string]string{
		"assignees": strings.Join(assignees, ","),
	}
	res, err := p.client.do("POST", path, data, nil)
	return res, err
}

func (p pullRequestReviewService) DeleteAssignees(repo string, number int, assignees []string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/assignees", repo, number)
	data := map[string]string{
		"assignees": strings.Join(assignees, ","),
	}
	res, err := p.client.do("DELETE", path, data, nil)
	return res, err
}
