package image

import (
	cont "context"
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func setupMockClient(t *testing.T) (*client.MockDockerClient, func()) {
	mock := gomock.NewController(t)
	mockClient := client.NewMockDockerClient(mock)
	return mockClient, func() { mock.Finish() }
}

func setupCtxAndConfig(
	mockClient *client.MockDockerClient,
) (*context.ExecuteContext, *config.ImageConfig) {
	ctx := &context.ExecuteContext{
		Client:     mockClient,
		WorkingDir: "/dir",
	}
	config := &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
	return ctx, config
}

func TestTagImageNothingToTag(t *testing.T) {
	ctx := &context.ExecuteContext{}
	config := &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
	err := tagImage(ctx, config, "imagename:tag")
	assert.NilError(t, err)
}

func TestTagImageWithTag(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().ImageTag(cont.Background(), "imagename:tag", "imagename:foo")

	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "imagename:foo")
	assert.NilError(t, err)
}

func TestTagImageWithFullImageName(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().ImageTag(cont.Background(), "imagename:tag", "othername:bar")
	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "othername:bar")
	assert.NilError(t, err)
}

func TestTagImageWithFullImageNameAndHost(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().ImageTag(cont.Background(), "imagename:tag", "localhost:3030/othername:bar")
	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "localhost:3030/othername:bar")
	assert.NilError(t, err)
}
