package wharfmaster_test

import (
	. "github.com/kr/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/telamon/wharfmaster/models"
	. "github.com/telamon/wharfmaster/util"
)

var _ = Describe("Service", func() {
	Context("The service state-machine", func() {
		It("should start in state `init`", func() {
			Println("")
			service := CreateService("magical-dishwasher")
			Expect(service).ToNot(BeNil())
			Expect(service.Current).To(BeNil())
			Expect(service.State).To(Equal(Initializing))
		})
	})
	Context("Latching a Service onto running Container", func() {

		It("Should find and set Container and Image", func() {
			service, err := CreateServiceFromImage("busybox")
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Current).ToNot(BeNil())
			Expect(service.Current.Container).ToNot(BeNil())
			Expect(service.Current.Image).ToNot(BeNil())
			Expect(service.State).To(Equal(Running))
		})
		It("Should go to into borked state if container isn't running", func() {
			service, err := CreateServiceFromImage("")
			Expect(err).To(HaveOccurred())
			Expect(service.State).To(Equal(Borked))
		})
		It("should know how to clone container confiuration", func() {
			service, err := CreateServiceFromImage("telamon/wharftest-version:1.1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Current).ToNot(BeNil())
			//Println(service.Current.Container.Config)
			conf := service.Current.MigrateConfigTo("1.2.0")

			Expect(conf.Config.Image).To(Equal("telamon/wharftest-version:1.2.0"))
			Expect(service.Current.Container.Config.Image).To(Equal("telamon/wharftest-version:1.1.0"))
		})
		PIt("should go into state `pulling` when signaled to redeploy", func() {
			service, err := CreateServiceFromImage("telamon/wharftest-version:1.1.0")
			Expect(err).NotTo(HaveOccurred())
			service.Redeploy("1.2.0")
			Expect(service.State == Pulling)
		})
	})
	Context("External interactions", func() {
		It("should generate a new configuration", func() {
			_, err := RegenerateConf()
			Expect(err).ToNot(HaveOccurred())
		})
		It("should be able reload nginx", func() {
			workers, err := NginxWorkerPids()
			Expect(err).ToNot(HaveOccurred())
			c := make(chan error)
			go func() {
				c <- ReloadNginx()
			}()
			err = <-c
			Expect(err).ToNot(HaveOccurred())
			workersNew, err := NginxWorkerPids()
			Expect(err).ToNot(HaveOccurred())
			for _, op := range workers {
				for _, np := range workersNew {
					Expect(op).NotTo(Equal(np))
				}
			}
		})
	})

})
