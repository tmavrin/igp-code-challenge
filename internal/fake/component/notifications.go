// Code generated by counterfeiter. DO NOT EDIT.
package component

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/tmavrin/igp-code-challenge/internal/component/notifications"
)

type NotificationsProvider struct {
	RegisterClientStub        func(uuid.UUID, *websocket.Conn)
	registerClientMutex       sync.RWMutex
	registerClientArgsForCall []struct {
		arg1 uuid.UUID
		arg2 *websocket.Conn
	}
	SendNotificationToUserStub        func(uuid.UUID, string)
	sendNotificationToUserMutex       sync.RWMutex
	sendNotificationToUserArgsForCall []struct {
		arg1 uuid.UUID
		arg2 string
	}
	UnregisterClientStub        func(*websocket.Conn)
	unregisterClientMutex       sync.RWMutex
	unregisterClientArgsForCall []struct {
		arg1 *websocket.Conn
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *NotificationsProvider) RegisterClient(arg1 uuid.UUID, arg2 *websocket.Conn) {
	fake.registerClientMutex.Lock()
	fake.registerClientArgsForCall = append(fake.registerClientArgsForCall, struct {
		arg1 uuid.UUID
		arg2 *websocket.Conn
	}{arg1, arg2})
	stub := fake.RegisterClientStub
	fake.recordInvocation("RegisterClient", []interface{}{arg1, arg2})
	fake.registerClientMutex.Unlock()
	if stub != nil {
		fake.RegisterClientStub(arg1, arg2)
	}
}

func (fake *NotificationsProvider) RegisterClientCallCount() int {
	fake.registerClientMutex.RLock()
	defer fake.registerClientMutex.RUnlock()
	return len(fake.registerClientArgsForCall)
}

func (fake *NotificationsProvider) RegisterClientCalls(stub func(uuid.UUID, *websocket.Conn)) {
	fake.registerClientMutex.Lock()
	defer fake.registerClientMutex.Unlock()
	fake.RegisterClientStub = stub
}

func (fake *NotificationsProvider) RegisterClientArgsForCall(i int) (uuid.UUID, *websocket.Conn) {
	fake.registerClientMutex.RLock()
	defer fake.registerClientMutex.RUnlock()
	argsForCall := fake.registerClientArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *NotificationsProvider) SendNotificationToUser(arg1 uuid.UUID, arg2 string) {
	fake.sendNotificationToUserMutex.Lock()
	fake.sendNotificationToUserArgsForCall = append(fake.sendNotificationToUserArgsForCall, struct {
		arg1 uuid.UUID
		arg2 string
	}{arg1, arg2})
	stub := fake.SendNotificationToUserStub
	fake.recordInvocation("SendNotificationToUser", []interface{}{arg1, arg2})
	fake.sendNotificationToUserMutex.Unlock()
	if stub != nil {
		fake.SendNotificationToUserStub(arg1, arg2)
	}
}

func (fake *NotificationsProvider) SendNotificationToUserCallCount() int {
	fake.sendNotificationToUserMutex.RLock()
	defer fake.sendNotificationToUserMutex.RUnlock()
	return len(fake.sendNotificationToUserArgsForCall)
}

func (fake *NotificationsProvider) SendNotificationToUserCalls(stub func(uuid.UUID, string)) {
	fake.sendNotificationToUserMutex.Lock()
	defer fake.sendNotificationToUserMutex.Unlock()
	fake.SendNotificationToUserStub = stub
}

func (fake *NotificationsProvider) SendNotificationToUserArgsForCall(i int) (uuid.UUID, string) {
	fake.sendNotificationToUserMutex.RLock()
	defer fake.sendNotificationToUserMutex.RUnlock()
	argsForCall := fake.sendNotificationToUserArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *NotificationsProvider) UnregisterClient(arg1 *websocket.Conn) {
	fake.unregisterClientMutex.Lock()
	fake.unregisterClientArgsForCall = append(fake.unregisterClientArgsForCall, struct {
		arg1 *websocket.Conn
	}{arg1})
	stub := fake.UnregisterClientStub
	fake.recordInvocation("UnregisterClient", []interface{}{arg1})
	fake.unregisterClientMutex.Unlock()
	if stub != nil {
		fake.UnregisterClientStub(arg1)
	}
}

func (fake *NotificationsProvider) UnregisterClientCallCount() int {
	fake.unregisterClientMutex.RLock()
	defer fake.unregisterClientMutex.RUnlock()
	return len(fake.unregisterClientArgsForCall)
}

func (fake *NotificationsProvider) UnregisterClientCalls(stub func(*websocket.Conn)) {
	fake.unregisterClientMutex.Lock()
	defer fake.unregisterClientMutex.Unlock()
	fake.UnregisterClientStub = stub
}

func (fake *NotificationsProvider) UnregisterClientArgsForCall(i int) *websocket.Conn {
	fake.unregisterClientMutex.RLock()
	defer fake.unregisterClientMutex.RUnlock()
	argsForCall := fake.unregisterClientArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NotificationsProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.registerClientMutex.RLock()
	defer fake.registerClientMutex.RUnlock()
	fake.sendNotificationToUserMutex.RLock()
	defer fake.sendNotificationToUserMutex.RUnlock()
	fake.unregisterClientMutex.RLock()
	defer fake.unregisterClientMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *NotificationsProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ notifications.Provider = new(NotificationsProvider)
