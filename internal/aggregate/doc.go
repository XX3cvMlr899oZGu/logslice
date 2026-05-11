// Package aggregate counts distinct values for a named log field across
// the lines processed by the slicer pipeline.
//
// Usage:
//
//	c := aggregate.New("level")
//	c.Add("info")
//	c.Add("error")
//	c.Add("info")
//	c.WriteTo(os.Stdout)
//
// Output:
//
//	field: level  total: 3
//	  info                            2  ( 66.7%)
//	  error                           1  ( 33.3%)
package aggregate
