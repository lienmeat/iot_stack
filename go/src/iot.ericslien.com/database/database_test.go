package database

import (
	// "github.com/google/uuid"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"iot.ericslien.com/config"

	"iot.ericslien.com/logging"

	_ "github.com/jinzhu/gorm"
	// "time"
	"github.com/jinzhu/gorm"
)

var count int = 0

type Something struct {
	Data string `json:"data"`
}

func (s *Something) Value() (driver.Value, error) {
	return JSONFielderValue(s)
}

func (s *Something) Scan(src interface{}) error {
	return JSONFielderScan(s, src)
}

type DBTESTS struct {
	ID        uint64
	Name      string
	Number    uint64
	Thing     string
	Thing1    string
	Thing2    string
	Thing3    string
	Thing4    string
	Thing5    string
	Thing6    string
	Thing7    string
	Thing8    string
	Thing9    string
	Thing10   string
	Thing11   string
	Thing12   string
	Thing13   string
	Thing14   string
	Thing15   string
	Thing16   string
	Thing17   string
	Thing18   string
	Thing19   string
	Thing20   string
	Something *Something `gorm:"type:json"`
}

var db_test_conf config.SQLDatabase

func init() {
	db_test_conf = config.SQLDatabase{User: "root", Host: "localhost", Password: "123", Driver: "mysql", Port: "3306", Database: "test", MaxConnections: 100}

	//this module is not part of an app instance, so you must set up a database to test against that is separate (shouldn't include app's config here)
	SetConfig(db_test_conf)
	logging.SetupLogging("test", "debug")
}

func TestMigrate(t *testing.T) {
	db, _ := Open()

	if db.AutoMigrate(&DBTESTS{}).Error != nil {
		t.Error("Migration failed")
	}

	if !db.HasTable(&DBTESTS{}) {
		t.Error("Table doesn't exist after migration")
	}

	dbt := DBTESTS{Name: "MYNAME", Number: 1000}
	fmt.Printf("record: %v", dbt)
	if db.Create(&dbt).Error != nil {
		t.Error("Could not create a new record on migrated table")
	}
	db.Delete(&dbt, &dbt)
}

func TestOpen(t *testing.T) {
	_, err := Open()
	if err != nil {
		t.Error("Could not Get DB")
	}
}

func TestGetMany(t *testing.T) {
	db, err := Open()

	if err != nil {
		t.Error("Could not Get DB initially")
	}
	for i := 1; i < 10000; i++ {
		db, err = Open()
		if err != nil {
			t.Error(fmt.Sprintf("Could not get DB %d times", i))
			return
		}
		e := db.Exec("SELECT * FROM dbtests").Error
		if e != nil {
			t.Error(fmt.Sprintf("%s @ %d", e, i))
			return
		}
	}
}

func TestJsonField(t *testing.T) {
	db, _ := Open()

	err := db.AutoMigrate(&DBTESTS{}).Error
	if err != nil {
		t.Error(fmt.Sprintf("Couldn't migrate %s", err))
	}

	// s_a := make([]*Something, 10)
	// for s := range s_a {
	// 	s_a[s] = &Something{Data: "datas"}
	// }
	s := Something{Data: "dddddjfajkfjaf"}
	// s_j, _ := json.Marshal(s)
	obj := DBTESTS{Name: "JSONTST", Something: &s}

	if err := db.Debug().Create(&obj).Error; err != nil {
		t.Error(fmt.Sprintf("Could not create %s", err))
	}
	dbt2 := DBTESTS{}
	if err2 := db.Debug().Where("id = ?", obj.ID).First(&dbt2).Error; err2 != nil || dbt2.Something.Data != obj.Something.Data {
		t.Error(fmt.Sprintf("Could get %s, %+v", err2, dbt2.Something))
	}
}

func TestConcurrent(t *testing.T) {
	db, _ := Open()

	var i uint64 = 0
	for i = 1; i < 9000; i++ {
		go q(db, i)
	}
	i++
	q(db, i)
}

func q(db *gorm.DB, i uint64) {
	e := db.Exec("SELECT * FROM dbtests").Error
	count++
	fmt.Printf("i: %d", count)
	if e != nil {
		panic(fmt.Sprintf("err: %s", e))
	}
}

func TestRawInsert(t *testing.T) {
	db, _ := Open()

	//make sure we've got a table to play with
	db.AutoMigrate(&DBTESTS{})
	fields := []string{"name", "number"}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?, ?)", "dbtests", strings.Join(fields, ", "))
	if db.Exec(query, "TESTNOW", 33).Error != nil {
		t.Error("Insert didn't work")
	}
}

func BenchmarkInsert(b *testing.B) {
	db, _ := Open()

	//make sure we've got a table to play with
	db.AutoMigrate(&DBTESTS{})

	tx := db.Begin()
	defer tx.Commit()

	fields := []string{"name", "number", "thing", "thing1", "thing2", "thing3", "thing4", "thing5", "thing6", "thing7", "thing8", "thing9", "thing10", "thing11", "thing12", "thing13", "thing14", "thing15", "thing16", "thing17", "thing18", "thing19", "thing20"}
	count := 0
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", "dbtests", strings.Join(fields, ", "))
	values := make([]interface{}, 0)
	var e error
	for i := 0; i < 15000; i++ {
		count++
		if count > 1 {
			query += ", "
		}
		query += " (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		values = append(values, fmt.Sprintf("Name_%d", i), i, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
		if i >= 999 { //mssql won't let you do more than 1k rows per insert...
			e = db.Exec(query, values...).Error
			query = fmt.Sprintf("INSERT INTO %s (%s) VALUES ", "dbtests", strings.Join(fields, ", "))
			values = make([]interface{}, 0)
			count = 0
		}
	}
	if count > 0 {
		e = db.Exec(query, values...).Error
	}
	if e != nil {
		tx.Rollback()
		panic(fmt.Sprintf("%s", e))
	}
}
