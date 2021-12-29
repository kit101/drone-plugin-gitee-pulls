package provider

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/kit101/drone-plugin-gitee-pulls/config"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/transport"
)

type Client struct {
	mu        sync.Mutex
	scmClient *scm.Client

	// HTTP client used to communicate with the API.
	Client *wrapper

	// Base URL for API requests.
	BaseURL *url.URL

	Repositories       scm.RepositoryService
	Linker             scm.Linker
	Labels             LabelService
	PullRequests       PullRequestService
	PullRequestLabels  PullRequestLabelService
	PullRequestReviews PullRequestReviewService
	PullRequestTests   PullRequestTestService
}

func NewClient(config config.Config) (*Client, error) {
	base, err := url.Parse(config.GiteeApiUrl)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	client := new(Client)
	client.scmClient, err = gitee.New(config.GiteeApiUrl)
	if err != nil {
		return nil, err
	}

	client.scmClient.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: config.AccessToken,
		},
	}
	client.Client = &wrapper{Client: client.scmClient, Ctx: context.Background()}

	// initialization services
	client.Repositories = client.scmClient.Repositories
	client.Linker = client.scmClient.Linker
	client.Labels = &labelService{client.Client}
	client.PullRequests = &pullRequestService{client: client.Client, wrapper: client.scmClient.PullRequests}
	client.PullRequestLabels = &pullRequestLabelService{client.Client}
	client.PullRequestReviews = &pullRequestReviewService{client.Client}
	client.PullRequestTests = &pullRequestTestService{client.Client}

	return client, nil
}
