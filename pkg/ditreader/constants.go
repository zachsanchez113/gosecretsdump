package ditreader

const (
	nuSNCreated              = "ATTq131091"
	nuSNChanged              = "ATTq131192"
	nname                    = "ATTm3"
	nobjectGUID              = "ATTk589826"
	nobjectSid               = "ATTr589970"
	nuserAccountControl      = "ATTj589832"
	nprimaryGroupID          = "ATTj589922"
	naccountExpires          = "ATTq589983"
	nlogonCount              = "ATTj589993"
	nsAMAccountName          = "ATTm590045"
	nsAMAccountType          = "ATTj590126"
	nlastLogonTimestamp      = "ATTq589876"
	nuserPrincipalName       = "ATTm590480"
	nunicodePwd              = "ATTk589914"
	ndBCSPwd                 = "ATTk589879"
	nntPwdHistory            = "ATTk589918"
	nlmPwdHistory            = "ATTk589984"
	npekList                 = "ATTk590689"
	nsupplementalCredentials = "ATTk589949"
	npwdLastSet              = "ATTq589920"
)

var nnToInternal = map[string]string{
	"uSNCreated":              "ATTq131091",
	"uSNChanged":              "ATTq131192",
	"name":                    "ATTm3",
	"objectGUID":              "ATTk589826",
	"objectSid":               "ATTr589970",
	"userAccountControl":      "ATTj589832",
	"primaryGroupID":          "ATTj589922",
	"accountExpires":          "ATTq589983",
	"logonCount":              "ATTj589993",
	"sAMAccountName":          "ATTm590045",
	"sAMAccountType":          "ATTj590126",
	"lastLogonTimestamp":      "ATTq589876",
	"userPrincipalName":       "ATTm590480",
	"unicodePwd":              "ATTk589914",
	"dBCSPwd":                 "ATTk589879",
	"ntPwdHistory":            "ATTk589918",
	"lmPwdHistory":            "ATTk589984",
	"pekList":                 "ATTk590689",
	"supplementalCredentials": "ATTk589949",
	"pwdLastSet":              "ATTq589920",
}

var accTypes = map[int32]string{
	0x30000000: "SAM_NORMAL_USER_ACCOUNT",
	0x30000001: "SAM_MACHINE_ACCOUNT",
	0x30000002: "SAM_TRUST_ACCOUNT",
}