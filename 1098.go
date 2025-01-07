package tax1099

// Submit1098Request represents the JSON structure for submitting 1098 forms
type Submit1098Request struct {
	TaxYear string `json:"taxYear"`
	Items   []Item `json:"items"`
}

// Item represents a single payer and their associated forms
type Item struct {
	PayerInfo PayerInfo `json:"payerInfo"` //Payer Info is the seller/lender of the property
	Forms     []Form    `json:"forms"`
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
	UnMaskRecipientTin     bool    `json:"unMaskRecipientTin"`
	CombinedFedStateFiling bool    `json:"combinedFedStateFiling"`
}

// Form represents the details of a single 1098 form
type Form struct {
	RecipientInfo       RecipientInfo `json:"recipientInfo"`      //Recipient Info is the buyer/borrower of the property
	TaxYear             string        `json:"taxYear"`            //Tax Year is the year for which the form is being filed
	AcctNo              string        `json:"acctNo"`             //Account Number is required if you have multiple accounts for a payer/borrower for whom you are filing more than one Form 1098
	MortgageInterest    float64       `json:"mortgageInterest"`   //Mortgage Interest is the amount of interest received from the borrower during the tax year
	PrincipalResidence  float64       `json:"principalResidence"` //Principal Residence is the points paid by the borrower for the residence
	OverpaidInterest    float64       `json:"overpaidInterest"`   //Overpaid Interest is the amount of interest received from the borrower that was refunded or credited during the tax year
	MortgagePremiums    float64       `json:"mortgagePremiums"`   //Mortgage Premiums is not used currently.
	MortgagePrincipal   float64       `json:"mortgagePrincipal"`  //Mortgage Principal is amount of principal as the start of the calendar year
	MortgageDate        string        `json:"mortgageDate"`       //Mortgage Date is the date the mortgage was originated
	IsAddressSame       bool          `json:"isAddressSame"`      //Is Address Same is used to indicate if the property address is the same as the recipient address
	PropertyAddress     string        `json:"propertyAddress"`    //Property Address is the address of the property for which the form is being filed
	PropertyDescription string        `json:"propertyDescription"`
	USPSMail            bool          `json:"uspsMail"`        //USPS Mail is used to indicate if the form should be mailed to the payer
	TINCheck            bool          `json:"tinCheck"`        //TIN Check is used to indicate if the TIN should be checked
	EDelivery           bool          `json:"eDelivery"`       //E-Delivery is used to indicate if the form should be delivered electronically
	CorrectedReturn     bool          `json:"correctedReturn"` //Corrected Return is used to indicate if the form is a corrected return
}

// RecipientInfo represents the details of a recipient
type RecipientInfo struct {
	ID                     int    `json:"payerId"`                     //PayerID is the unique identifier for the payer in Tax1099's system
	ClientID               string `json:"clientRecipientId,omitempty"` //ClientRecipientID is the unique identifier for the recipient in your system
	TinType                string `json:"tinType"`                     //TinType is the type of TIN for the recipient, whether business or individual
	TaxIdentifer           string `json:"recipientTin"`                //TaxIdentifier is the recipient's TIN, SSN, or EIN with no dashes.
	FirstName              string `json:"firstName,omitempty"`
	MiddleName             string `json:"middleName,omitempty"`
	LastNameOrBusinessName string `json:"lastNameOrBusinessName"`
	Suffix                 string `json:"suffix,omitempty"`
	Address                string `json:"address"`
	Address2               string `json:"address2,omitempty"`
	City                   string `json:"city"`
	State                  string `json:"state"`   //State is the two-letter abbreviation for the state
	ZipCode                string `json:"zipCode"` //ZipCode is the 5-digit zip code or 9-digit zip code with a hyphen
	Country                string `json:"country"` //Country is the two-letter abbreviation for the country
	Email                  string `json:"email"`
	RecipientID            int    `json:"recipientId"`
	AttentionTo            string `json:"attentionTo"`
	IsActive               bool   `json:"isActive"`
	Phone                  string `json:"phone"`
}

// Submit1098Response represents the response for the Submit 1098 API
type Submit1098Response struct {
	Result             []SubmissionResult `json:"result"`
	TotalCount         int                `json:"totalCount"`
	ValidationErrors   []ValidationError  `json:"validationErrors"`
	Message            string             `json:"message"`
	StatusCode         int                `json:"statusCode"`
	OriginalStatusCode int                `json:"originalStatusCode"`
	IsError            bool               `json:"isError"`
}

// SubmissionResult represents individual submission results
type SubmissionResult struct {
	ID         int  `json:"id"`
	IsInserted bool `json:"isInserted"`
}
