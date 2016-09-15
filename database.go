package main

func fetchAll(query string) (results []map[string]interface{}) {
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
		res := map[string]interface{}{}
		for idx, col := range cols {
			res[col] = vars[idx]
		}
		results = append(results, res)
	}
	rows.Close()
	return results
}
