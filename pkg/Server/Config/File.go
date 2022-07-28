package Config

/*
File is an arbitrary config file meant to configure anything not related
*/
type File struct {
	Name       string
	Contents   string
	TargetPath string
}

/*
Parse uses the file with configuration context to do substitution for any placeholders used
*/
func (f File) Parse(cfg SubstitutionsGetter) {

	//Substitute contents, todo: replace by placeholder? {KEY} ?

	//Parse contents

}

type SubstitutionsGetter interface {
	List() map[string]string
}
