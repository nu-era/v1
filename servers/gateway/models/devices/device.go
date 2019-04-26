package devices

import (
	"errors"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

//ErrUserNotFound is returned when the user can't be found
var ErrDeviceNotFound = errors.New("Device not found")

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

// Device represents the client for our alerting system
// this stores necessary information for sending/receiving alerts
// and other useful/identifying information
type Device struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `json:"name"` // ideally want org or business name i.e. Seattle School district
	Lat      float64       `json:"latitude"`
	Long     float64       `json:"longitude"`
	PassHash []byte        `json:"-"`
	Email    string        `json:"-"`
	Phone    string        `json:"-"`
	Status   string        `json:"status"`
}

//Credentials represents device authorization credentials
// COMMENT: Do we need this if we want persistent connections?
//// Maybe just a password that we can reference in a lookup table via hash
//// for authorized users?
type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

//NewDevice represents a new device being stood up
// COMMENT: Do we want to make credentials for new users and send via email?
type NewDevice struct {
	Name         string  `json:"name"`
	Lat          float64 `json:"latitude"`
	Long         float64 `json:"longitude"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Password     string  `json:"password"`
	PasswordConf string  `json:"passwordConf"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	Name         string  `json:"Name"`
	Lat          float64 `json:"latitude"`
	Long         float64 `json:"longitude"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Status       string  `json:"status"`
	OldPassword  string  `json:"oldPassword"`
	Password     string  `json:"password"`
	PasswordConf string  `json:"passwordConf"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewDevice) Validate() error {
	if _, err := mail.ParseAddress(nu.Email); err != nil {
		return fmt.Errorf("got parse error for email: %v", err)
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters, got %s", nu.Password)
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password must match PasswordConf, got Password: %s and PasswordConf: %s", nu.Password, nu.PasswordConf)
	}
	if len(nu.Name) < 1 {
		return fmt.Errorf("Name must have non-zero length, got %s", nu.Name)
	}
	if len(nu.Email) < 1 {
		return fmt.Errorf("Email must be provided, got %s", nu.Email)
	}
	if nu.Lat == 0 || nu.Long == 0 {
		return fmt.Errorf("Location must be provided, got lat:%f, long:%f", nu.Lat, nu.Long)
	}
	//TODO: Restrict to only washington state area

	return nil
}

//ToDevice converts the NewDevice to a Device, setting the
//PassHash, Status, and other fields appropriately
func (nu *NewDevice) ToDevice() (*Device, error) {
	if err := nu.Validate(); err != nil {
		return nil, err
	}

	dev := &Device{ // make new device
		ID:     bson.NewObjectId(),
		Email:  nu.Email,
		Name:   nu.Name,
		Lat:    nu.Lat,
		Long:   nu.Long,
		Status: "down",
	}
	// hash and set passHash field of device
	if err := dev.SetPassword(nu.Password); err != nil {
		return nil, err
	}
	return dev, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (d *Device) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	d.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (d *Device) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword(d.PassHash, []byte(password)); err != nil {
		return err
	}
	return nil
}

//ApplyUpdates applies the updates to the device. An error
//is returned if the updates are invalid
func (d *Device) ApplyUpdates(updates *Updates) error {
	if updates == nil { // nil struct passed
		return errors.New("Update struct must contain desired updates to fields")
	}
	if updates.Email != "" {
		if _, err := mail.ParseAddress(updates.Email); err != nil {
			return fmt.Errorf("got parse error for email: %v", err)
		} else {
			d.Email = updates.Email
		}
	}
	if updates.Name != "" {
		d.Name = updates.Name
	}
	if updates.Lat != 0 {
		d.Lat = updates.Lat
	}
	if updates.Long != 0 {
		d.Long = updates.Long
	}
	if updates.Status != "" && (updates.Status == "up" || updates.Status == "down") {
		d.Status = updates.Status
	}
	if len(updates.Phone) > 10 && len(updates.Phone) < 12 {
		d.Phone = updates.Phone
	} else {
		return errors.New("Not a valid phone number")
	}
	if len(updates.OldPassword) != 0 {
		if err := d.Authenticate(updates.OldPassword); err != nil {
			return err
		}
		if len(updates.Password) < 6 {
			return errors.New("New password must be at least 6 characters")
		} else if len(updates.Password) != 0 && updates.Password != updates.PasswordConf {
			return errors.New("Passwords do not match")
			if err := d.SetPassword(updates.Password); err != nil {
				return err
			}
		}
	}
	return nil
}
