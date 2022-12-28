// Package dbms provides database connectivity.
package dbms

// SQL  http://go-database-sql.org/modifying.html
// SQLx http://jmoiron.github.io/sqlx/

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"go.opentelemetry.io/otel/trace"
)

// Config is the required properties to use the database.
type Config struct {
	User       string `json:"user,omitempty"`
	Password   string `json:"pass,omitempty"`
	Host       string `json:"host,omitempty"`
	Name       string `json:"name,omitempty"`
	DisableTLS bool   `json:"disable_tls,omitempty"`
	//PathSQL    string
}

// Store wraps `DB`
// type Store struct {
// 	*sqlx.DB
// 	PathSQL string
// }

// Open knows how to open a database connection based on the configuration.
//func Open(cfg Config) (DB, error) {
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", dbURL.String())
}

// Status returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func Status(ctx context.Context, db *sqlx.DB) error {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.dbms.status")
	defer span.End()

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	// Docker-healthcheck verion : curl -s $_PREFIX_PER_SVC/readiness | grep ok || exit 1
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// Log provides a pretty print version of the SQL query and parameters.
func Log(query string, args ...interface{}) string {
	for i, arg := range args {
		n := fmt.Sprintf("$%d", i+1)

		var a string
		switch v := arg.(type) {
		case string:
			a = fmt.Sprintf("%q", v)
		case []byte:
			a = string(v)
		case []string:
			a = strings.Join(v, ",")
		default:
			a = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, n, a, 1)
		query = strings.Replace(query, "\t", "", -1)
		query = strings.Replace(query, "\n", " ", -1)
	}

	return query
}

// ===  NULLable scheme : DEPRICATED  ===

// var (
// 	// Insert dbUNIL instead of NULL; validates as uuid, but not uuid4
// 	// @ Golang: `uuid.Nil`; @ Postgres "uuid-oosp" ext: `uuid_nil()`
// 	//dbUNIL = fmt.Sprintf("%s", uuid.Nil)
// 	dbUNIL sql.NullString = sql.NullString{String: uuid.Nil.String(), Valid: true}
// 	//dbNULL sql.NullString = sql.NullString{}
// )

// // IfEmptyToNullString ...
// // https://stackoverflow.com/questions/40266633/golang-insert-null-into-sql-instead-of-empty-string
// func IfEmptyToNullString(s string) sql.NullString {
// 	if len(s) == 0 {
// 		return sql.NullString{}
// 	}
// 	return sql.NullString{
// 		String: s,
// 		Valid:  true,
// 	}
// }

// // IfEmptyToUNIL ...
// func IfEmptyToUNIL(s string) sql.NullString {
// 	if len(s) == 0 {
// 		return sql.NullString{}
// 	}
// 	return sql.NullString{String: uuid.Nil.String(), Valid: true}
// }
