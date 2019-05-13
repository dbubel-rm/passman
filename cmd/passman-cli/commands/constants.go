package commands

const (
	REGISTER_ACCOUNT  = "init"
	NUKE_ACCOUNT      = "nuke"
	HELP              = "Help"
	VERSION           = "version"
	GEN_PASS          = "rand"
	INSERT_CREDENTIAL = "insert"
	LOGIN             = "login"
	PASSMAN_MASTER    = "PASSMAN_MASTER"
	GET_CREDENTIAL    = "get"
	RM_CREDENTIAL     = "rm"
	GET_SERVICES      = "list"
	UPDATE_CREDENTIAL = "update"
	UPDATE_MASTER     = "master"
	baseURL = "http://localhost:3000"
)

var PassmanHome = "/.passman/session.json"
