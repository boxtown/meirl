package api

import "testing"
import "net/http"

func TestListOptionsFromRequest(t *testing.T) {
	r, _ := http.NewRequest("", "/test", nil)
	options := ListOptionsFromRequest(r)
	if options.Offset != 0 {
		t.Errorf("Expected offset 0, parsed %d", options.Offset)
		t.Fail()
	}
	if options.Limit != 10 {
		t.Errorf("Expected limit 10, parsed %d", options.Limit)
		t.Fail()
	}
	if options.Desc {
		t.Error("Expected desc to be false, got true")
		t.Fail()
	}
	if options.Marker != "" {
		t.Error("Marker should have been empty")
		t.Fail()
	}

	r, _ = http.NewRequest("", "/test?offset=-3&limit=0", nil)
	options = ListOptionsFromRequest(r)
	if options.Offset != 0 {
		t.Errorf("Expected offset 0, parsed %d", options.Offset)
		t.Fail()
	}
	if options.Limit != 10 {
		t.Errorf("Expected limit 10, parsed %d", options.Limit)
		t.Fail()
	}
	if options.Desc {
		t.Error("Expected desc to be false, got true")
		t.Fail()
	}
	if options.Marker != "" {
		t.Error("Marker should have been empty")
		t.Fail()
	}

	r, _ = http.NewRequest("", "/test?offset=3&limit=5&desc=true&marker=01-01-2016%2000:00:00", nil)
	options = ListOptionsFromRequest(r)
	if options.Offset != 3 {
		t.Errorf("Expected offset 3, parsed %d", options.Offset)
		t.Fail()
	}
	if options.Limit != 5 {
		t.Errorf("Expected limit 5, parsed %d", options.Limit)
		t.Fail()
	}
	if !options.Desc {
		t.Error("Expected desc to be true, got false")
		t.Fail()
	}
	if options.Marker != "01-01-2016 00:00:00" {
		t.Errorf("Expected marker 01-01-2016 00:00:00, got %s", options.Marker)
		t.Fail()
	}
}

func TestAuth(t *testing.T) {
	auth := NewAuth()
	password := "test"
	secured, err := auth.SecurePassword(password)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	ok := auth.CheckPassword(password, secured)
	if !ok {
		t.Error("Secured password was not checked successfully")
		t.Fail()
	}
}
