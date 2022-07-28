package DLL

/**
DLL a structure to define dlls to be added to the server.
This is to add anything such as primalprofiler or discord extension or anything of that kind
*/
type DLL struct {
	Name string
	//Path is relative to the dll's directory
	FileName string
	Binary   []byte
}
