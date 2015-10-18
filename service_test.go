package wharfmaster_test

import (
	_ "github.com/telamon/wharfmaster"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/ryanfaerman/fsm"
	//docker "github.com/fsouza/go-dockerclient"
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

})
