package utils

import (
	"fmt"
)

func Jump(pageIndex, pageSize int, queryString string) string {
	if queryString == "" {
		return fmt.Sprintf("pageSize=%d&pageIndex=%d", pageSize, pageIndex)
	}
	return fmt.Sprintf("%s&pageSize=%d&pageIndex=%d", queryString, pageSize, pageIndex)
}

func Prefix(queryString string, currentPage, pageSize int) string {
	if currentPage > 1 {
		return Jump(currentPage - 1, pageSize, queryString)
	}
	return Jump(1, pageSize, queryString)
}

func Suffix(queryString string, currentPage, pageSize, totalPage int) string {
	if currentPage < totalPage {
		return Jump(currentPage + 1, pageSize, queryString)
	}
	return Jump(totalPage, pageSize, queryString)
}

func ShowPrefix(currentPage int) bool {
	return currentPage > 1
}

func ShowSuffix(currentPage, totalPage int) bool {
	return currentPage < totalPage
}

func ShowFirst(currentPage int) bool {
	return currentPage > 3
}

func ShowLast(currentPage, totalPage int) bool {
	return currentPage < totalPage - 2
}

func GetPageNumber(currentPage, totalPage int) []int {
	pageNumber := make([]int, 5)
	if totalPage <= 5 {
		pageNumber = make([]int, totalPage)
		for index := range pageNumber {
			pageNumber[index] = index + 1
		}
	} else {
		switch {
		case currentPage <= 3:
			pageNumber[0] = 1
			pageNumber[1] = 2
			pageNumber[2] = 3
			pageNumber[3] = 4
			pageNumber[4] = 5
		case currentPage >= totalPage-2:
			pageNumber[0] = totalPage - 4
			pageNumber[1] = totalPage - 3
			pageNumber[2] = totalPage - 2
			pageNumber[3] = totalPage - 1
			pageNumber[4] = totalPage
		default:
			// 前后各取2个
			pageNumber[0] = currentPage - 2
			pageNumber[1] = currentPage - 1
			pageNumber[2] = currentPage
			pageNumber[3] = currentPage + 1
			pageNumber[4] = currentPage + 2
		}
	}
	return pageNumber
}

func GetBeginIndex(currentPage, pageSize int) int {
	if currentPage > 1 {
		return pageSize * (currentPage - 1)
	}
	return 0
}

func GetEndIndex(pageSize, currentPage, totalRecord int) int {
	endIndex := pageSize * currentPage
	if endIndex > totalRecord {
		endIndex = totalRecord
	}
	return endIndex
}

func GetTotalPage(pageSize, totalRecord int) int {
	if pageSize > 0 {
		return (totalRecord - 1 ) / pageSize + 1
	}
	return 0
}

