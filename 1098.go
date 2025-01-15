package tax1099

import "log"

// Submit1098Request represents the JSON structure for submitting 1098 forms
type Submit1098Request struct {
	TaxYear string     `json:"taxYear"`
	Items   []Item1098 `json:"items"`
}

// Item represents a single payer and their associated forms
type Item1098 struct {
	PayerInfo PayerInfo  `json:"payerInfo"` //Payer Info is the seller/lender of the property
	Forms     []Form1098 `json:"forms"`
}

// Form1098 represents the details of a single 1098 form
type Form1098 struct {
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
	PropertyDescription string        `json:"propertyDescription,omitempty"`
	USPSMail            bool          `json:"uspsMail"`        //USPS Mail is used to indicate if the form should be mailed to the payer
	TINCheck            bool          `json:"tinCheck"`        //TIN Check is used to indicate if the TIN should be checked
	EDelivery           bool          `json:"eDelivery"`       //E-Delivery is used to indicate if the form should be delivered electronically
	CorrectedReturn     bool          `json:"correctedReturn"` //Corrected Return is used to indicate if the form is a corrected return
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

func (t *tax1099Impl) Validate1098(payload Submit1098Request) (Submit1098Response, error) {
	log.Println("Submitting the 1098 form for validation...")

	var res Submit1098Response
	if err := t.post(t.generateFullUrl(Url1098, "form/1098/validate"), payload, &res); err != nil {
		return res, err
	}

	log.Printf("Validation response: %+v", res)

	return res, nil
}
