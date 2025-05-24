package container

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func Build(engine *gin.Engine) (*di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	// Register base Gin engine as a dependency
	builder.Add((di.Def{
		Name: "engine",
		Build: func(ctn di.Container) (interface{}, error) {
			return engine, nil
		},
	}))

	// Register core/shared dependencies
	if _, err := RegisterDependency(builder); err != nil {
		return nil, err
	}

	// Register feature/module-specific
	RegisterModul(builder)

	// Build the container
	ctn := builder.Build()
	var (
		_ = ctn.Get("user-handler")
		_ = ctn.Get("auth-handler")
		_ = ctn.Get("auth-middleware")
	)

	return &ctn, nil
}
