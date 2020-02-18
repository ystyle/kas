package file

import "fmt"

func FormatBytesLength(length int) string {
	if length < 1024*1024 {
		return fmt.Sprintf("%d K", length/(1024))
	} else if length < 1024*1024*1024 {
		return fmt.Sprintf("%d M", length/(1024*1024))
	} else if length < 1024*1024*1024*1024 {
		return fmt.Sprintf("%d G", length/(1024*1024*1024))
	} else {
		return fmt.Sprintf("%d T", length/(1024*1024*1024*1024))
	}
}
