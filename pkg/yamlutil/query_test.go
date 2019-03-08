package yamlutil

import "testing"

func TestQueryInto(t *testing.T) {
	str := `foo:
  - user-a
  - user-b`
	var users []string
	err := QueryInto(str, &users, "foo", 0)
	if err != nil {
		t.Errorf("error while querying: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("len(users): expected %d but got %d", 2, len(users))
	}
}
