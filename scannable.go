package preql

import (
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ScannableType holds information about a type that is scannable.
type ScannableType struct {
	// Name is the name of the type.
	Name string

	// Fields are the types field names, indexed by sql name (as defined in the sql struct tag).
	Fields map[string]string
}

// parseScannableTypes parses all scannable types and returns them, omitting types that
// are not scannable.
func parseScannableTypes(pkg *packages.Package) []*ScannableType {
	byName := make(map[string]*ScannableType)
	for _, typ := range pkg.TypesInfo.Types {
		if t := parseScannableType(typ); t != nil {
			byName[t.Name] = t
		}
	}

	resp := make([]*ScannableType, 0, len(byName))
	for _, t := range byName {
		resp = append(resp, t)
	}

	return resp
}

// parseScannableType parses a type into a ScannableType definition.
// Returns nil if the type is not Scannable.
func parseScannableType(typ types.TypeAndValue) *ScannableType {
	if !typ.IsType() {
		return nil
	}

	named, ok := typ.Type.(*types.Named)
	if !ok {
		return nil
	}

	flds := fields(typ.Type)
	if len(flds) == 0 {
		return nil
	}

	return &ScannableType{
		Name:   named.Obj().Name(),
		Fields: flds,
	}
}

// parseStructFields parses a struct type for its fields.
// It returns a map of field names indexed by sql tag.
func parseStructFields(strct *types.Struct) map[string]string {
	flds := make(map[string]string)
	for i := 0; i < strct.NumFields(); i++ {
		field := strct.Field(i)
		tag := sqlTag(strct.Tag(i))
		name := field.Name()

		if tag == "" {
			continue
		}

		flds[tag] = name
	}

	return flds
}

// sqlTag parses a struct tag to get the sql tag value.
// Returns the empty string if no tag was matched.
func sqlTag(tag string) string {
	return reflect.StructTag(tag).Get("sql")
}

// Receiver returns the receiver name of the scannable type.
func (t ScannableType) Receiver() string {
	return strings.ToLower(t.Name)[0:1]
}
