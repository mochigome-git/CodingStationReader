// pkg/thd/text_picker.go
package thd

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// GetLog searches for a given string in a log file. If the string is found, it returns the portion of the log file containing the string.
// If the string is not found in the specified log file, it will attempt to find and search in additional log files matching the given pattern.
func GetLog(fileName, findStr, filePath string) (string, error) {
	const perRead int64 = 512

	filePath = filePath + fileName

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var contents []byte
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting file stat: %w", err)
	}

	findBytes := []byte(findStr)
	findLength := len(findBytes)
	if int64(findLength) > perRead {
		return "", fmt.Errorf("findStr length is larger than a read")
	}

	lastRead := stat.Size()

	for lastRead > 0 {
		readSize := perRead
		if lastRead < readSize {
			readSize = lastRead
		}

		readBytes := make([]byte, readSize)
		_, err = file.ReadAt(readBytes, lastRead-readSize)
		if err != nil {
			return "", fmt.Errorf("error reading file: %w", err)
		}

		contents = append(readBytes, contents...)
		indexOf := bytes.LastIndex(contents, findBytes)

		if indexOf != -1 {
			log.Printf("text_picker.go: %s", string(contents[indexOf+findLength:]))
			return string(contents[indexOf+findLength:]), nil
		}

		lastRead -= readSize
	}

	return "", nil
}

//func GetLog(fileName string, findStr string, filePath string) (string, error) {
//	const perRead int64 = 512
//
//	file, err := os.Open(filePath + fileName)
//	if err != nil {
//		// Error opening file
//		return "", fmt.Errorf("error opening file: %w", err)
//	}
//	defer file.Close()
//
//	stat, err := file.Stat()
//	if err != nil {
//		// Error getting file stat
//		return "", fmt.Errorf("error getting file stat: %w", err)
//	}
//
//	// Convert findStr to bytes for fast searching.
//	findBytes := []byte(findStr)
//	findLength := len(findBytes)
//	// The length of findStr can't be larger than a read.
//	if int64(findLength) > perRead {
//		return "", fmt.Errorf("findStr length is larger than a read")
//	}
//
//	var lastRead = stat.Size()
//	var contents = make([][]byte, lastRead/perRead+1)
//	var lastIndex = len(contents) - 1
//	var saveIndex = lastIndex
//
//	for {
//		var readBytes []byte
//
//		if lastRead == 0 {
//			break
//		}
//		if lastRead-perRead > -1 {
//			readBytes = make([]byte, perRead)
//			lastRead = lastRead - perRead
//		} else {
//			readBytes = make([]byte, lastRead-0)
//			lastRead = 0
//		}
//
//		_, err = file.ReadAt(readBytes, lastRead)
//		if err != nil {
//			// Error reading file
//			return "", fmt.Errorf("error reading file: %w", err)
//		}
//
//		var indexOf = bytes.Index(readBytes, findBytes)
//
//		if indexOf != -1 {
//			contents[saveIndex] = readBytes[indexOf+findLength:]
//			saveIndex -= 1
//			break
//		} else {
//			if saveIndex < lastIndex {
//				// Take a small chunk of the beginning of the last found (equal to findStr's length)
//				// add to a small ending chunk of this found (equal to findStr's length)
//				// However, if this found is less than findStr length, grab whatever is available.
//				var halfpart []byte
//				if len(readBytes) < findLength {
//					halfpart = append(readBytes, contents[saveIndex+1][:findLength]...)
//				} else {
//					halfpart = append(readBytes[len(readBytes)-findLength:], contents[saveIndex+1][:findLength]...)
//				}
//
//				var indexOf2 = bytes.Index(halfpart, findBytes)
//				if indexOf2 != -1 {
//					saveIndex = saveIndex + 1
//					contents[saveIndex] = append(halfpart[indexOf2+findLength:], contents[saveIndex][findLength:]...)
//					saveIndex -= 1
//					break
//				}
//			}
//			contents[saveIndex] = readBytes
//			saveIndex -= 1
//		}
//	}
//
//	for i := saveIndex; i > -1; i-- {
//		contents[saveIndex] = []byte{}
//	}
//
//	log.Printf("text_picker.go: %s", string(bytes.Join(contents, []byte{})))
//
//	return string(bytes.Join(contents, []byte{})), nil
//}
