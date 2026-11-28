package helpers

import "net/http"

type Recorder struct {
	http.ResponseWriter
	Status      int
	Bytes       int64
	wroteHeader bool
}

func (recorder *Recorder) WriteHeader(Status int) {
	if recorder.wroteHeader {
		return
	}

	recorder.Status = Status
	recorder.wroteHeader = true
	recorder.ResponseWriter.WriteHeader(Status)
}

func (recorder *Recorder) Write(payload []byte) (int, error) {
	if !recorder.wroteHeader {
		recorder.WriteHeader(http.StatusOK)
	}

	written, err := recorder.ResponseWriter.Write(payload)
	recorder.Bytes += int64(written)
	return written, err
}

func (recorder *Recorder) Unwrap() http.ResponseWriter {
	return recorder.ResponseWriter
}

