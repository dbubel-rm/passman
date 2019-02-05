// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ErrInvalidDBProvided is returned in the event that an uninitialized db is
// used to perform actions against.
// var ErrInvalidDBProvided = errors.New("invalid DB provided")

// DB is a collection of support for different DB technologies. Currently
// only MongoDB has been implemented. We want to be able to access the raw
// database support for the given DB so an interface does not work. Each
// database is too different.
type DB struct {
	Database *sqlx.DB
}

// New returns a new DB value for use with MongoDB based on a registered
// master session.
func New(url string) (*DB, error) {
	// db, err := sqlx.Connect("mysql", "passman:wtfthispasswordNeedsLonger29@tcp(production-database.cfneifgjtyib.us-east-1.rds.amazonaws.com:3306)/passman")
	mysql, err := sqlx.Connect("mysql", url)

	if err != nil {
		return nil, errors.Wrap(err, "Error connecting to DB when application started")
	}
	// mysql.SetMaxIdleConns(10)
	// mysql.DB.SetMaxIdleConns(10)
	// mysql.DB.SetMaxOpenConns(10)

	db := DB{
		Database: mysql,
	}

	return &db, nil
}

// Close closes a DB value being used with MongoDB.
func (db *DB) Close() {
	db.Database.Close()
}

// // Copy returns a new DB value for use with MongoDB based on master session.
// func (db *DB) Copy() *DB {
// 	ses := db.session.Copy()

// 	// As per the mgo documentation, https://godoc.org/gopkg.in/mgo.v2#Session.DB
// 	// if no database name is specified, then use the default one, or the one that
// 	// the connection was dialed with.
// 	newDB := DB{
// 		database: ses.DB(""),
// 		session:  ses,
// 	}

// 	return &newDB
// }

// // Execute is used to execute MongoDB commands.
// func (db *DB) Execute(ctx context.Context, collName string, f func(*mgo.Collection) error) error {
// 	ctx, span := trace.StartSpan(ctx, "platform.DB.Execute")
// 	defer span.End()

// 	if db == nil || db.session == nil {
// 		return errors.Wrap(ErrInvalidDBProvided, "db == nil || db.session == nil")
// 	}

// 	return f(db.database.C(collName))
// }

// // ExecuteTimeout is used to execute MongoDB commands with a timeout.
// func (db *DB) ExecuteTimeout(ctx context.Context, timeout time.Duration, collName string, f func(*mgo.Collection) error) error {
// 	ctx, span := trace.StartSpan(ctx, "platform.DB.ExecuteTimeout")
// 	defer span.End()

// 	if db == nil || db.session == nil {
// 		return errors.Wrap(ErrInvalidDBProvided, "db == nil || db.session == nil")
// 	}

// 	db.session.SetSocketTimeout(timeout)

// 	return f(db.database.C(collName))
// }

// // StatusCheck validates the DB status good.
// func (db *DB) StatusCheck(ctx context.Context) error {
// 	ctx, span := trace.StartSpan(ctx, "platform.DB.StatusCheck")
// 	defer span.End()

// 	return nil
// }

// // Query provides a string version of the value
// func Query(value interface{}) string {
// 	json, err := json.Marshal(value)
// 	if err != nil {
// 		return ""
// 	}

// 	return string(json)
// }
