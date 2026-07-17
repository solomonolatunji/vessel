package utils

import "os"

func IsDryRun() bool {
	return os.Getenv("DEPLOY_DRY_RUN") == "true"
}
