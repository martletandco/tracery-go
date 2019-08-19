package tracery

/** ModifierFunc is provided as a convience for using plain functions as Modifiers
* Although most of the time you'll likely want `Grammar.AddModifyFunc`
 */
type ModifierFunc func(value string, params ...string) string

/** Modify implements the single member of Modifier for ModifierFn
 */
func (f ModifierFunc) Modify(value string, params ...string) string {
	if f == nil {
		return value
	}
	return f(value, params...)
}
