package broker

import (
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	b, err := New(10)
	errIsNil(t, err)

	errIsNil(t, b.Put("pet", "cat"))
	errIsNil(t, b.Put("pet", "dog"))
	errIsNil(t, b.Put("role", "manager"))
	errIsNil(t, b.Put("role", "executive"))

	assertGetEq(t, b, "pet", 0, "cat", true, "")
	assertGetEq(t, b, "pet", 0, "dog", true, "")
	assertGetEq(t, b, "pet", 1, "", false, "")
	assertGetEq(t, b, "pet", 1, "", false, "")
	assertGetEq(t, b, "role", 0, "manager", true, "")
	assertGetEq(t, b, "role", 0, "executive", true, "")
	assertGetEq(t, b, "role", 1, "", false, "")
}

func errIsNil(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func get(t *testing.T, b *Broker, qName string, timeout int) (string, bool) {
	msgs, err := b.Subscribe(qName)
	errIsNil(t, err)

	// Ждем сообщение до конца
	if timeout == 0 {
		return <-msgs, true
	}

	// Ждем сообщение некоторое время
	tmr := time.NewTimer(time.Duration(timeout) * time.Second)
	defer tmr.Stop()

	select {
	case <-tmr.C:
		return "", false

	case msg := <-msgs:
		return msg, true
	}
}

func assertGetEq(t *testing.T, b *Broker, qName string, timeout int, expV string, expOK bool, msg string) {
	if v, ok := get(t, b, qName, timeout); ok != expOK || v != expV {
		t.Errorf("assertGetEq: %s (actual: %s, %t)", msg, v, ok)
	}
}
