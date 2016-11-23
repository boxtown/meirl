package postgres

import "testing"

func TestSeekingPaginationBuild(t *testing.T) {
	query := ""
	p := paginator{}

	p.field = "test"
	result := p.seekingQuery(query, 0, false)
	expected := " WHERE test > $0 ORDER BY test ASC LIMIT 0"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}

	p.desc = true
	result = p.seekingQuery(query, 0, false)
	expected = " WHERE test < $0 ORDER BY test DESC LIMIT 0"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}

	result = p.seekingQuery(query, 0, true)
	expected = " AND test < $0 ORDER BY test DESC LIMIT 0"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}

	result = p.seekingQuery(query, 2, true)
	expected = " AND test < $2 ORDER BY test DESC LIMIT 0"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}

	p.limit = 10
	result = p.seekingQuery(query, 2, true)
	expected = " AND test < $2 ORDER BY test DESC LIMIT 10"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}

	query = "SELECT FROM posts WHERE id=$1"
	result = p.seekingQuery(query, 2, true)
	expected = "SELECT FROM posts WHERE id=$1 AND test < $2 ORDER BY test DESC LIMIT 10"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		t.Fail()
	}
}
