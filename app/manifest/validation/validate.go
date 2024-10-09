package validation

type ValidateResult struct {
	errors []string
}

func (v ValidateResult) AddError(err string) {
	v.errors = append(v.errors, err)
}

func (v ValidateResult) HasError() bool {
	return len(v.errors) > 0
}

func (v ValidateResult) IsSuccess() bool {
	return !v.HasError()
}
