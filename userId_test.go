package chatwork

import (
	"testing"
)

func TestToString(t *testing.T) {
	var ids UserIds
	ids = []UserId{UserId(1), UserId(2), UserId(3)}
	result := ids.toString(",")
	if result != "1,2,3" {
		t.Errorf("'%s' is not '1,2,3'", result)
	}
}
