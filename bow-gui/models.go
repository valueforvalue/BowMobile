package main

const Version = "1.1.1"

type Part struct {
	KeyNo       string
	BasePart    string
	Revision    string
	Qty         string
	Description string
	Remarks     string
}

type PartOccurrence struct {
	ModelSeries     string
	ManualRevision  string
	FigureID        string
	KeyNo           string
	FullPartNumber  string
	Revision        string
	Description     string
	Remarks         string
}

type GroupedResult struct {
	BasePart    string
	Description string
	Occurrences []PartOccurrence
}

type ManualInfo struct {
	ModelSeries string
	Revision    string
}
