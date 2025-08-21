package validation

type validable interface {
	Valid() ValidateResult
}

type ValidateResult struct {
	errors []string
}

func (v ValidateResult) GetErrors() []string {
	return v.errors
}

func (v *ValidateResult) AddError(err string) {
	v.errors = append(v.errors, err)
}

func (v *ValidateResult) AddErrors(errs []error) {
	for _, err := range errs {
		v.errors = append(v.errors, err.Error())
	}
}

func (v ValidateResult) HasError() bool {
	return len(v.errors) > 0
}

func (v ValidateResult) IsSuccess() bool {
	return !v.HasError()
}

func (v *ValidateResult) Append(validate ValidateResult) {
	for _, err := range validate.errors {
		v.AddError(err)
	}
}

func (v *ValidateResult) AppendValidable(validable validable) {
	var result = validable.Valid()
	v.Append(result)
}
