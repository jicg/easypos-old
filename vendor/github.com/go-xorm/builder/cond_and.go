package builder

import "fmt"

type condAnd []Cond

var _ Cond = condAnd{}

func And(conds ...Cond) Cond {
	var result = make(condAnd, 0, len(conds))
	for _, cond := range conds {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}
	return result
}

func (and condAnd) WriteTo(w Writer) error {
	for i, cond := range and {
		if _, ok := cond.(condOr); ok {
			fmt.Fprint(w, "(")
		}

		err := cond.WriteTo(w)
		if err != nil {
			return err
		}

		if _, ok := cond.(condOr); ok {
			fmt.Fprint(w, ")")
		}

		if i != len(and)-1 {
			fmt.Fprint(w, " AND ")
		}
	}

	return nil
}

func (and condAnd) And(conds ...Cond) Cond {
	return And(and, And(conds...))
}

func (and condAnd) Or(conds ...Cond) Cond {
	return Or(and, Or(conds...))
}

func (and condAnd) IsValid() bool {
	return len(and) > 0
}
