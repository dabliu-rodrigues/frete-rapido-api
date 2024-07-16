package quotes

import "net/http"

func CreateQuote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello /quote!"))
}

func GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello /metrics!"))
}
