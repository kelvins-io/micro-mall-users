package vars

type EmailConfigSettingS struct {
	Enable   bool   `json:"enable"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type JwtSettingS struct {
	Secret            string
	TokenExpireSecond int
}

type EmailNoticeSettingS struct {
	Receivers []string `json:"receivers"`
}
