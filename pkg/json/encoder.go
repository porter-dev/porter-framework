package json

import (
	"fmt"
	"strconv"

	v "github.com/porterdev/ego/internal/value"
)

// ToJSON converts a Porter value to a JSON string
func ToJSON(v1 v.Value) (string, error) {
	if v1 == nil {
		return "null", nil
	}

	switch v1.(type) {
	case v.Boolean:
		res, _ := v1.(v.Boolean)

		return strconv.FormatBool(bool(res)), nil
	case v.Float:
		res, _ := v1.(v.Float)

		return strconv.FormatFloat(float64(res), 'e', -1, 64), nil
	case v.Integer:
		res, _ := v1.(v.Integer)

		return strconv.FormatInt(int64(res), 10), nil
	case v.String:
		res, _ := v1.(v.String)

		return "\"" + string(res) + "\"", nil
	case v.Array:
		res, _ := v1.(v.Array)

		str := "["

		for i, val := range res {
			valStr, err := ToJSON(val)

			if err != nil {
				return "", err
			}

			str += valStr

			if i+1 < len(res) {
				str += ","
			}
		}

		str += "]"

		return str, nil
	case v.Object:
		res, _ := v1.(v.Object)

		str := "{"

		count := 0

		for k, v := range res {
			keyStr, keyErr := ToJSON(k)

			if keyErr != nil {
				return "", keyErr
			}

			str += keyStr + ":"

			valStr, valErr := ToJSON(v)

			if valErr != nil {
				return "", valErr
			}

			str += valStr

			if count+1 < len(res) {
				str += ","
			}

			count++
		}

		str += "}"

		return str, nil
	}

	return "", fmt.Errorf("Value does not contain a supported Porter type")
}
