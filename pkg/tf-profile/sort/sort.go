package tfprofile

import (
	"reflect"
	"sort"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

type (
	SortSpecItem struct {
		col   string
		order string
	}

	// Fake record we construct to allow sorting on (multiple) custom columns.
	// See Sort() for usage.
	ProxyRecord struct {
		resource string
		items    []float64
	}
)

// Parse a sort_spec into a map
// e.g "n=asc,tot_time=desc" => {n: asc, tot_time: desc}
func parseSortSpec(in string) []SortSpecItem {
	tokens := strings.Split(in, ",")

	result := []SortSpecItem{}
	for _, spec := range tokens {
		split := strings.Split(spec, "=")
		result = append(result, SortSpecItem{split[0], split[1]})
	}
	return result
}

// Sort a parsed log according to the provided sort_spec
func Sort(log ParsedLog, sort_spec string) []string {
	// Because we can not construct a custom sort function upfront,
	// we "rebuild" the log such that the "sorting" metrics come first,
	// and values for columns that are to be sorted descendingly are
	// inverted. This way, the sorting function is always the same
	proxy_log := []ProxyRecord{}

	sort_spec_p := parseSortSpec(sort_spec)

	for k, v := range log.Resources {
		proxy_item_values := []float64{0, 0, 0, 0}

		// With values in the order of the sort_spec, create a proxy record
		for idx, sort_item := range sort_spec_p {
			column := sort_item.col
			order := sort_item.order
			var value float64
			if column == "n" {
				value = float64(v.NumCalls)
			} else if column == "tot_time" {
				value = float64(v.TotalTime)
			} else if column == "idx_creation" {
				value = float64(v.CreationStartedIndex)
			} else if column == "idx_created" {
				value = float64(v.CreationCompletedIndex)
			} else if column == "status" {
				value = float64(v.AfterStatus)
			}
			if order == "desc" {
				value = -value
			}
			proxy_item_values[idx] = value
		}

		proxy_log = append(proxy_log, ProxyRecord{k, proxy_item_values})
	}

	N := reflect.TypeOf(ProxyRecord{}).NumField()

	// Sort the proxy log
	sort.Slice(proxy_log, func(i, j int) bool {
		// Custom sort function: sort by all values in 'items'
		for item := 0; item < N; item++ {
			if proxy_log[i].items[item] != proxy_log[j].items[item] {
				return proxy_log[i].items[item] < proxy_log[j].items[item]
			}
		}
		return false // Everything is equal
	})

	// Finally, extract the resource names out of the sorted slice
	result := []string{}
	for _, v := range proxy_log {
		result = append(result, v.resource)
	}
	return result
}
