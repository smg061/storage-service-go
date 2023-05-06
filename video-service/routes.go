package main

import "net/http"

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<body style="">
			<div style="display:flex; justify-content:center; align-items:center">
				<h1 style="color: blue">Ewiwa my beloved</h1>
			</div>
		</body>
		`))
	})
	mux.HandleFunc("/videos", ShowVideo)
	return mux
}