package Plugin

/*
Type to define paths and directories
Used to specify specific files such as HitIcon.json in oxide data dir
or specify directory such as Arkan in oxide data dir to remove contents of whole directory
*/
type DataPath struct {
	Directory bool
	Path      string
}
