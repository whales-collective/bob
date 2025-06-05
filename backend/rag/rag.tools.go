package rag

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetChunksOfCloneDocuments(cloneName string) ([]string, error) {
	contents, err := GetContentFiles("/app/docs/"+cloneName, ".md")
	if err != nil {
		return nil, fmt.Errorf("error getting content files for %s agent: %w", cloneName, err)
	}
	for _, content := range contents {
		fmt.Println("📄", cloneName, "content file:", content)
	}
	chunks := []string{}
	for _, content := range contents {
		chunks = append(chunks, ChunkText(content, 512, 210)...)
	}
	return chunks, nil
}

// ChunkText takes a text string and divides it into chunks of a specified size with a given overlap.
// It returns a slice of strings, where each string represents a chunk of the original text.
//
// Parameters:
//   - text: The input text to be chunked.
//   - chunkSize: The size of each chunk.
//   - overlap: The amount of overlap between consecutive chunks.
//
// Returns:
//   - []string: A slice of strings representing the chunks of the original text.
func ChunkText(text string, chunkSize, overlap int) []string {
	chunks := []string{}
	for start := 0; start < len(text); start += chunkSize - overlap {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[start:end])
	}
	return chunks
}

// GetContentFiles searches for files with a specific extension in the given directory and its subdirectories.
//
// Parameters:
// - dirPath: The directory path to start the search from.
// - ext: The file extension to search for.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func GetContentFiles(dirPath string, ext string) ([]string, error) {
	content := []string{}
	_, err := ForEachFile(dirPath, ext, func(path string) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content = append(content, string(data))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ForEachFile iterates over all files with a specific extension in a directory and its subdirectories.
//
// Parameters:
// - dirPath: The root directory to start the search from.
// - ext: The file extension to search for.
// - callback: A function to be called for each file found.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func ForEachFile(dirPath string, ext string, callback func(string) error) ([]string, error) {
	var textFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			textFiles = append(textFiles, path)
			err = callback(path)
			// generate an error to stop the walk
			if err != nil {
				return err
			}
		}
		return nil
	})
	return textFiles, err
}
