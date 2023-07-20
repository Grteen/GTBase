package utils

import "GtBase/src/page"

func IsBucketFilePath(filePath string) bool {
	return filePath == page.BucketPageFilePathToDo
}
