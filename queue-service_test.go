package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var response = `{"success":true,"outputScenario":"Default","data":{"review_queues":[{"queue":{"count":"9.3k","name":"Close Votes"}},{"queue":{"count":"196","name":"Triage"}},{"queue":{"count":"160","name":"First Posts"}},{"queue":{"count":"157","name":"Suggested Edits"}},{"queue":{"count":"151","name":"Help and Improvement"}},{"queue":{"count":"125","name":"Reopen Votes"}},{"queue":{"count":"75","name":"Low Quality Posts"}},{"queue":{"count":"3","name":"Late Answers"}},{"queue":{"count":"0","name":"Documentation: Proposed Changes"}},{"queue":{"count":null,"name":"Meta Reviews"}}]},"stateToken":"eyJqYXIiOnsidmVyc2lvbiI6InRvdWdoLWNvb2tpZUAyLjMuMiIsInN0b3JlVHlwZSI6Ik1lbW9yeUNvb2tpZVN0b3JlIiwicmVqZWN0UHVibGljU3VmZml4ZXMiOnRydWUsImNvb2tpZXMiOlt7ImtleSI6InByb3YiLCJ2YWx1ZSI6IjkyZmVjNTI1LTVmNDctNjg4My1hN2VkLTUzNzdkNjc2ZTFkMCIsImV4cGlyZXMiOiIyMDU1LTAxLTAxVDAwOjAwOjAwLjAwMFoiLCJkb21haW4iOiJzdGFja292ZXJmbG93LmNvbSIsInBhdGgiOiIvIiwiaHR0cE9ubHkiOnRydWUsImhvc3RPbmx5IjpmYWxzZSwiY3JlYXRpb24iOiIyMDE3LTA1LTAxVDIwOjQwOjA2LjY4M1oiLCJsYXN0QWNjZXNzZWQiOiIyMDE3LTA1LTAxVDIwOjQwOjA2LjY4M1oifV19LCJ2ZXJzaW9uIjoxfQ=="}`

func TestResponseUnmarshal(t *testing.T) {
	mockClient := NewMockHttpOperations(gomock.NewController(t))
	mockClient.EXPECT().Get(url).Times(1).Return([]byte(response), nil)

	s := NewQueueService(mockClient)

	queues, err := s.GetQueues()

	assert.Equal(t, 10, len(queues))
	assert.Nil(t, err)
}
