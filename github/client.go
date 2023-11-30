package github

import (
	"github.com/google/go-github/v56/github"

	"github.com/starudream/go-lib/resty/v2"
)

func C() *github.Client {
	return github.NewClient(resty.C().GetClient()).WithAuthToken("")
}
