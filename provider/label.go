package provider

import (
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
)

type LabelService interface {
	// List 仓库任务标签列表
	List(repo string) (*scm.Response, []*Label, error)
	// Get 获取指定任务标签
	Get(repo string, name string) (*scm.Response, *Label, error)
	// Create 创建任务标签
	Create(repo string, in Label) (*scm.Response, *Label, error)
	// Delete 删除任务标签
	Delete(repo string, name string) (*scm.Response, error)
	// Update 更新任务标签
	Update(repo string, name string, in Label) (*scm.Response, *Label, error)
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type labelService struct {
	client *wrapper
}

func (l *labelService) List(repo string) (*scm.Response, []*Label, error) {
	path := fmt.Sprintf("repos/%s/labels", repo)
	var out []*Label
	res, err := l.client.do("GET", path, nil, &out)
	if err != nil {
		return res, nil, err
	}
	return res, out, nil
}

func (l *labelService) Get(repo string, name string) (*scm.Response, *Label, error) {
	path := fmt.Sprintf("repos/%s/labels/%s", repo, url.QueryEscape(name))
	var out *Label
	res, err := l.client.do("GET", path, nil, &out)
	if err != nil {
		return res, nil, err
	}
	return res, out, nil
}

func (l *labelService) Create(repo string, in Label) (*scm.Response, *Label, error) {
	path := fmt.Sprintf("repos/%s/labels", repo)
	var out *Label
	res, err := l.client.do("POST", path, in, &out)
	if err != nil {
		return res, nil, err
	}
	return res, out, nil
}

func (l *labelService) Delete(repo string, name string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/labels/%s", repo, name)
	res, err := l.client.do("DELETE", path, nil, nil)
	return res, err
}

func (l *labelService) Update(repo string, originalName string, in Label) (*scm.Response, *Label, error) {
	path := fmt.Sprintf("repos/%s/labels/%s", repo, originalName)
	var out *Label
	res, err := l.client.do("PUT", path, in, &out)
	if err != nil {
		return res, nil, err
	}
	return res, out, nil
}
