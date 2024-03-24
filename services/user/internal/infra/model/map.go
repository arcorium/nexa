package model

// DataAccessModelMapOption function pointer to be used for mapping configuration.
// D is domain model type and R is repository or database model type
type DataAccessModelMapOption[D, R any] func(D, R)
