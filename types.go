package tax1099

// Environment defines the allowed values for the Environment field
type Environment string

var (
	EnvironmentStaging    Environment = "staging"
	EnvironmentProduction Environment = "production"
)

type UrlType string

var (
	UrlMain    UrlType = "main"
	UrlPayment UrlType = "payment"
	Url1098    UrlType = "1098"
)

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

// PayerInfo represents the details of a payer
type PayerInfo struct {
	ID                     int     `json:"payerId"`                 //PayerID is the unique identifier for the payer in Tax1099's system
	ClientID               string  `json:"clientPayerId,omitempty"` //ClientPayerID is the unique identifier for the payer in your system
	TinType                TinType `json:"tinType"`                 //TinType is the type of TIN for the payer, whether business or individual
	TaxIdentifer           string  `json:"payerTin"`                //TaxIdentifier is the payer's TIN, SSN, or EIN with no dashes
	FirstName              string  `json:"firstName,omitempty"`     //FirstName is the first name of the payer, if individual
	MiddleName             string  `json:"middleName,omitempty"`    //MiddleName is the middle name of the payer, if individual
	LastNameOrBusinessName string  `json:"lastNameOrBusinessName"`  //LastNameOrBusinessName is the last name of the payer, if individual, or the business name
	Suffix                 string  `json:"suffix,omitempty"`        //Suffix is the suffix of the payer, if individual
	Address                string  `json:"address"`                 //Address is the street address (line1) of the payer
	Address2               string  `json:"address2,omitempty"`      //Address2 is the street address (line2) of the payer
	City                   string  `json:"city"`                    //City is the city of the payer
	State                  string  `json:"state"`                   //State is the two-letter abbreviation for the state
	ZipCode                string  `json:"zipCode"`                 //ZipCode is the 5-digit zip code or 9-digit zip code with a hyphen
	Country                string  `json:"country"`                 //Country is the two-letter abbreviation for the country
	Email                  string  `json:"email,omitempty"`         //Email is the email address of the payer, optional
	PhoneNumber            string  `json:"phone"`                   //PhoneNumber is the phone number of the payer, is required
	LastFiling             bool    `json:"lastFiling"`              //LastFiling is used to indicate if this is the last filing for the payer with no more expected in the future
	DisregardedEntity      string  `json:"disregardedEntity"`       //DisregardedEntity is used to indicate if the payer is a disregarded entity where the payer is disregarded
	UnMaskRecipientTin     bool    `json:"unMaskRecipientTin"`      //UnMaskRecipientTin is used to indicate whether to mask or unmask the TIN
}

// RecipientInfo represents the details of a recipient
type RecipientInfo struct {
	PayerID                int    `json:"payerId"`                     //PayerID is the unique identifier for the payer in Tax1099's system
	RecipientID            int    `json:"recipientId"`                 //RecipientID is the unique identifier for the recipient in Tax1099's system
	ClientID               string `json:"clientRecipientId,omitempty"` //ClientRecipientID is the unique identifier for the recipient in your system
	TinType                string `json:"tinType"`                     //TinType is the type of TIN for the recipient, whether business or individual
	TaxIdentifer           string `json:"recipientTin"`                //TaxIdentifier is the recipient's TIN, SSN, or EIN with no dashes.
	FirstName              string `json:"firstName,omitempty"`         //FirstName is the first name of the recipient, if individual
	MiddleName             string `json:"middleName,omitempty"`        //MiddleName is the middle name of the recipient, if individual
	LastNameOrBusinessName string `json:"lastNameOrBusinessName"`      //LastNameOrBusinessName is the last name of the recipient, if individual, or the business name
	Suffix                 string `json:"suffix,omitempty"`            //Suffix is the suffix of the recipient, if individual
	Address                string `json:"address"`                     //Address is the street address (line1) of the recipient
	Address2               string `json:"address2,omitempty"`          //Address2 is the street address (line2) of the recipient
	City                   string `json:"city"`                        //City is the city of the recipient
	State                  string `json:"state"`                       //State is the two-letter abbreviation for the state
	ZipCode                string `json:"zipCode"`                     //ZipCode is the 5-digit zip code or 9-digit zip code with a hyphen
	Country                string `json:"country"`                     //Country is the two-letter abbreviation for the country
	Email                  string `json:"email"`                       //Email is the email address of the recipient, optional
	PhoneNumber            string `json:"phone"`                       //PhoneNumber is the phone number of the recipient, is required
	AttentionTo            string `json:"attentionTo,omitempty"`       //AttentionTo is the name of the person to whom the form should be addressed, if a business
	IsActive               bool   `json:"isActive"`                    //IsActive is used to indicate if the recipient is active or inactive
}
