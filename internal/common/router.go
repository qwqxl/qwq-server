package common

type HTTPRouter struct {
	API []string
}

type Route struct {
	PathGroup string
	PathName  []string
}
