package captcha

import (
	"testing"
	"sort"
)

func Test(t *testing.T) {
	results := [][]*Result{}
	results = append(results, []*Result{NewResult("a", 1.0), NewResult("b", 0.5)})
	results = append(results, []*Result{NewResult("c", 0.9), NewResult("d", 0.3)})
	var c ResultSorter
	c = GenerateResults(results)
	sort.Sort(c)
	if c[0].Label != "ac" {
		t.Error()
	}
	if c[1].Label != "bc" {
		t.Error()
	}
	if c[2].Label != "ad" {
		t.Error()
	}
	if c[3].Label != "bd" {
		t.Error()
	}
}
