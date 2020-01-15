package manager

import "gfs/app/common"

const (
	formParamError      = "Request from param is missing"
	paramError          = "Request param is missing"
	loginNameOrPsdError = "Login name or password error"
	IllegalRequest      = "Illegal request"
	tokenError          = "AuthToken is Expired"
	userInfoError       = "User info error"
	ormOptionsFailed    = "orm options failed"
	databaseError       = "database error"
)

const (
	initPassword = "123456"
)

var (
	response = common.ResponseInstance()
)
