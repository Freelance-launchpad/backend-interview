package pkg

import "github.com/Freelance-launchpad/backend-interview/common/jdb"

type ListFileInfosFilters struct {
	Tags []string `form:"tag"` // Used as ?tag=foo&tag=bar
	jdb.QueryPagination
}

type ResponseListFileInfos struct {
	Count int64      `json:"count"`
	Files []FileInfo `json:"files"`
}
