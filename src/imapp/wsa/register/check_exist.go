package register

type CheckExist struct {
}

// CheckExistModel 检查手机号是否已有注册账号
type CheckExistModel struct {
	Cc       string `json:"cc"`
	In       string `json:"in"`
	Rc       string `json:"rc"`
	Lg       string `json:"lg"`
	Lc       string `json:"lc"`
	AuthKey  string `json:"authkey"`
	Eregid   string `json:"e_regid"`
	Ekeytype string `json:"e_keytype"`
	Eident   string `json:"e_ident"`
	EskeyId  string `json:"e_skey_id"`
	EskeyVal string `json:"e_skey_val"`
	EskeySig string `json:"e_skey_sig"`
	Fdid     string `json:"fdid"`
	Expid    string `json:"expid"`
	//TosVersion  string `json:"tos_version"`
	OfflineAb   string `json:"offline_ab"`
	Id          string `json:"id"`
	BackupToken string `json:"backup_token"`
}

type metricsModel struct {
	ExpidCd int32 `json:"expid_cd"`
	ExpidMd int32 `json:"expid_md"`
	RcC     bool  `json:"rc_c"`
}

type OfflineAbModel struct {
	Exposure []string      `json:"exposure"`
	Metrics  *metricsModel `json:"metrics"`
}

func MakeCheckExistUrl() string {
	return ""
}
