package tfprofile

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
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

// Print a parsed log in tabular format, optionally sorting by certain columns
// sort_spec is a comma-separated list of "column_name=(asc|desc)", e.g. "n=asc,tot_time=desc"
func Table(log ParsedLog, sort_spec string) error {
	headerFmt := color.New(color.FgHiBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgBlue).SprintfFunc()

	tbl := table.New("resource", "n", "tot_time", "idx_creation", "idx_created")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Sort the resources according to the sort_spec and create rows
	for _, r := range Sort(log, sort_spec) {
		for resource, metric := range log {
			if r == resource {
				tbl.AddRow(
					resource,
					(metric.NumCalls),
					(metric.TotalTime),
					(metric.CreationIndex),
					(metric.CreatedIndex),
				)
				break
			}
		}
	}

	fmt.Println() // Create space above the table
	tbl.Print()

	return nil
}

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

	for k, v := range log {
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
				value = float64(v.CreationIndex)
			} else if column == "idx_created" {
				value = float64(v.CreatedIndex)
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
