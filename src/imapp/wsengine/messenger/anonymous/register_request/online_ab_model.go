package registerRequest

type metricsModel struct {
	ExpidCd int32 `json:"expid_cd"`
	ExpidMd int32 `json:"expid_md"`
	RcC     bool  `json:"rc_c"`
}

type offlineAbModel struct {
	Exposure []string      `json:"exposure"`
	Metrics  *metricsModel `json:"metrics"`
}
