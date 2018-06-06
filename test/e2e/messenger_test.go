package e2e_test

import (
	"os"

	"github.com/kubeware/messenger/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	api "github.com/kubeware/messenger/apis/messenger/v1alpha1"
)

var _ = Describe("Messenger", func() {
	var (
		f *framework.Invocation

		labels          map[string]string
		name, namespace string

		secret, notifierConfig *core.Secret
		service, svc     *core.Service

		messagingServiceObj, messagingService *api.MessagingService
		messageObj, message *api.Message

		drive string
		to []string
		err              error

		authTokenToSendMessage     string
		authTokenToSeeHistory string
		skipSend, skipSeeHist bool
	)

	BeforeEach(func() {
		f = root.Invoke()
		name = f.App()
		namespace = f.Namespace()
		labels = map[string]string{
			"app": f.App(),
		}

		skipSend = false
		skipSeeHist = false

	})

	Describe("Send Message in Hipchat", func() {
		BeforeEach(func() {
			authTokenToSendMessage, skipSend = os.LookupEnv("AUTH_TOKEN_TO_SEND_MSG")
			if skipSend	{
				secret = f.NewSecret(name+"-secret-notifier-config", namespace, authTokenToSendMessage, labels)
			}

			authTokenToSeeHistory, skipSeeHist = os.LookupEnv("AUTH_TOKEN_TO_SEE_HIST")

			drive = "Hipchat"
			to = []string{"ops-alerts"}
			messagingServiceObj = f.NewMessagingService(name, namespace, labels, drive, secret.Name, to)
		})

		JustBeforeEach(func() {
			if skipSend	{
				By("Creating secret...")
				notifierConfig, err := f.CreateSecret(secret)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			By("Deleting secrets...")
			f.DeleteAllSecrets()
		})

		Context("To \"ops-alerts\" room", func() {

		})
	})
})
