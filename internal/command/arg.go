package command

type ValidatedArgs struct {
	values map[string]interface{}
}

func NewValidatedArgs() *ValidatedArgs {
	return &ValidatedArgs{values: make(map[string]interface{})}
}

func (v *ValidatedArgs) String(label string) string {
	return v.values[label].(string)
}

func (v *ValidatedArgs) Int(label string) int {
	return v.values[label].(int)
}

func (v *ValidatedArgs) Bool(label string) bool {
	return v.values[label].(bool)
}

func (v *ValidatedArgs) StringOr(label string, defaultValue string) string {
	if val, exists := v.values[label]; exists {
		return val.(string)
	}

	return defaultValue
}

func (v *ValidatedArgs) IntOr(label string, defaultValue int) int {
	if val, exists := v.values[label]; exists {
		return val.(int)
	}

	return defaultValue
}

func (v *ValidatedArgs) BoolOr(label string, defaultValue bool) bool {
	if val, exists := v.values[label]; exists {
		return val.(bool)
	}

	return defaultValue
}

type Args []Arg
type Arg interface {
	Label() string
	Description() string
	Required() bool
}

type StringArg struct {
	label       string
	description string
	required    bool
}

func NewStringArg(label, description string, required bool) *StringArg {
	return &StringArg{
		label:       label,
		description: description,
		required:    required,
	}
}

func (a *StringArg) Label() string {
	return a.label
}

func (a *StringArg) Description() string {
	return a.description
}

func (a *StringArg) Required() bool {
	return a.required
}

type IntArg struct {
	label       string
	description string
	required    bool
}

func NewIntArg(label, description string, required bool) *IntArg {
	return &IntArg{
		label:       label,
		description: description,
		required:    required,
	}
}

func (a *IntArg) Label() string {
	return a.label
}

func (a *IntArg) Description() string {
	return a.description
}

func (a *IntArg) Required() bool {
	return a.required
}

type BoolArg struct {
	label       string
	description string
	required    bool
}

func NewBoolArg(label, description string, required bool) *BoolArg {
	return &BoolArg{
		label:       label,
		description: description,
		required:    required,
	}
}

func (a *BoolArg) Label() string {
	return a.label
}

func (a *BoolArg) Description() string {
	return a.description
}

func (a *BoolArg) Required() bool {
	return a.required
}
