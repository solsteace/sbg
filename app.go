package main

const (
	fLAG_HELP         = "--help"
	fLAG_HELP_SHORT   = "-h"
	fLAG_CONFIG       = "--config"
	fLAG_CONFIG_SHORT = "-C"
)

const (
	sTATE_READY       = "ready"
	sTATE_NEED_PARAMS = "need_params"
	sTATE_INSTA_END   = "insta_end"
)

const (
	tEXT_HELP_SHORT = `sbg - braille-art-to-SVG converter

sbg is a tool to turn a braille art into SVG that you could embed to HTML files.

Quick usage:
	cat <your-braille-art> | sbg > result.html

For complete details and option list, use sbg --help
`

	tEXT_HELP = `sbg - braille-art-to-SVG converter

sbg is a tool to turn a braille art into SVG that you could embed to HTML files.

Usage:
	sbg --source [path] --destination [path]

Flags:
	-C [path]
	--config [path]
	
	Sets the path of the config file. If [path] omitted, "./config.json"
	would be used instead.
	
	If succesfully used, the value for the emitted configuration flags would be 
	overidden by the file content.
`
)
