package until

import (
	"fmt"
    "github.com/google/uuid"
)

func UidString() string {
	v1, _ := uuid.NewUUID()
	return Md5Sum(fmt.Sprintf("%s", v1))
}
