package preql

import (
	"fmt"
	"go/types"
)

// fields returns the fields of a type, indexed by their sql tag.
// Returns nil if the type does not have fields.
func fields(t types.Type) map[string]string {
	if t == nil {
		return nil
	}

	switch v := t.(type) {
	case *types.Named:
		return fields(v.Underlying())
	case *types.Struct:
		n := v.NumFields()
		res := make(map[string]string)
		for i := 0; i < n; i++ {
			name := v.Field(i).Name()
			if name == "" {
				continue
			}

			tag := sqlTag(v.Tag(i))
			if tag == "" {
				continue
			}

			res[tag] = name
		}
		return res
	case *types.Pointer:
		return fields(v.Elem())
	case *types.Basic:
		return nil
	case *types.Interface:
		return nil
	default:
		fmt.Printf("Unexpected type %T\n", t)
		return nil
	}
}
