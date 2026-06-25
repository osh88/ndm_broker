// go test -v ./...

package queue

import "testing"

func TestQueue(t *testing.T) {
	q, err := NewQueue[string](3)
	if err != nil {
		t.Error(err)
	}

	q.Put("1")
	q.Put("2")
	q.Put("3")

	assertGetEq(t, q, "1", true, "v != 1")
	assertGetEq(t, q, "2", true, "v != 2")

	q.Put("4")
	q.Put("5")

	assertGetEq(t, q, "3", true, "v != 3")
	assertGetEq(t, q, "4", true, "v != 4")
	assertGetEq(t, q, "5", true, "v != 5")
	assertGetEq(t, q, "", false, "v != ''")

	q.Put("1")
	q.Put("2")
	q.Put("3")
	q.Put("4")
	q.Put("5")
	q.Put("6")

	assertGetEq(t, q, "1", true, "v != 1")
	assertGetEq(t, q, "2", true, "v != 2")
	assertGetEq(t, q, "3", true, "v != 3")
	assertGetEq(t, q, "4", true, "v != 4")
	assertGetEq(t, q, "5", true, "v != 5")
	assertGetEq(t, q, "6", true, "v != 6")
	assertGetEq(t, q, "", false, "v != ''")
	assertGetEq(t, q, "", false, "v != ''")
}

func assertGetEq(t *testing.T, q *Queue[string], expV string, expOK bool, msg string) {
	if v, ok := q.Get(); ok != expOK || v != expV {
		t.Errorf("assertGetEq: %s (actual: %s, %t)", msg, v, ok)
	}
}
