package filters

import "github.com/asaskevich/govalidator"

type Validator struct {
}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
