package dao

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrQueryingDatabase     = errors.New("error querying database")
	ErrParsingRow           = errors.New("error parsing row from database")
	ErrHashingPassword      = errors.New("error hashing password")
	ErrSavingToDatabase     = errors.New("error saving user in database")
	ErrUpdatingToDatabase   = errors.New("error updating user in database")
	ErrDeletingFromDatabase = errors.New("error deleting user in database")
)

func buildWhereClause(params map[string]string) (string, []any) {
	clauses := make([]string, 0)
	vals := make([]any, 0)
	i := 1
	for _, validParam := range validUserParams {
		_, ok := params[validParam]
		if ok {
			databaseField := paramToColumn[validParam]
			vals = append(vals, params[validParam])
			clause := fmt.Sprintf("%s = $%d", databaseField, i)
			clauses = append(clauses, clause)
			i++
		}
	}

	return strings.Join(clauses, " and "), vals
}

func Close(err *error, c io.Closer) {
	if e := c.Close(); err != nil {
		*err = e
	}
}
