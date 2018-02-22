package interfaces

import (
	"os"
	"testing"
	"github.com/go-gorp/gorp"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
)

type DatasourcesTestSuite struct {
	suite.Suite
}

// Go testing framework entry point.
func TestDatasourcesTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourcesTestSuite))
}

func (suite *DatasourcesTestSuite) TestInitAlbaDatasource() {
	// Rewrite datasource config for test.
	testDataSourceFile := os.TempDir() + "testDataSource.db"
	viper.Set("DB.file", testDataSourceFile)
	viper.SetDefault("DB.Driver", "sqlite3")

	// Test creating the datasource.
	ds, err := InitAlbaDatasource()
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), ds)

	// Close datasource and remove test file.
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		dbmap.Db.Close()
	}

	// Test with an invalid path.
	viper.Set("DB.file", "/whatever/test.db")
	viper.SetDefault("DB.Driver", "sqlite3")

	// Test creating the datasource.
	ds, err = InitAlbaDatasource()
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), ds)

	// Close datasource and remove test file.
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		dbmap.Db.Close()
	}

	// Test with an invalid driver.
	viper.Set("DB.file", testDataSourceFile)
	viper.SetDefault("DB.Driver", "")

	// Test creating the datasource.
	ds, err = InitAlbaDatasource()
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), ds)

	// Close datasource and remove test file.
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		dbmap.Db.Close()
	}
}
