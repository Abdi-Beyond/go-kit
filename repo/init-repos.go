package repo

// setup/repository.go
// -----------------------------------############ Note on Repository Initialization ############----------------------------------- //
// 🛠 Repository Initialization Guidelines
//
// 1 Initialize all repositories here.
// ⚠️ Use alias imports to avoid naming conflicts.

import (
	"github.com/Abdi-Beyond/go-kit/core/dependency"
	limitcheck "github.com/Abdi-Beyond/go-kit/modules/limitclient/repo"
)

type Repositories struct {
	Limitcheck *limitcheck.LimitCheckerRepo
}

func NewRepositories(deps *dependency.AppDependencies) *Repositories {
	return &Repositories{
		Limitcheck: limitcheck.NewLimitCheckerRepo(deps),
	}
}
