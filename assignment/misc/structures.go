package misc

type Address struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	City      string `json:"city"`
	Pincode   string `json:"pincode"`
	State     string `json:"state"`
	Country   string `json:"country"`
}

type UserDetailsRequest struct {
	UserId       string  `json:"userId"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Age          int     `json:"age"`
	Gender       string  `json:"gender"`
	Address      Address `json:"address"`
	Email        string  `json:"email"`
	MobileNumber string  `json:"mobileNumber"`
	AccessLevel  string  `json:"accessLevel"`
	UserName     string  `json:"username"`
}
