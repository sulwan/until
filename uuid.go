package until

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func UidString() string {
	return Md5Sum(fmt.Sprintf("%s", uuid.Must(uuid.NewV4())))
}
