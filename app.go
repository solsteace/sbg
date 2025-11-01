package main

const (
	fLAG_SOURCE            = "--source"
	fLAG_SOURCE_SHORT      = "-S"
	fLAG_DESTINATION       = "--destination"
	fLAG_DESTINATION_SHORT = "-D"
	fLAG_VARIATION         = "--variation"
	fLAG_VARIATION_SHORT   = "-R"
	fLAG_HELP              = "--help"
	fLAG_HELP_SHORT        = "-h"
)

const (
	tEXT_HELP_SHORT = `sbg - braille-art-to-SVG converter

sbg is a tool to turn a braille art into SVG that you could embed to HTML files.

Quick usage:
	sbg --source [path] --destination [path]

For complete details and option list, use sbg --help
`

	tEXT_HELP = `sbg - braille-art-to-SVG converter

sbg is a tool to turn a braille art into SVG that you could embed to HTML files.

Usage:
	sbg --source [path] --destination [path]

Flags:
	-S [path]
	--source [path]

	The input of braille art that would be converted. [file] is an optional 
	filepath to the file containing the braille art, typically stored as .txt.
	If [file] omitted, stdin would be used instead
	
	Note that the input MUST NOT contain any characters other than braille 
	characters (U+2800 - U+28FF).


	-D [path]
	--destination [path]

	The path for the result file. The resulting SVG would be saved in a HTML file. 
	If [path] omitted, stdout would be used instead.


	-R [name]
	--variation [name]

	This flag refers to the resulting pattern variation visible in the SVG. If 
	[name] omitted, line-horizontal would be used. [name] possible values are:
		line-horizontal
		line-vertical
		diagonal-up
		diagonal-down
`
)
