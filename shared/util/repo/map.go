package repo

// DataAccessModelMapOption function pointer to be used for mapping configuration.
// D is domain model type and R is repository or database model type
type DataAccessModelMapOption[D, R any] func(D, R)

func MapOptionFunc[D, M any](d *D, m *M) func(option *DataAccessModelMapOption[*D, *M]) {
	return func(option *DataAccessModelMapOption[*D, *M]) {
		(*option)(d, m)
	}
}
