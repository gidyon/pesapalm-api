package user

type User_ struct {
	ID            uint64         `json:"id,omitempty"`
	Phone         string         `json:"phone,omitempty"`
	Email         string         `json:"email,omitempty"`
	Names         string         `json:"names,omitempty"`
	BirthDate     string         `json:"birth_date,omitempty"`
	Gender        string         `json:"gender,omitempty"`
	ProfileURL    string         `json:"profile_url,omitempty"`
	Country       string         `json:"country,omitempty"`
	CountryCode   string         `json:"country_code,omitempty"`
	GroupId       int64          `json:"group_id,omitempty"`
	GeneralData   map[string]any `json:"general_data,omitempty"`
	PrimaryGroup  string         `json:"primary_group,omitempty"`
	Password      string         `json:"password,omitempty"`
	AccountStatus string         `json:"account_status,omitempty"`
	LastLogin     string         `json:"last_login,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RequestOtpRequest struct {
	Phone string `json:"phone,omitempty"`
}

type ValidateOtpRequest struct {
	Phone string `json:"phone,omitempty"`
	Otp   string `json:"otp,omitempty"`
}

type RefreshRequest struct {
	RefreshUuid string `json:"refresh_uuid,omitempty"`
}
