package repository

import (
	"bufio"
	"deus.ai-code-challenge/domain"
	"fmt"
	"os"
)

type pageURL = string

type pageInMemoryRepository struct {
	data map[pageURL]struct{}
}

func newPageInMemoryRepository(pages []domain.PageURL) domain.PageRepository {
	data := make(map[pageURL]struct{}, len(pages))
	for _, page := range pages {
		data[pageURL(page)] = struct{}{}
	}

	return &pageInMemoryRepository{
		data: data,
	}
}

func (i *pageInMemoryRepository) Exists(url domain.PageURL) (bool, error) {
	_, found := i.data[pageURL(url)]

	return found, nil
}

// ReadPages simply reads a file line by line and returns each line as a page url
func ReadPages(filePath string) ([]domain.PageURL, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer func(readFile *os.File) {
		_ = readFile.Close()
	}(readFile)

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	ret := make([]domain.PageURL, 0)
	for fileScanner.Scan() {
		t := fileScanner.Text()
		if t != "" {
			ret = append(ret, domain.PageURL(t))
		}
	}

	return ret, nil
}
