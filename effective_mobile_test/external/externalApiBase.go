package external

type IExternalAPI[T any] interface {
	Get(subdirectory, queryString string) ([]T, error)
}

func FetchData[T any](api IExternalAPI[T], subdirectory, queryString string, ch chan T, errorCh chan error) {
	res, err := api.Get(subdirectory, queryString)
	if err != nil {
		errorCh <- err
	}
	ch <- res[0]
}
