package tax1099

// TinType defines the allowed values for the TinType field
type TinType string

const (
	TinTypeIndividual TinType = "Individual"
	TinTypeBusiness   TinType = "Business"
)

// ValidationError represents details about validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Source  string `json:"source"`
	Message string `json:"message"`
}
