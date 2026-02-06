package cursor

import "encoding/base64"
import "strconv"

// EncodeCursor Encodes the cursor id on the client URL
func EncodeCursor(id int64) string {
	idBytes := []byte(strconv.FormatInt(id, 10))
	return base64.StdEncoding.EncodeToString(idBytes)
}

// DecodeCursor Dencodes the cursor id on the client URL
func DecodeCursor(cursor string) (int64, error) {
	b64, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(b64), 10, 64)
}
