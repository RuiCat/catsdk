package fontscan

import (
	"fmt"
	"os"

	"sdk/typesetting/font"
	"sdk/typesetting/opentype/api"
	meta "sdk/typesetting/opentype/api/metadata"
	"sdk/typesetting/opentype/loader"
	"sdk/typesetting/opentype/tables"
)

// Location identifies where a font.Face is stored.
type Location = api.FontID

// footprint is a condensed summary of the main information
// about a font, serving as a lightweight surrogate
// for the original font file.
type footprint struct {
	// Location stores the adress of the font resource.
	Location Location

	// Family is the general nature of the font, like
	// "Arial"
	// Note that, for performance reason, we store the
	// normalized version of the family name.
	Family string

	// Runes is the set of runes supported by the font.
	Runes runeSet

	// set of scripts deduced from Runes
	scripts scriptSet

	// set of languages deduced from Runes
	langs langSet

	// Aspect precises the visual characteristics
	// of the font among a family, like "Bold Italic"
	Aspect meta.Aspect

	// isUserProvided is set to true for fonts add manually to
	// a FontMap
	// User fonts will always be tried if no other fonts match,
	// and will have priority among font with same family name.
	//
	// This field is not serialized in the index, since it is always false
	// for system fonts.
	isUserProvided bool
}

func newFootprintFromFont(f font.Font, location Location, md meta.Description) (out footprint) {
	out.Runes, out.scripts, _ = newCoveragesFromCmap(f.Cmap, nil)
	out.langs = newLangsetFromCoverage(out.Runes)
	out.Family = meta.NormalizeFamily(md.Family)
	out.Aspect = md.Aspect
	out.Location = location
	out.isUserProvided = true
	return out
}

func newFootprintFromLoader(ld *loader.Loader, isUserProvided bool, buffer scanBuffer) (out footprint, _ scanBuffer, err error) {
	raw := buffer.tableBuffer

	// since raw is shared, special car must be taken in the parsing order

	raw, _ = ld.RawTableTo(loader.MustNewTag("OS/2"), raw)
	fp := tables.FPNone
	if os2, _, err := tables.ParseOs2(raw); err != nil {
		fp = os2.FontPage()
	}

	// we can use the buffer since ProcessCmap do not keep any reference on
	// the input slice
	raw, err = ld.RawTableTo(loader.MustNewTag("cmap"), raw)
	if err != nil {
		return footprint{}, buffer, err
	}
	tb, _, err := tables.ParseCmap(raw)
	if err != nil {
		return footprint{}, buffer, err
	}
	cmap, _, err := api.ProcessCmap(tb, fp)
	if err != nil {
		return footprint{}, buffer, err
	}

	out.Runes, out.scripts, buffer.cmapBuffer = newCoveragesFromCmap(cmap, buffer.cmapBuffer) // ... and build the corresponding rune set

	out.langs = newLangsetFromCoverage(out.Runes)

	family, aspect, raw := meta.Describe(ld, raw)
	out.Family = meta.NormalizeFamily(family)
	out.Aspect = aspect
	out.isUserProvided = isUserProvided

	buffer.tableBuffer = raw

	return out, buffer, nil
}

// loadFromDisk assume the footprint location refers to the file system
func (fp *footprint) loadFromDisk() (font.Face, error) {
	location := fp.Location

	file, err := os.Open(location.File)
	if err != nil {
		return nil, err
	}

	faces, err := font.ParseTTC(file)
	if err != nil {
		return nil, err
	}

	if index := int(location.Index); len(faces) <= index {
		// this should only happen if the font file as changed
		// since the last scan (very unlikely)
		return nil, fmt.Errorf("invalid font index in collection: %d >= %d", index, len(faces))
	}

	return faces[location.Index], nil
}
