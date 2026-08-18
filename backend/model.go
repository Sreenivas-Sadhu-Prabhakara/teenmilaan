package backend

import (
	"encoding/json"
	"fmt"
	"math"
)

// Match is a counter-side check of what was ordered, what the invoice says, and
// what was actually counted. The order line is optional — the core value works
// as a 2-way match (invoice vs counted) even when the verbal order wasn't logged.
type Match struct {
	Item         string  `json:"item"`
	OrderedQty   float64 `json:"orderedQty"` // optional; 0 = not logged
	InvoicedQty  float64 `json:"invoicedQty"`
	CountedQty   float64 `json:"countedQty"`
	InvoicedRate float64 `json:"invoicedRate"`
	AgreedRate   float64 `json:"agreedRate"`
}

// QtyGap is invoiced minus counted (positive = short-delivered vs invoice).
func (m Match) QtyGap() float64 { return m.InvoicedQty - m.CountedQty }

// RateGap is invoiced minus agreed rate (positive = overcharged).
func (m Match) RateGap() float64 { return m.InvoicedRate - m.AgreedRate }

// Mismatch reports whether quantity or rate disagree.
func (m Match) Mismatch() bool {
	return math.Abs(m.QtyGap()) > 1e-9 || math.Abs(m.RateGap()) > 1e-9
}

// ValueLeak is the rupee value at risk: short quantity valued at the agreed rate
// plus the rate overcharge on the quantity actually received.
func (m Match) ValueLeak() float64 {
	return m.QtyGap()*m.AgreedRate + m.RateGap()*m.CountedQty
}

// Validate reports whether the Match is well formed.
func (m Match) Validate() error {
	if m.Item == "" {
		return fmt.Errorf("item is required")
	}
	if m.InvoicedQty < 0 || m.CountedQty < 0 || m.InvoicedRate < 0 || m.AgreedRate < 0 {
		return fmt.Errorf("quantities and rates cannot be negative")
	}
	return nil
}

// Summary aggregates the match log.
type Summary struct {
	Total        int     `json:"total"`
	Mismatches   int     `json:"mismatches"`
	TotalValueLeak float64 `json:"totalValueLeak"`
}

// Summarize counts mismatches and totals the value at risk.
func Summarize(records []Record) Summary {
	var s Summary
	for _, r := range records {
		var m Match
		if json.Unmarshal(r.Input, &m) != nil {
			continue
		}
		s.Total++
		if m.Mismatch() {
			s.Mismatches++
			s.TotalValueLeak += m.ValueLeak()
		}
	}
	return s
}

// parseEntry decodes+validates a match; headline is its value leak, label its state.
func parseEntry(raw []byte) (float64, string, error) {
	var m Match
	if err := json.Unmarshal(raw, &m); err != nil {
		return 0, "", fmt.Errorf("invalid json")
	}
	if err := m.Validate(); err != nil {
		return 0, "", err
	}
	label := "match"
	if m.Mismatch() {
		label = "mismatch"
	}
	return m.ValueLeak(), label, nil
}
