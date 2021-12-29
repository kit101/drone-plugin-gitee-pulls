package main

import (
	"errors"
	"os"
	"strings"

	"github.com/kit101/drone-plugin-gitee-pulls/config"
	"github.com/kit101/drone-plugin-gitee-pulls/plugins"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitee plugin"
	app.Usage = "gitee plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug",
			EnvVar: "PLUGIN_DEBUG",
		},

		// gitee
		cli.StringFlag{
			Name:   "api-server",
			Usage:  "gitee api server url",
			EnvVar: "PLUGIN_API_SERVER",
			Value:  "https://gitee.com/api/v5",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "gitee api access-token to access gitee api",
			EnvVar: "PLUGIN_ACCESS_TOKEN,PLUGIN_GITEE_ACCESS_TOKEN,GITEE_ACCESS_TOKEN",
		},

		// drone
		cli.StringFlag{
			Name:   "host",
			Usage:  "drone server address",
			EnvVar: "DRONE_SYSTEM_HOST",
		},
		cli.StringFlag{
			Name:   "proto",
			Usage:  "drone server access proto: http/https",
			EnvVar: "DRONE_SYSTEM_PROTO",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "DRONE_REPO",
			EnvVar: "DRONE_REPO",
		},
		cli.IntFlag{
			Name:   "pull-request",
			Usage:  "drone pull request number",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "build-link",
			Usage:  "drone build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build-status",
			Usage:  "drone stage status",
			EnvVar: "DRONE_STAGE_STATUS",
		},
		cli.StringFlag{
			Name:   "ref",
			Usage:  "drone commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},

		// is running
		cli.BoolFlag{
			Name:   "is-running",
			Usage:  "build is running",
			EnvVar: "PLUGIN_IS_RUNNING",
		},

		// comment
		cli.BoolFlag{
			Name:   "comment",
			Usage:  "comment disabled",
			EnvVar: "PLUGIN_COMMENT_DISABLED",
		},

		// label
		cli.BoolFlag{
			Name:   "label",
			Usage:  "label disabled",
			EnvVar: "PLUGIN_LABEL_DISABLED",
		},
		cli.StringFlag{
			Name:   "label-running",
			Usage:  "running label",
			EnvVar: "PLUGIN_RUNNING_LABEL,PLUGIN_GITEE_RUNNING_LABEL,PLUGIN_LABEL_RUNNING",
			Value:  "drone-build/running,E6A23C",
		},
		cli.StringFlag{
			Name:   "label-success",
			Usage:  "success label",
			EnvVar: "PLUGIN_SUCCESS_LABEL,PLUGIN_GITEE_SUCCESS_LABEL,PLUGIN_LABEL_SUCCESS",
			Value:  "drone-build/success,67C23A",
		},
		cli.StringFlag{
			Name:   "label-failure",
			Usage:  "failure label",
			EnvVar: "PLUGIN_FAILURE_LABEL,PLUGIN_GITEE_FAILURE_LABEL,PLUGIN_LABEL_FAILURE",
			Value:  "drone-build/failure,DB2828",
		},

		// test
		cli.BoolFlag{
			Name:   "test",
			Usage:  "test disabled",
			EnvVar: "PLUGIN_TEST_DISABLED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	isDebug := c.Bool("debug")
	server := c.String("api-server")
	token := c.String("token")
	host := c.String("host")
	proto := c.String("proto")
	repo := c.String("repo")
	prNumber := c.Int("pull-request")
	buildLink := c.String("build-link")
	droneBuildStatus := c.String("build-status")
	ref := c.String("ref")

	isRunning := c.Bool("is-running")

	commentDisabled := c.Bool("comment")

	labelDisabled := c.Bool("label")
	runningLabel := c.String("label-running")
	successLabel := c.String("label-success")
	failureLabel := c.String("label-failure")

	testDisabled := c.Bool("test")

	buildStatus := getBuildStatus(isRunning, droneBuildStatus)

	conf := config.Config{
		GiteeApiUrl:       server,
		AccessToken:       token,
		DroneHost:         host,
		DroneProto:        proto,
		Repo:              repo,
		PullRequestNumber: prNumber,
		BuildLink:         buildLink,
		BuildStatus:       buildStatus,
		CommitRef:         ref,
		IsRunning:         isRunning,
		PluginComment: config.Comment{
			Disabled: commentDisabled,
		},
		PluginLabel: config.Label{
			Disabled: labelDisabled,
			Running:  runningLabel,
			Success:  successLabel,
			Failure:  failureLabel,
		},
		PluginTest: config.Test{
			Disabled: testDisabled,
		},
	}

	initLogger(isDebug)
	logrus.WithFields(logrus.Fields{
		"api_server":          conf.GiteeApiUrl,
		"drone_proto":         conf.DroneProto,
		"drone_host":          conf.DroneHost,
		"repo":                conf.Repo,
		"pull_request_number": conf.PullRequestNumber,
		"build_link":          conf.BuildLink,
		"build_status":        conf.BuildStatus.String(),
		"commit_ref":          conf.CommitRef,
		"is_running":          conf.IsRunning,
		"comment_disabled":    conf.PluginComment.Disabled,
		"label_disabled":      conf.PluginLabel.Disabled,
		"label_running":       conf.PluginLabel.Running,
		"label_success":       conf.PluginLabel.Success,
		"label_failure":       conf.PluginLabel.Failure,
		"test_disabled":       conf.PluginTest.Disabled,
	}).Debug("args")

	if strings.TrimSpace(token) == "" {
		return errors.New("token must not be blank")
	}
	plugin := plugins.NewPlugin(conf)
	return plugin.Exec()
}

func initLogger(debug bool) {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	//logrus.SetReportCaller(true)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		for _, env := range os.Environ() {
			logrus.Debug(env)
		}
	}
}

func getBuildStatus(isRunning bool, droneBuildStatus string) config.BuildStatus {
	if isRunning {
		return config.BuildStatusRunning
	} else {
		return config.BuildStatusOfValue(droneBuildStatus)
	}
}
