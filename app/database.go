package app

import db "github.com/sonyarouje/simdb"

type Customer struct {
	CustID  string `json:"custid"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact Contact
}

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// ID any struct that needs to persist should implement this function defined
// in Entity interface.
func (c Customer) ID() (jsonField string, value interface{}) {
	value = c.CustID
	jsonField = "custid"
	return
}

var Driver *db.Driver
var err error

func Database() {
	Driver, err = db.New("data")
	if err != nil {
		panic(err)
	}
}

func crud() {
	customer := Customer{
		CustID:  "CUST1",
		Name:    "sarouje",
		Address: "address",
		Contact: Contact{
			Phone: "45533355",
			Email: "someone@gmail.com",
		},
	}

	//creates a new Customer file inside the directory passed as the parameter to New()
	//if the Customer file already exist then insert operation will add the customer data to the array
	err = Driver.Insert(customer)
	if err != nil {
		panic(err)
	}

	//GET ALL Customer
	//opens the customer json file and filter all the customers with name sarouje.
	//AsEntity takes a pointer to Customer array and fills the result to it.
	//we can loop through the customers array and retireve the data.
	var customers []Customer
	err = Driver.Open(Customer{}).Where("name", "=", "sarouje").Get().AsEntity(&customers)
	if err != nil {
		panic(err)
	}

	//GET ONE Customer
	//First() will return the first record from the results
	//AsEntity takes a pointer to Customer variable (not an array pointer)
	var customerFrist Customer
	err = Driver.Open(Customer{}).Where("custid", "=", "CUST1").First().AsEntity(&customerFrist)
	if err != nil {
		panic(err)
	}

	//Update function uses the ID() to get the Id field/value to find the record and update the data.
	customerFrist.Name = "Sony Arouje"
	err = Driver.Update(customerFrist)
	if err != nil {
		panic(err)
	}

	//Delete
	toDel := Customer{
		CustID: "CUST1",
	}
	err = Driver.Delete(toDel)
}