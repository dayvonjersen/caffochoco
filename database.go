package main

import "fmt"

func fetchAll(query string) (results []map[string]string) {
	rows, err := db.Query(query)
	checkErr(err)
	cols, err := rows.Columns()
	checkErr(err)
	for rows.Next() {
		vars := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var v string
			vars[i] = &v
		}
		checkErr(rows.Scan(vars...))
		res := map[string]string{}
		for idx, col := range cols {
			res[col] = fmt.Sprintf("%s", *vars[idx].(*string))
		}
		results = append(results, res)
	}
	rows.Close()
	return results
}
