package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"iot.ericslien.com/config"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/**
* Our database connection instance
 */
var db_connection *gorm.DB

var next_ping time.Time

var db_conf config.SQLDatabase

var db_initialized bool

//Makes sure the db connection stays open, by pinging and reconnecting if no ping every 5 seconds
//We will keep one database connection open in the app, never closing it.
//however!  If we get disconnected for any reason, we will only be disconnected for a max of 5 seconds unless the
//db is actually down or unavailable to us
//Below is a link (one of many) dicussing how database/sql best practice is either to terminate a connection after ever request
//or keep one open ALWAYS in the global context, and NEVER explicitly close it.  Given workers need to share this connection,
//we just cannot close the connection at the end of every request, and it's considered bad practice to Open() more than one
//instance in gorm.  That leaves us with only the option of keeping a persistent connection going.  Obviously, this is counter to how
//mongo worked, but it's a well-documented difference.  I'm sticking with what others have found.
//https://github.com/go-sql-driver/mysql/issues/461
func keepalive() {
	t := time.Now()
	if t.After(next_ping) {
		logrus.Debug("Pinging DB")
		if db_connection == nil || db_connection.DB().Ping() != nil {
			logrus.Debug("Reconnecting to database")
			db, err := New()
			if err != nil {
				e := fmt.Sprintf("Database connection was closed and could not be brought back up: %s", err)
				logrus.Fatal(e)
				panic(e)
			}
			db_connection = db
		}
		//rate limit pings, so that we never do more once every 10 seconds
		next_ping = t.Add(time.Second * 10)
	}
}

/**
* Provides a standard way of getting our db in an app
 */
func new(conf config.SQLDatabase) (*gorm.DB, error) {
	logrus.Debug("Opening a new db connection")
	var db *gorm.DB
	var err error

	switch conf.Driver {
	case "mysql":
		//mysql connection string
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		//user:password@tcp(localhost:5555)/dbname
		connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
		db, err = gorm.Open(conf.Driver, connString)
		// db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(conf.MaxConnections)
	}
	db_initialized = true
	return db, err
}

/**
* Sets the globally-used db configuration for the running instance/app
 */
func SetConfig(conf config.SQLDatabase) {
	db_conf = conf
	db_initialized = false
	logrus.Debug(fmt.Sprintf("Database configuration set to %+v", conf))
}

/**
* Instantiates a new database connection that this module does not try to keep track of
* You can run out of connections using this function
 */
func New() (*gorm.DB, error) {
	return new(db_conf)
}

//Open Get or open a connection to the db
func Open() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if db_connection == nil || !db_initialized {
		if db_conf.User == "" {
			e := "Database configuration not initialized"
			logrus.Fatal(e)
			panic(e)
		} else {
			db, err = New()
			if err == nil {
				db_connection = db
				return db_connection, nil
			}
		}
	}
	keepalive()
	return db_connection, err
}

//this interface just ensures we aught to be able to use both the
// JSONDBField* functions in database calls.  Without both,
// we won't be able to insert and get data from a field
// see the database/sql scanner and valuer interfaces
// I'm combining them for "integrity"
// If this turns out to be an anti-pattern (cause it feels strange)
// I'm sorry, but it does actually work in tests
type JSONFielder interface {
	Value() (driver.Value, error)
	Scan(src interface{}) error
}

func JSONFielderValue(j JSONFielder) (driver.Value, error) {
	return json.Marshal(j)
}

func JSONFielderScan(j JSONFielder, src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}
	return json.Unmarshal(source, &j)
}
