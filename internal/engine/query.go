package engine

import (
	"fmt"
	"log"
	"time"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
	"github.com/influxdata/influx-stress/internal/report"
	"github.com/influxdata/influxdb1-client/models"
)

// QueryStatement is a Statement Implementation to run queries on the target InfluxDB instance
type QueryStatement struct {
	StatementID string
	Name        string

	// TemplateString is a query template that can be filled in by Args
	TemplateString string
	Args           []string

	// Number of queries to run
	Count int

	// Tracer for tracking returns
	Tracer *influxclient.Tracer

	// track time for all queries
	runtime time.Duration
}

// This function adds tags to the recording points
func (i *QueryStatement) tags() map[string]string {
	tags := make(map[string]string)
	return tags
}

// SetID statisfies the Statement Interface
func (i *QueryStatement) SetID(s string) {
	i.StatementID = s
}

// Run statisfies the Statement Interface
func (i *QueryStatement) Run(s *influxclient.StressTest) {

	i.Tracer = influxclient.NewTracer(i.tags())

	vals := make(map[string]interface{})

	var point models.Point

	runtime := time.Now()

	for j := 0; j < i.Count; j++ {

		// If the query is a simple query, send it.
		if len(i.Args) == 0 {
			b := []byte(i.TemplateString)

			// Make the package
			p := influxclient.NewPackage(influxclient.Query, b, i.StatementID, i.Tracer)

			// Increment the tracer
			i.Tracer.Add(1)

			// Send the package
			s.SendPackage(p)

		} else {
			// Otherwise cherry pick field values from the commune?

			// TODO: Currently the program lock up here if s.GetPoint
			//       cannot return a value, which can happen.
			// See insert.go
			s.Lock()
			point = s.GetPoint(i.Name, s.Precision)
			s.Unlock()

			setMapValues(vals, point)

			// Set the template string with args from the commune
			b := []byte(fmt.Sprintf(i.TemplateString, setArgs(vals, i.Args)...))

			// Make the package
			p := influxclient.NewPackage(influxclient.Query, b, i.StatementID, i.Tracer)

			// Increment the tracer
			i.Tracer.Add(1)

			// Send the package
			s.SendPackage(p)

		}
	}

	// Wait for all operations to finish
	i.Tracer.Wait()

	// Stop time timer
	i.runtime = time.Since(runtime)
}

// Report statisfies the Statement Interface
func (i *QueryStatement) Report(s *influxclient.StressTest) string {
	// Pull data via StressTest client
	allData := s.GetStatementResults(i.StatementID, "query")

	if len(allData) == 0 || allData[0].Series == nil {
		log.Fatalf("No data returned for query report\n  Statement Name: %v\n  Statement ID: %v\n", i.Name, i.StatementID)
	}

	return report.Query(i.Name, allData[0].Series[0].Columns, allData[0].Series[0].Values)
}

func getRandomTagPair(m models.Tags) string {
	for k, v := range m {
		return fmt.Sprintf("%v='%v'", k, v)
	}

	return ""
}

func getRandomFieldKey(m map[string]interface{}) string {
	for k := range m {
		return fmt.Sprintf("%v", k)
	}

	return ""
}

func setMapValues(m map[string]interface{}, p models.Point) {
	fields, err := p.Fields()
	if err != nil {
		panic(err)
	}
	m["%f"] = getRandomFieldKey(fields)
	m["%m"] = string(p.Name())
	m["%t"] = getRandomTagPair(p.Tags())
	m["%a"] = p.UnixNano()
}

func setArgs(m map[string]interface{}, args []string) []interface{} {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		values[i] = m[arg]
	}
	return values
}
