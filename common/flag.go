package common

import "flag"

func ParseFlag() {
	debugFlag := flag.Bool("d", false, "Enable debug mode")
	ARCHIVEFILELOACTIONFlag := flag.String("p", "./static/archive/", "Assign HTML file path")
	MEILIHOSTFlag := flag.String("h", "http://127.0.0.1:7700", "Assign MeiliSearch host")
	flag.Parse()
	DEBUG = *debugFlag
	ARCHIVEFILELOACTION = *ARCHIVEFILELOACTIONFlag
	MEILIHOST = *MEILIHOSTFlag
}
