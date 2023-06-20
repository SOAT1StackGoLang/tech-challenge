package setup

type (
	ServiceParams struct {
		PostgresConnParamsAcquirer func() sqlpersis
	}
)
