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
			service, err := CreateService("_/debian")
			Expect(err).NotTo(HaveOccurred())
			Expect(service).ToNot(BeNil())
			Expect(service.State).To(Equal(Initializing))
		})
	})
	Context("Latching a Service onto running Container", func() {

		It("should find ImageID,ContainerID and Tag", func() {
			service, err := CreateService("busybox")
			Expect(err).NotTo(HaveOccurred())
			err = service.Latch()
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Current).ToNot(BeNil())
			Expect(service.Current.Container).ToNot(BeNil())
			Expect(service.State).To(Equal(Running))

		})
	})

	Context("Docker daemon communication", func() {
	})
})
