package validators_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
	authvalidator "grovehub/services/identity/validators/impl"
)

var _ = Describe("Registration DTO stages", func() {
	var validator = authvalidator.NewAuthValidator()

	It("trims name fields and canonicalizes email before validating tags", func() {
		middleName := "  Devi "
		request := identityrequest.Register{
			FirstName: "  Sima  ", MiddleName: &middleName, LastName: "  Bharti ",
			Email: "  SIMA@example.com  ", Password: "securePassword123",
		}
		validated, err := validator.ValidateRegister(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(validated.FirstName).To(Equal("Sima"))
		Expect(*validated.MiddleName).To(Equal("Devi"))
		Expect(validated.LastName).To(Equal("Bharti"))
		Expect(validated.Email).To(Equal("sima@example.com"))
		Expect(*request.MiddleName).To(Equal("  Devi "))
	})

	It("uses DTO tags for required names and explicit nonnull middle name", func() {
		middleName := "ok"
		request := identityrequest.Register{
			FirstName: "   ", MiddleName: &middleName, LastName: "Bharti",
			Email: "sima@example.com", Password: "securePassword123",
		}
		_, err := validator.ValidateRegister(request)
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(identityerror.InvalidRegistration.Code()))
		request.FirstName = "Sima"
		request.MiddleName = nil
		_, err = validator.ValidateRegister(request)
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(identityerror.InvalidRegistration.Code()))
	})

	It("runs password checks as custom domain validation", func() {
		middleName := "ok"
		request := identityrequest.Register{
			FirstName: "Sima", MiddleName: &middleName, LastName: "Bharti",
			Email: "sima@example.com", Password: "short",
		}
		_, err := validator.ValidateRegister(request)
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(identityerror.InvalidRegistration.Code()))
	})
})
