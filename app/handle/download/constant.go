package download

import "gfs/app/common"

var (
	response = common.ResponseInstance()
)

const (
	resourceNotFount = "resource not found"
	resourceIsBreak  = "resource is breakdown"
	tokenError       = "token error"
	exError          = "token is expired"
)
