package streams

import (
	"bytes"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type (
	// RootModuleBlob is the module that will be registered with the runtime.
	RootModuleBlob struct{}

	// ModuleInstanceBlob is the module instance that will be created for each VU.
	ModuleInstanceBlob struct {
		vu modules.VU
	}
)

// Ensure the interfaces are implemented correctly
var (
	_ modules.Instance = &ModuleInstanceBlob{}
	_ modules.Module   = &RootModuleBlob{}
)

// NewBlob creates a new RootModuleBlob instance.
func NewBlob() *RootModuleBlob {
	return &RootModuleBlob{}
}

// NewModuleInstance creates a new instance of the module for a specific VU.
func (rm *RootModuleBlob) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstanceBlob{
		vu: vu,
	}
}

// Exports returns the module exports, that will be available in the runtime.
func (mi *ModuleInstanceBlob) Exports() modules.Exports {
	return modules.Exports{Named: map[string]interface{}{
		"Blob": mi.NewBlob,
	}}
}

func (mi *ModuleInstanceBlob) NewBlob(call sobek.ConstructorCall) *sobek.Object {
	rt := mi.vu.Runtime()

	// Validate constructor call arguments.
	//if len(call.Arguments) > 0 && !sobek.IsUndefined(call.Arguments[0]) {
	//	if !isObject(call.Arguments[0]) {
	//		throw(rt, newTypeError(rt, "first argument must be an object"))
	//	}
	//}

	b := &Blob{}
	if len(call.Arguments) > 0 {
		if parts, ok := call.Arguments[0].Export().([]interface{}); ok {
			for _, part := range parts {
				var err error
				switch v := part.(type) {
				case []uint8:
					_, err = b.data.Write(v)
				case string:
					_, err = b.data.WriteString(v)
				}
				if err != nil {
					common.Throw(rt, newError(RuntimeError, err.Error()))
				}
			}
		}
	}

	obj := rt.NewObject()

	if err := obj.DefineAccessorProperty("type", rt.ToValue(func() sobek.Value {
		return rt.ToValue(b.typ)
	}), nil, sobek.FLAG_FALSE, sobek.FLAG_TRUE); err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	if err := obj.DefineAccessorProperty("size", rt.ToValue(func() sobek.Value {
		return rt.ToValue(b.data.Len())
	}), nil, sobek.FLAG_FALSE, sobek.FLAG_TRUE); err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	if err := obj.Set("text", func(call sobek.FunctionCall) sobek.Value {
		return rt.ToValue(b.text())
	}); err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	if err := obj.Set("arrayBuffer", func(call sobek.FunctionCall) sobek.Value {
		return rt.ToValue(b.arrayBuffer())
	}); err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	proto := call.This.Prototype()
	err := obj.SetPrototype(proto)
	if err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	if err := proto.Set("toString", func(call sobek.FunctionCall) sobek.Value {
		return rt.ToValue("[object Blob]")
	}); err != nil {
		common.Throw(rt, newError(RuntimeError, err.Error()))
	}

	return obj
}

type Blob struct {
	typ  string
	data bytes.Buffer
}

func (b *Blob) text() string {
	return b.data.String()
}

func (b *Blob) arrayBuffer() []byte {
	return b.data.Bytes()
}
