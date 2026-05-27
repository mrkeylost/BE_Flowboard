package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (uuidArray *UUIDArray) Scan(value interface{}) error {

	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("Failed to parse uuid array: unsupported data type")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	elements := strings.Split(str, ",")

	*uuidArray = make(UUIDArray, 0, len(elements))

	for _, item := range elements {
		item = strings.TrimSpace(strings.Trim(item, `"`))
		if item == "" {
			continue
		}

		id, err := uuid.Parse(item)
		if err != nil {
			return fmt.Errorf("Invalid UUID Format: %v", err)
		}

		*uuidArray = append(*uuidArray, id)
	}

	return nil
}

func (uuidArray UUIDArray) Value() (driver.Value, error) {

	if len(uuidArray) == 0 {
		return "{}", nil
	}

	postgreFormat := make([]string, 0, len(uuidArray))

	for _, value := range uuidArray {
		postgreFormat = append(postgreFormat, fmt.Sprintf(`"%s"`, value.String()))
	}

	return "{" + strings.Join(postgreFormat, ",") + "}", nil
}

func (UUIDArray) GORMDataType() string {
	return "uuid[]"
}
