package postgres

import (
	"bytes"
	"strconv"
)

// seek-based paginator
type paginator struct {
	field string
	desc  bool
	limit int
}

// Build builds the paginated query from the given query
func (p *paginator) seekingQuery(query string, pIdx int, skipWhere bool) string {
	buf := bytes.NewBufferString(query)
	if skipWhere {
		buf.WriteString(" AND ")
	} else {
		buf.WriteString(" WHERE ")
	}
	buf.WriteString(p.field)
	if p.desc {
		buf.WriteString(" < ")
	} else {
		buf.WriteString(" > ")
	}
	buf.WriteRune('$')
	buf.WriteString(strconv.Itoa(pIdx))
	buf.WriteString(" ORDER BY ")
	buf.WriteString(p.field)
	if p.desc {
		buf.WriteString(" DESC LIMIT ")
	} else {
		buf.WriteString(" ASC LIMIT ")
	}
	buf.WriteString(strconv.Itoa(p.limit))
	return buf.String()
}
