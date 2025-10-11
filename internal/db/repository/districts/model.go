package districts

type District struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
	Area []byte `db:"area"`
}
