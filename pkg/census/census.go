package census

// The census command is so simple that we just put the implementation in
// `main.go` for now. If it grows, we can move it here.

type CensusFileEntry struct {
	// The DID of the file
	Did  string
	Head string
	Rev  string
}
