package checker

import (
	"errors"
	"fmt"

	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"
)

// The assumption is that the input to the checker is a simplified, unified
// policy that is ready to be compiled. The checker does the final step of making sure that
// all the rules are valid and type checks.
// The checker will not do anything with the macros defined.
// It will assume all calls and variable references left are errors (but that should have been caught
// in the phases before).
// Except for checking type validity, the checker will also make sure we don't have
// more than one rule for the same syscall. This is also the place where we make sure
// all the syscalls with rules are defined.

// EnsureValid takes a policy and returns all the errors encounterered for the given rules
// If everything is valid, the return will be empty
func EnsureValid(p tree.Policy) []error {
	v := &validityChecker{rules: p.Rules, seen: make(map[string]bool)}
	return v.check()
}

type validityChecker struct {
	rules []tree.Rule
	seen  map[string]bool
}

type ruleError struct {
	syscallName string
	err         error
}

func (e *ruleError) Error() string {
	return fmt.Sprintf("[%s] %s", e.syscallName, e.err)
}

func checkValidSyscall(r tree.Rule) error {

	if _, ok := constants.GetSyscall(r.Name); !ok {
		return errors.New("invalid syscall")
	}
	return nil
}

func (v *validityChecker) check() []error {
	result := []error{}

	for _, r := range v.rules {
		var res error
		if v.seen[r.Name] {
			res = errors.New("duplicate definition of syscall rule")
		}
		v.seen[r.Name] = true
		if res == nil {
			res = checkValidSyscall(r)
		}
		if res == nil {
			res = v.checkRule(r)
		}
		if res != nil {
			result = append(result, &ruleError{syscallName: r.Name, err: res})
		}
	}

	return result
}

func (v *validityChecker) checkRule(r tree.Rule) error {
	return typeCheckExpectingBoolean(r.Body)
}
