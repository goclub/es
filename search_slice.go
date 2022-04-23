package es

import (
	"fmt"
	xjson "github.com/goclub/json"
	es7 "github.com/olivere/elastic/v7"
)

func SearchSlice7[T any](r *es7.SearchResult, element T) (slice []T, err error) {
	if r.Hits == nil || r.Hits.Hits == nil || len(r.Hits.Hits) == 0 {
		return
	}
	for _, hit := range r.Hits.Hits {
		item := element
		err = xjson.Unmarshal(hit.Source, &item)
		if err != nil {
			var hitJSON []byte
			var marshalErr error
			hitJSON, marshalErr = xjson.Marshal(hit)
			if marshalErr != nil {
				hitJSON = []byte(fmt.Sprintf(`%#v`, hit))
			}
			decodeFailMessage := err.Error() + ";" + fmt.Sprintf("hit.Index(%s) hit.Id(%s)", hit.Index, hit.Id)

			err = newDecodeSearchResultError(err, hitJSON, element, decodeFailMessage)
			return
		}
		slice = append(slice, item)
	}
	return
}
