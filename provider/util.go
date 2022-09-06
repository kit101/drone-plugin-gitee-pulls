package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

// wrapper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	Client *scm.Client
	Ctx    context.Context
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}
	logRequest(method, path, in)
	// execute the http request
	res, err := c.Client.Do(c.Ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// parse the gitee request id.
	res.ID = res.Header.Get("X-Request-Id")

	// gitee pageValues
	populatePageValues(req, res)

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 300 {
		logResponse(res)
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		logResponse(res)
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	logResponse(res)
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a Gitee error.
type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// populatePageValues parses the HTTP Link response headers
// and populates the various pagination link values in the
// Response.
// response header: total_page, total_count
func populatePageValues(req *scm.Request, resp *scm.Response) {
	last, totalError := strconv.Atoi(resp.Header.Get("total_page"))
	reqURL, err := url.Parse(req.Path)
	if err != nil {
		return
	}
	current, currentError := strconv.Atoi(reqURL.Query().Get("page"))
	if totalError != nil && currentError != nil {
		return
	}
	resp.Page.First = 1
	if last != 0 {
		resp.Page.Last = last
	}
	if current != 0 {
		if current < resp.Page.Last {
			resp.Page.Next = current + 1
		} else {
			resp.Page.Next = resp.Page.Last
		}
		if current > resp.Page.First {
			resp.Page.Prev = current - 1
		} else {
			resp.Page.Prev = resp.Page.First
		}
	}
}

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func logRequest(method, path string, in interface{}) {
	logrus.Debug("request method: ", method, ", path: ", path, ", in: ", in)
}

func logResponse(resp *scm.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("read body error", err)
	}
	err = resp.Body.Close()
	if err != nil {
		logrus.Error("close body error", err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	logrus.Debug("response status: ", resp.Status, ", body: ", string(body), ", header: ", resp.Header)
}
