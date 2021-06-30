// event package contains structs for events between microservices. this events are used for publish/subscribe pattern.
package event

// Registration represents event when user starter registration.
type Registration struct {
	RequestID string // id of APM transaction
	Login     string
	Email     string
}

// todo add Publisher and Subscriber interfaces
