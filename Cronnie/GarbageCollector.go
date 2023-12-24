package Cronnie

import (
	"fmt"
	"strings"
	"time"
)

func (ci *Instance) garbageCollector() {
	for {

		expired := time.Now().Add(-ci.keepCompleted)
		q := `SELECT id FROM jobs where completed_at < $1`

		r, e := ci.conn.Query(ci.ctx, q, expired.UTC())
		if e != nil {
			ci.logger.Printf("Error while running garbage collection. Error :: %s\n", e)
			continue
		}

		var oldItems []int
		for r.Next() {
			var id int
			e := r.Scan(&id)
			if e != nil {
				fmt.Printf("Cant unmarshal :: %#v\n", e)
			}
			oldItems = append(oldItems, id)
		}

		if len(oldItems) == 0 {
			time.Sleep(time.Second * 30)
			continue
		}

		// convert old items to csv
		stringValues := make([]string, len(oldItems))
		for i, v := range oldItems {
			stringValues[i] = fmt.Sprint(v)
		}
		csv := strings.Join(stringValues, ",")

		q = `DELETE FROM jobs WHERE id IN(` + csv + `);`
		_, e = ci.conn.Query(ci.ctx, q)
		if e != nil {
			ci.logger.Printf("Error while dropping items %s", e)
			continue
		}

		time.Sleep(time.Second * 30)
	}
}
