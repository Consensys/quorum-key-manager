package utils

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/utils"
)

type Utilities struct {
	logger log.Logger
}

var _ utils.Utilities = &Utilities{}

func New(logger log.Logger) *Utilities {
	return &Utilities{
		logger: logger,
	}
}

// nolint
func NewUncoveredMethod() {
	for {
		fmt.Sprintf("no test")
		fmt.Sprintf("no test")
		fmt.Sprintf("no test")
		fmt.Sprintf("no test")
		fmt.Sprintf("no test")
		fmt.Sprintf("no test")
	}
}
