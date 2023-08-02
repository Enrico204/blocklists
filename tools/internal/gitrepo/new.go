package gitrepo

import (
	"go.uber.org/zap"
)

func New(logger *zap.SugaredLogger, repoURL string, tmpdir string) *GitRepo {
	var repo GitRepo
	repo.url = repoURL
	repo.dir = tmpdir
	repo.logger = logger
	return &repo
}
