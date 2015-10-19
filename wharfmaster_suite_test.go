package wharfmaster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/telamon/wharfmaster"
	"os"
	//dio "github.com/telamon/wharfmaster/util/dio"
	"testing"
)

var app = wharfmaster.New()

func TestWharfmaster(t *testing.T) {
	BeforeSuite(func() {
		//	dio.Initialize()
		os.Setenv("GO_ENV", "testing")
	})
	AfterSuite(func() {
		wharfmaster.StopNginx()
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wharfmaster Suite")
}
