package optparser

import "testing"

func TestParseOpts(t *testing.T) {
	pos, kv := ParseOpts("foo,bar,key=value,ab=cd", true)

	if len(pos) != 2 {
		t.Errorf("len(pos): expected %d but got %d", 2, len(pos))
	}

	if pos[0] != "foo" {
		t.Errorf("pos[0]: expected %q but got %q", "foo", pos[0])
	}

	if pos[1] != "bar" {
		t.Errorf("pos[0]: expected %q but got %q", "bar", pos[0])
	}

	if len(kv) != 2 {
		t.Errorf("len(kv): expected %d but got %d", 2, len(kv))
	}

	if kv["key"] != "value" {
		t.Errorf("kv[key]: expected %q but got %q", "value", kv["key"])
	}

	if kv["ab"] != "cd" {
		t.Errorf("kv[ab]: expected %q but got %q", "cd", kv["ab"])
	}
}

func TestEmptyParseCLIOpts(t *testing.T) {
	pos, kv := ParseOpts(",", true)

	if pos != nil {
		t.Errorf("pos: expected %v but got %+v", nil, pos)
	}

	if kv != nil {
		t.Errorf("kv: expected %v but got %+v", nil, kv)
	}
}
