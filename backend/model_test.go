package backend

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memStore struct{ items []Record }

func (m *memStore) Save(r Record) (Record, error) {
	r.ID = int64(len(m.items) + 1)
	m.items = append([]Record{r}, m.items...)
	return r, nil
}
func (m *memStore) List(limit int) ([]Record, error) { return m.items, nil }

func TestMismatchAndValueLeak(t *testing.T) {
	// invoiced 100 @ ₹12, counted 95 @ agreed ₹11: short 5 units + ₹1 overcharge.
	m := Match{Item: "bolt", InvoicedQty: 100, CountedQty: 95, InvoicedRate: 12, AgreedRate: 11}
	if !m.Mismatch() {
		t.Fatal("expected mismatch")
	}
	// leak = 5*11 + 1*95 = 55 + 95 = 150
	if math.Abs(m.ValueLeak()-150) > 1e-9 {
		t.Fatalf("valueLeak=%v want 150", m.ValueLeak())
	}
}

func TestCleanMatch(t *testing.T) {
	m := Match{Item: "x", InvoicedQty: 10, CountedQty: 10, InvoicedRate: 5, AgreedRate: 5}
	if m.Mismatch() {
		t.Fatal("clean match flagged")
	}
}

func TestSummarize(t *testing.T) {
	in1, _ := json.Marshal(Match{Item: "a", InvoicedQty: 100, CountedQty: 95, InvoicedRate: 12, AgreedRate: 11})
	in2, _ := json.Marshal(Match{Item: "b", InvoicedQty: 10, CountedQty: 10, InvoicedRate: 5, AgreedRate: 5})
	s := Summarize([]Record{{Input: in1, Label: "mismatch"}, {Input: in2, Label: "match"}})
	if s.Total != 2 || s.Mismatches != 1 {
		t.Fatalf("summary=%+v", s)
	}
}

func TestLogEndpoint(t *testing.T) {
	srv := NewServer(&memStore{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log",
		strings.NewReader(`{"item":"bolt","invoicedQty":100,"countedQty":95,"invoicedRate":12,"agreedRate":11}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("log %d", rec.Code)
	}
}
