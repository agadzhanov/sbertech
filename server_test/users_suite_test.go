package server_test_test

import (
	"context"
	"crud/api/users"
	"crud/db_stub"
	"crud/logger"
	"crud/server"
	"crud/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"reflect"
	"testing"
)

var (
	usersServer users.UsersServer
)

var _ = BeforeSuite(func() {
	var (
		err          error
		log          *logrus.Entry
		usersStorage storage.Users
	)
	usersStorage, err = db_stub.NewUsersStorage()
	Expect(err).To(BeNil())

	// получаем логгер
	log, err = logger.NewLogger()
	Expect(err).NotTo(HaveOccurred())
	// глушим, чтобы при запуске тестов не спамил в командной строке
	log.Logger.SetOutput(ioutil.Discard)
	Expect(err).To(BeNil())
	Expect(log).NotTo(BeNil())

	usersServerOptions := []interface{}{
		log,
		usersStorage,
	}
	usersServer = server.NewUsersServer(usersServerOptions)
})

func TestUsers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Suite")
}

var _ = Describe("UsersServer.CreateUser", func() {
	var (
		err            error
		createResponse *users.CreateUserResponse
		createRequest  *users.CreateUserRequest
	)

	Context("Request field `User` is nil", func() {
		BeforeEach(func() {
			createRequest = &users.CreateUserRequest{}
			createResponse, err = usersServer.CreateUser(context.Background(), createRequest)
		})
		When("Request is executed", func() {
			It("Response should be nil", func() {
				Expect(createResponse).To(BeNil())
			})

			It("Status code should be `InvalidArgument`", func() {
				actualCode := status.Code(err)
				Expect(actualCode).To(Equal(codes.InvalidArgument))
			})

			It("should cause an error `Invalid Argument`", func() {
				actualCode := status.Code(err)
				//Expect(err.Error()).To(Equal("rpc error: code = Internal desc = request failed"))
				Expect(actualCode).To(Equal(codes.InvalidArgument))
			})
		})
	})

	Context("Request field User.Id is not empty", func() {
		BeforeEach(func() {
			createRequest = &users.CreateUserRequest{
				User: &users.User{
					Id: 100500,
				},
			}
			createResponse, err = usersServer.CreateUser(context.Background(), createRequest)
		})
		When("Request is executed", func() {
			It("Status code should be AlreadyExists", func() {
				Expect(status.Code(err)).To(Equal(codes.AlreadyExists))
			})

			It("Response should be nil", func() {
				Expect(createResponse).To(BeNil())
			})
		})
	})

	var (
		expectedFirstName = "some name"
		expectedLastName  = "some last name"
	)

	Context("Request field User.Id is  empty", func() {
		var previousCreatedUserId uint64
		BeforeEach(func() {
			createRequest = &users.CreateUserRequest{
				User: &users.User{
					Id:        0,
					FirstName: expectedFirstName,
					LastName:  expectedLastName,
				},
			}
			previousCreatedUserId = createResponse.GetUser().GetId()
			createResponse, err = usersServer.CreateUser(context.Background(), createRequest)
		})
		When("Request is executed", func() {
			It("Status code should be `OK`", func() {
				Expect(status.Code(err)).To(Equal(codes.OK))
			})

			It("Response has User field with valid data", func() {
				Expect(createResponse.GetUser().GetId()).NotTo(BeNil())
				Expect(createResponse.GetUser().GetFirstName()).To(Equal(expectedFirstName))
				Expect(createResponse.GetUser().GetLastName()).To(Equal(expectedLastName))
			})
		})

		When("Request is executed again", func() {
			It("`Id` field of last created user is greater, then previous", func() {
				Expect(createResponse.GetUser().GetId() > previousCreatedUserId).To(BeTrue())
			})
		})
	})
})

var _ = Describe("UsersServer.GetUser", func() {
	Context("Getting existing user", func() {
		var (
			firstName    = "fName"
			lastName     = "lName"
			createdUser  *users.User
			expectedUser *users.User
		)

		BeforeEach(func() {
			createRequest := &users.CreateUserRequest{
				User: &users.User{
					FirstName: firstName,
					LastName:  lastName,
				},
			}
			response, err := usersServer.CreateUser(context.Background(), createRequest)
			Expect(err).NotTo(HaveOccurred())
			createdUser = response.GetUser()
			Expect(createdUser.GetId()).NotTo(BeZero())
			expectedUser = &users.User{
				Id:        createdUser.GetId(),
				FirstName: firstName,
				LastName:  lastName,
			}
		})
		When("User is created", func() {
			var actualUser *users.User

			BeforeEach(func() {
				getUserRequest := &users.GetUserRequest{
					Id: createdUser.GetId(),
				}
				response, err := usersServer.GetUser(context.Background(), getUserRequest)
				Expect(err).NotTo(HaveOccurred())
				actualUser = response.GetUser()
			})
			It("User, obtained by Id, should contain expected data", func() {
				Expect(reflect.DeepEqual(expectedUser, actualUser)).To(BeTrue())
			})
		})
	})

	Context("Getting not existing user", func() {
		var (
			response *users.GetUserResponse
			err      error
		)
		BeforeEach(func() {
			getUserRequest := &users.GetUserRequest{
				Id: 0,
			}
			response, err = usersServer.GetUser(context.Background(), getUserRequest)
		})
		It("Status code is NotFound", func() {
			Expect(response).To(BeNil())
			statusCode := status.Code(err)
			Expect(statusCode).To(Equal(codes.NotFound))
		})
	})
})

var _ = Describe("UsersServer.UpdateUser", func() {
	Context("Updating existing user", func() {
		var (
			firstName    = "fName"
			lastName     = "lName"
			createdUser  *users.User
			UpdatingUser *users.User
		)

		BeforeEach(func() {
			createRequest := &users.CreateUserRequest{
				User: &users.User{
					FirstName: firstName,
					LastName:  lastName,
				},
			}
			response, err := usersServer.CreateUser(context.Background(), createRequest)
			Expect(err).NotTo(HaveOccurred())
			createdUser = response.GetUser()
			Expect(createdUser.GetId()).NotTo(BeZero())
			UpdatingUser = &users.User{
				Id:        createdUser.GetId(),
				FirstName: firstName + "updated",
				LastName:  lastName + "updated",
			}
		})
		When("User is created", func() {
			var actualUser *users.User

			BeforeEach(func() {
				updateUserRequest := &users.UpdateUserRequest{
					User: UpdatingUser,
				}
				response, err := usersServer.UpdateUser(context.Background(), updateUserRequest)
				Expect(err).NotTo(HaveOccurred())
				actualUser = response.GetUser()
			})
			It("Updated user from response body, should contain expected data", func() {
				Expect(reflect.DeepEqual(UpdatingUser, actualUser)).To(BeTrue())
			})

			BeforeEach(func() {
				getUserRequest := &users.GetUserRequest{
					Id: actualUser.GetId(),
				}
				response, err := usersServer.GetUser(context.Background(), getUserRequest)
				Expect(err).NotTo(HaveOccurred())
				actualUser = response.GetUser()
			})
			It("Updated user, obtained by Id, should contain expected data", func() {
				Expect(reflect.DeepEqual(UpdatingUser, actualUser)).To(BeTrue())
			})
		})
	})

	Context("Updating not existing user", func() {
		var (
			response *users.UpdateUserResponse
			err      error
		)
		BeforeEach(func() {
			UpdateUserRequest := &users.UpdateUserRequest{
				User: &users.User{},
			}
			response, err = usersServer.UpdateUser(context.Background(), UpdateUserRequest)
		})
		It("Status code is NotFound", func() {
			Expect(response).To(BeNil())
			statusCode := status.Code(err)
			Expect(statusCode).To(Equal(codes.NotFound))
		})
	})
})

var _ = Describe("UsersServer.DeleteUser", func() {
	Context("Deleting existing user", func() {
		var (
			createdUser *users.User
		)

		BeforeEach(func() {
			createRequest := &users.CreateUserRequest{
				User: &users.User{},
			}
			response, err := usersServer.CreateUser(context.Background(), createRequest)
			Expect(err).NotTo(HaveOccurred())
			createdUser = response.GetUser()
			Expect(createdUser.GetId()).NotTo(BeZero())
		})
		When("User is created", func() {
			var (
				err error
			)

			BeforeEach(func() {
				deleteUserRequest := &users.DeleteUserRequest{
					Id: createdUser.GetId(),
				}
				_, err = usersServer.DeleteUser(context.Background(), deleteUserRequest)
			})
			It("Response status code after user deletion is `OK`", func() {
				statusCode := status.Code(err)
				Expect(statusCode).To(Equal(codes.OK))
			})
		})
		Context("Deleting not existing user", func() {
			var err error
			BeforeEach(func() {
				deleteUserRequest := &users.DeleteUserRequest{
					Id: 0,
				}
				_, err = usersServer.DeleteUser(context.Background(), deleteUserRequest)
			})
			It("Response status code after deletion non-existing user is `NotFound`", func() {
				statusCode := status.Code(err)
				Expect(statusCode).To(Equal(codes.NotFound))
			})
		})
	})
})
