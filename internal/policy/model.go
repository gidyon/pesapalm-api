package policy

type PolicyRequest struct {
	Policy  string `json:"policy,omitempty"`
	Subject string `json:"subject,omitempty"`
	Object  string `json:"object,omitempty"`
	Action  string `json:"action,omitempty"`
}
