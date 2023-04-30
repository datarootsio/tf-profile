package tfprofile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

func Graph(args []string, w int, h int, OutFile string) error {
	var file *bufio.Scanner
	var err error

	if len(args) == 1 {
		file, err = FileReader{File: args[0]}.Read()
	} else {
		file, err = StdinReader{}.Read()
	}
	if err != nil {
		return err
	}
	tflog, err := Parse(file, false)
	if err != nil {
		return err
	}
	tflog, err = Aggregate(tflog)
	if err != nil {
		return err
	}

	CleanFailedResources(tflog)
	_, err = PrintGNUPlotOutput(tflog, w, h, OutFile)

	return nil
}

// For failed resources, CreationCompletedEvent will always be -1, since we never
// detect the end of their modifications. We manually set their CreationCompletedEvent
// to the maximum value, leading to a long red bar.
func CleanFailedResources(tflog ParsedLog) {
	max := 0

	// Find max creation value
	for _, metrics := range tflog.Resources {
		if metrics.CreationCompletedEvent > max {
			max = metrics.CreationCompletedEvent
		}
	}

	// Update all non-successful resources to end at that index
	for resource, metrics := range tflog.Resources {
		if metrics.CreationStatus != Created && metrics.CreationStatus != AllCreated {
			metrics.CreationCompletedEvent = max
			tflog.Resources[resource] = metrics
		}
	}
}

// Use plot.tpl and a ParsedLog to generate all output for gnuplot.
// This can be piped into gnuplot (optionally providing a filename at runtime)
func PrintGNUPlotOutput(tflog ParsedLog, w int, h int, OutFile string) (string, error) {
	// Context object for templating
	Context := map[string]interface{}{}
	Context["W"] = w
	Context["H"] = h
	Context["File"] = OutFile

	SortedResources := sortResourcesForGraph(tflog)
	Resources := []string{} // Lines passed into template

	// Build list of lines and let template do the looping
	for _, r := range SortedResources {
		metrics := tflog.Resources[r]

		NameForOutput := strings.Replace(r, "_", `\\\_`, -1)
		NameForOutput = strings.Replace(NameForOutput, `"`, `'`, -1)
		// Escape underscores and add the necessary metrics.
		line := fmt.Sprintf("%v %v %v %v",
			NameForOutput,
			metrics.CreationStartedEvent,
			metrics.CreationCompletedEvent,
			metrics.CreationStatus,
		)
		Resources = append(Resources, line)
	}
	Context["Resources"] = Resources

	template, _ := template.New("plot").Parse(Template)
	err := template.Execute(os.Stdout, Context) // To stdout
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = template.Execute(&output, Context) // To variable
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

// To create a nice graph, sort the resources chronologically
// according to CreationStartedEvent
func sortResourcesForGraph(log ParsedLog) []string {
	// Collect keys
	keys := []string{}
	for key := range log.Resources {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return log.Resources[keys[i]].CreationStartedEvent > log.Resources[keys[j]].CreationStartedEvent
	})
	return keys
}

const Template string = `
# GNUplot template for generating Gantt chart. $DATA will be provided at runtime
reset
set termoption dash
set terminal pngcairo  background "#ffffff" fontscale 1.0 dashed size {{ .W }}, {{ .H }}

# --- Output colors
green = 0x49A720;# 0xFFE599;
red = 0xD32F2F; # 0xF1C232;

# resource        start    end   status
$DATA << EOD 
{{range .Resources -}} 	
{{ . }}
{{ end }}
EOD     
                     
# set output
set output "{{ .File }}"

# grid and tics
set mxtics 
set mytics
set grid xtics
set grid ytics
set grid mxtics

# create list of keys
List = ''
set table $Dummy
    plot $DATA u (List=List.'"'.strcol(1).'" ',NaN) w table
unset table

# define functions for lookup/index and color
Lookup(s) = (Index = NaN, sum [i=1:words(List)] \
    (Index = s eq word(List,i) ? i : Index,0), Index)
Color(s) = (s eq "Created" || s eq "AllCreated") ?  green : red

# set range of x-axis and y-axis
set xrange [-1:]
set yrange [0.5:words(List)+0.5]

set label "(All)Created" at screen 0.86,0.93 tc rgb green
set label "Other" at screen 0.86,0.89 tc rgb red

plot $DATA u 2:(Idx=Lookup(strcol(1))): 3 : 2 :(Idx-0.2):(Idx+0.2): \
    (Color(strcol(4))): ytic(strcol(1)) w boxxyerror fill solid 0.7 lw 2.0 lc rgb var notitle`
