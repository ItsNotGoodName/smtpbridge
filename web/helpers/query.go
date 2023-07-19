package helpers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
)

func Query(queries map[string]string, vals ...any) string {
	tuple := []lo.Tuple2[string, any]{}
	for i := 0; i < len(vals); i += 2 {
		tuple = append(tuple, lo.Tuple2[string, any]{A: (vals[i]).(string), B: vals[i+1]})
	}

	var newQueries []string
	for _, t := range tuple {
		newQueries = append(newQueries, t.A+"="+fmt.Sprint(t.B))
	}

	for k, v := range queries {
		_, found := queryTupleFind(tuple, k)
		if !found {
			newQueries = append(newQueries, k+"="+v)
		}
	}

	sort.Slice(newQueries, func(i, j int) bool {
		return newQueries[i] < newQueries[j]
	})

	return strings.Join(newQueries, "&")
}

func queryTupleFind(tuple []lo.Tuple2[string, any], key string) (int, bool) {
	for i, t := range tuple {
		if t.A == key {
			return i, true
		}
	}

	return 0, false
}
