package provider

import (
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
)

// PullRequestTestService Pull Request 测试
type PullRequestTestService interface {
	// Pass 处理 Pull Request 测试
	Pass(repo string, number int, force bool) (*scm.Response, error)
	// Reset 重置 Pull Request 测试 的状态
	Reset(repo string, number int, resetAll bool) (*scm.Response, error)
	// AddTesters 指派用户测试 Pull Request
	AddTesters(repo string, number int, testers []string) (*scm.Response, error)
	// DeleteTesters 取消用户测试 Pull Request
	DeleteTesters(repo string, number int, testers []string) (*scm.Response, error)
}

type pullRequestTestService struct {
	client *wrapper
}

func (p pullRequestTestService) Pass(repo string, number int, force bool) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/test", repo, number)
	params := map[string]bool{
		"name": force,
	}
	res, err := p.client.do("POST", path, params, nil)
	return res, err
}

func (p pullRequestTestService) Reset(repo string, number int, resetAll bool) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/testers", repo, number)
	params := map[string]bool{
		"reset_all": resetAll,
	}
	res, err := p.client.do("PATCH", path, params, nil)
	return res, err
}

func (p pullRequestTestService) AddTesters(repo string, number int, testers []string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/testers", repo, number)
	data := map[string]string{
		"assignees": strings.Join(testers, ","),
	}
	res, err := p.client.do("POST", path, data, nil)
	return res, err
}

func (p pullRequestTestService) DeleteTesters(repo string, number int, testers []string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/testers", repo, number)
	data := map[string]string{
		"assignees": strings.Join(testers, ","),
	}
	res, err := p.client.do("DELETE", path, data, nil)
	return res, err
}
