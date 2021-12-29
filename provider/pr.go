package provider

import (
	"context"
	"fmt"
	"github.com/drone/go-scm/scm"
)

type PullRequestService interface {
	scm.PullRequestService
	// UpdateComment 更新pull request 评论
	UpdateComment(repo string, id int, body string) (*scm.Response, error)
}

type pullRequestService struct {
	client  *wrapper
	wrapper scm.PullRequestService
}

// UpdateComment 更新pull request 评论
func (p *pullRequestService) UpdateComment(repo string, id int, body string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/comments/%d", repo, id)
	in := map[string]string{"body": body}
	res, err := p.client.do("PATCH", path, in, nil)
	return res, err
}

func (p *pullRequestService) Find(ctx context.Context, s string, i int) (*scm.PullRequest, *scm.Response, error) {
	return p.wrapper.Find(ctx, s, i)
}

func (p *pullRequestService) FindComment(ctx context.Context, s string, i int, i2 int) (*scm.Comment, *scm.Response, error) {
	return p.wrapper.FindComment(ctx, s, i, i2)
}

func (p *pullRequestService) List(ctx context.Context, s string, options scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	return p.wrapper.List(ctx, s, options)
}

func (p *pullRequestService) ListChanges(ctx context.Context, s string, i int, options scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return p.wrapper.ListChanges(ctx, s, i, options)
}

func (p *pullRequestService) ListComments(ctx context.Context, s string, i int, options scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return p.wrapper.ListComments(ctx, s, i, options)
}

func (p *pullRequestService) ListCommits(ctx context.Context, s string, i int, options scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	return p.wrapper.ListCommits(ctx, s, i, options)
}

func (p *pullRequestService) Merge(ctx context.Context, s string, i int) (*scm.Response, error) {
	return p.wrapper.Merge(ctx, s, i)
}

func (p *pullRequestService) Close(ctx context.Context, s string, i int) (*scm.Response, error) {
	return p.wrapper.Close(ctx, s, i)
}

func (p *pullRequestService) Create(ctx context.Context, s string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	return p.wrapper.Create(ctx, s, input)
}

func (p *pullRequestService) CreateComment(ctx context.Context, s string, i int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return p.wrapper.CreateComment(ctx, s, i, input)
}

func (p *pullRequestService) DeleteComment(ctx context.Context, s string, i int, i2 int) (*scm.Response, error) {
	return p.wrapper.DeleteComment(ctx, s, i, i2)
}
