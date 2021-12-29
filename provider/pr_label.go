package provider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/drone/go-scm/scm"
)

type PullRequestLabelService interface {
	// List 获取某个 Pull Request 的所有标签
	List(repo string, number int, opts scm.ListOptions) (*scm.Response, []*Label, error)
	// Create 创建 Pull Request 标签
	Create(repo string, number int, names []string) (*scm.Response, []*Label, error)
	// Update 替换 Pull Request 所有标签
	Update(repo string, number int, names []string) (*scm.Response, []*Label, error)
	// Delete 删除 Pull Request 标签
	Delete(repo string, number int, names []string) (*scm.Response, []*Label, error)
}

type pullRequestLabelService struct {
	client *wrapper
}

func (p pullRequestLabelService) List(repo string, number int, opts scm.ListOptions) (*scm.Response, []*Label, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/labels?%s", repo, number, encodeListOptions(opts))
	var out []*Label
	res, err := p.client.do("GET", path, nil, &out)
	return res, out, err
}

func (p pullRequestLabelService) Create(repo string, number int, names []string) (*scm.Response, []*Label, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/labels", repo, number)
	var out []*Label
	res, err := p.client.do("POST", path, names, &out)
	return res, out, err
}

func (p pullRequestLabelService) Update(repo string, number int, names []string) (*scm.Response, []*Label, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/labels", repo, number)
	var out []*Label
	res, err := p.client.do("PUT", path, names, &out)
	return res, out, err
}

func (p pullRequestLabelService) Delete(repo string, number int, names []string) (*scm.Response, []*Label, error) {
	name := url.PathEscape(strings.Join(names, ","))
	path := fmt.Sprintf("repos/%s/pulls/%d/labels/%s", repo, number, name)
	var out []*Label
	res, err := p.client.do("DELETE", path, nil, &out)
	return res, out, err
}
