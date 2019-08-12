package tracery

/** ModifierFn is provided as a convience for using plain functions as Modifiers
 */
type ModifierFn func(value string, params ...string) string

/** Modify implements the single member of Modifier for ModifierFn
 */
func (f ModifierFn) Modify(value string, params ...string) string {
	if f == nil {
		return value
	}
	return f(value, params...)
}
