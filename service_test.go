package wharfmaster_test

import (
	. "github.com/telamon/wharfmaster/util"

	. "fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/telamon/wharfmaster/models"
)

var _ = Describe("Service", func() {
	Context("The service state-machine", func() {
		It("should start in state `init`", func() {
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
	})
	Context("Docker-gen interaction", func() {
		It("should generate a new configuration", func() {
			_, err := RegenerateConf()
			Expect(err).ToNot(HaveOccurred())
		})
		PIt("should be able to start and stop nginx", func() {
			_, err := RegenerateConf()
			Expect(err).ToNot(HaveOccurred())
			Println("Attempting to start nginx")
			err = StartNginx()
			Expect(NginxPID()).ToNot(Equal(-1))
			Println("Attempting to stop nginx")
			StopNginx()
		})
	})

})
