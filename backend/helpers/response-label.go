package helpers

import "net/http"

func ResponseLabel(response http.ResponseWriter, flusher http.Flusher, label string, text string) {
	textToSendToFrontend := "<" + label + ">" + text + "</" + label + ">"
	response.Write([]byte(textToSendToFrontend))
	flusher.Flush()
}

func ResponseLabelNewLine(response http.ResponseWriter, flusher http.Flusher, label string, text string) {
	textToSendToFrontend := "<" + label + ">" + text + "</" + label + "><br>"
	response.Write([]byte(textToSendToFrontend))
	flusher.Flush()
}
