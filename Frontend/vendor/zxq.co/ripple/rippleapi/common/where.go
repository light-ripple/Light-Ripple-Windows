package common

// WhereClause is a struct representing a where clause.
// This is made to easily create WHERE clauses from parameters passed from a request.
type WhereClause struct {
	Clause string
	Params []interface{}
	useOr  bool
}

// Where adds a new WHERE clause to the WhereClause.
func (w *WhereClause) Where(clause, passedParam string, allowedValues ...string) *WhereClause {
	if passedParam == "" {
		return w
	}
	if len(allowedValues) != 0 && !contains(allowedValues, passedParam) {
		return w
	}
	w.addWhere()
	w.Clause += clause
	w.Params = append(w.Params, passedParam)
	return w
}

func (w *WhereClause) addWhere() {
	// if string is empty add "WHERE", else add AND
	if w.Clause == "" {
		w.Clause += "WHERE "
	} else {
		if w.useOr {
			w.Clause += " OR "
			return
		}
		w.Clause += " AND "
	}
}

// Or enables using OR instead of AND
func (w *WhereClause) Or() *WhereClause {
	w.useOr = true
	return w
}

// And enables using AND instead of OR
func (w *WhereClause) And() *WhereClause {
	w.useOr = false
	return w
}

// In generates an IN clause.
// initial is the initial part, e.g. "users.id".
// Fields are the possible values.
// Sample output: users.id IN ('1', '2', '3')
func (w *WhereClause) In(initial string, fields ...[]byte) *WhereClause {
	if len(fields) == 0 {
		return w
	}
	w.addWhere()
	w.Clause += initial + " IN (" + generateQuestionMarks(len(fields)) + ")"
	fieldsInterfaced := make([]interface{}, len(fields))
	for k, f := range fields {
		fieldsInterfaced[k] = string(f)
	}
	w.Params = append(w.Params, fieldsInterfaced...)
	return w
}

func generateQuestionMarks(x int) (qm string) {
	for i := 0; i < x-1; i++ {
		qm += "?, "
	}
	if x > 0 {
		qm += "?"
	}
	return qm
}

// ClauseSafe returns the clause, always containing something. If w.Clause is
// empty, it returns "WHERE 1".
func (w *WhereClause) ClauseSafe() string {
	if w.Clause == "" {
		return "WHERE 1"
	}
	return w.Clause
}

// Where is the same as WhereClause.Where, but creates a new WhereClause.
func Where(clause, passedParam string, allowedValues ...string) *WhereClause {
	w := new(WhereClause)
	return w.Where(clause, passedParam, allowedValues...)
}
