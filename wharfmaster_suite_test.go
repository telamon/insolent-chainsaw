package wharfmaster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wharfmaster "github.com/telamon/wharfmaster"
	util "github.com/telamon/wharfmaster/util"
	"os"
	"testing"
)

var app = wharfmaster.New()

func TestWharfmaster(t *testing.T) {
	BeforeSuite(func() {
		//	dio.Initialize()
		os.Setenv("GO_ENV", "testing")
	})
	AfterSuite(func() {
		util.StopNginx()
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wharfmaster Suite")
}
