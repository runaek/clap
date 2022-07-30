//go:build clap_mocks
// +build clap_mocks

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/runaek/clap (interfaces: FileReader,FileWriter,FileInfo,PositionalDeriver,KeyValueDeriver,FlagDeriver)

// Package clap is a generated GoMock package.
package clap

import (
	fs "io/fs"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockFileReader is a mock of FileReader interface.
type MockFileReader struct {
	ctrl     *gomock.Controller
	recorder *MockFileReaderMockRecorder
}

// MockFileReaderMockRecorder is the mock recorder for MockFileReader.
type MockFileReaderMockRecorder struct {
	mock *MockFileReader
}

// NewMockFileReader creates a new mock instance.
func NewMockFileReader(ctrl *gomock.Controller) *MockFileReader {
	mock := &MockFileReader{ctrl: ctrl}
	mock.recorder = &MockFileReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileReader) EXPECT() *MockFileReaderMockRecorder {
	return m.recorder
}

// Fd mocks base method.
func (m *MockFileReader) Fd() uintptr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fd")
	ret0, _ := ret[0].(uintptr)
	return ret0
}

// Fd indicates an expected call of Fd.
func (mr *MockFileReaderMockRecorder) Fd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fd", reflect.TypeOf((*MockFileReader)(nil).Fd))
}

// Read mocks base method.
func (m *MockFileReader) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockFileReaderMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockFileReader)(nil).Read), arg0)
}

// Stat mocks base method.
func (m *MockFileReader) Stat() (fs.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(fs.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockFileReaderMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockFileReader)(nil).Stat))
}

// MockFileWriter is a mock of FileWriter interface.
type MockFileWriter struct {
	ctrl     *gomock.Controller
	recorder *MockFileWriterMockRecorder
}

// MockFileWriterMockRecorder is the mock recorder for MockFileWriter.
type MockFileWriterMockRecorder struct {
	mock *MockFileWriter
}

// NewMockFileWriter creates a new mock instance.
func NewMockFileWriter(ctrl *gomock.Controller) *MockFileWriter {
	mock := &MockFileWriter{ctrl: ctrl}
	mock.recorder = &MockFileWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileWriter) EXPECT() *MockFileWriterMockRecorder {
	return m.recorder
}

// Fd mocks base method.
func (m *MockFileWriter) Fd() uintptr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fd")
	ret0, _ := ret[0].(uintptr)
	return ret0
}

// Fd indicates an expected call of Fd.
func (mr *MockFileWriterMockRecorder) Fd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fd", reflect.TypeOf((*MockFileWriter)(nil).Fd))
}

// Stat mocks base method.
func (m *MockFileWriter) Stat() (fs.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(fs.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockFileWriterMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockFileWriter)(nil).Stat))
}

// Write mocks base method.
func (m *MockFileWriter) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockFileWriterMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockFileWriter)(nil).Write), arg0)
}

// MockFileInfo is a mock of FileInfo interface.
type MockFileInfo struct {
	ctrl     *gomock.Controller
	recorder *MockFileInfoMockRecorder
}

// MockFileInfoMockRecorder is the mock recorder for MockFileInfo.
type MockFileInfoMockRecorder struct {
	mock *MockFileInfo
}

// NewMockFileInfo creates a new mock instance.
func NewMockFileInfo(ctrl *gomock.Controller) *MockFileInfo {
	mock := &MockFileInfo{ctrl: ctrl}
	mock.recorder = &MockFileInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileInfo) EXPECT() *MockFileInfoMockRecorder {
	return m.recorder
}

// IsDir mocks base method.
func (m *MockFileInfo) IsDir() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDir")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDir indicates an expected call of IsDir.
func (mr *MockFileInfoMockRecorder) IsDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDir", reflect.TypeOf((*MockFileInfo)(nil).IsDir))
}

// ModTime mocks base method.
func (m *MockFileInfo) ModTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ModTime indicates an expected call of ModTime.
func (mr *MockFileInfoMockRecorder) ModTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModTime", reflect.TypeOf((*MockFileInfo)(nil).ModTime))
}

// Mode mocks base method.
func (m *MockFileInfo) Mode() fs.FileMode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mode")
	ret0, _ := ret[0].(fs.FileMode)
	return ret0
}

// Mode indicates an expected call of Mode.
func (mr *MockFileInfoMockRecorder) Mode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mode", reflect.TypeOf((*MockFileInfo)(nil).Mode))
}

// Name mocks base method.
func (m *MockFileInfo) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFileInfoMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFileInfo)(nil).Name))
}

// Size mocks base method.
func (m *MockFileInfo) Size() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockFileInfoMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockFileInfo)(nil).Size))
}

// Sys mocks base method.
func (m *MockFileInfo) Sys() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sys")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Sys indicates an expected call of Sys.
func (mr *MockFileInfoMockRecorder) Sys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sys", reflect.TypeOf((*MockFileInfo)(nil).Sys))
}

// MockPositionalDeriver is a mock of PositionalDeriver interface.
type MockPositionalDeriver struct {
	ctrl     *gomock.Controller
	recorder *MockPositionalDeriverMockRecorder
}

// MockPositionalDeriverMockRecorder is the mock recorder for MockPositionalDeriver.
type MockPositionalDeriverMockRecorder struct {
	mock *MockPositionalDeriver
}

// NewMockPositionalDeriver creates a new mock instance.
func NewMockPositionalDeriver(ctrl *gomock.Controller) *MockPositionalDeriver {
	mock := &MockPositionalDeriver{ctrl: ctrl}
	mock.recorder = &MockPositionalDeriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPositionalDeriver) EXPECT() *MockPositionalDeriverMockRecorder {
	return m.recorder
}

// DerivePosition mocks base method.
func (m *MockPositionalDeriver) DerivePosition(arg0 interface{}, arg1 int, arg2 ...Option) (IPositional, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DerivePosition", varargs...)
	ret0, _ := ret[0].(IPositional)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DerivePosition indicates an expected call of DerivePosition.
func (mr *MockPositionalDeriverMockRecorder) DerivePosition(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DerivePosition", reflect.TypeOf((*MockPositionalDeriver)(nil).DerivePosition), varargs...)
}

// MockKeyValueDeriver is a mock of KeyValueDeriver interface.
type MockKeyValueDeriver struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueDeriverMockRecorder
}

// MockKeyValueDeriverMockRecorder is the mock recorder for MockKeyValueDeriver.
type MockKeyValueDeriverMockRecorder struct {
	mock *MockKeyValueDeriver
}

// NewMockKeyValueDeriver creates a new mock instance.
func NewMockKeyValueDeriver(ctrl *gomock.Controller) *MockKeyValueDeriver {
	mock := &MockKeyValueDeriver{ctrl: ctrl}
	mock.recorder = &MockKeyValueDeriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValueDeriver) EXPECT() *MockKeyValueDeriverMockRecorder {
	return m.recorder
}

// DeriveKeyValue mocks base method.
func (m *MockKeyValueDeriver) DeriveKeyValue(arg0 interface{}, arg1 string, arg2 ...Option) (IKeyValue, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeriveKeyValue", varargs...)
	ret0, _ := ret[0].(IKeyValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeriveKeyValue indicates an expected call of DeriveKeyValue.
func (mr *MockKeyValueDeriverMockRecorder) DeriveKeyValue(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeriveKeyValue", reflect.TypeOf((*MockKeyValueDeriver)(nil).DeriveKeyValue), varargs...)
}

// MockFlagDeriver is a mock of FlagDeriver interface.
type MockFlagDeriver struct {
	ctrl     *gomock.Controller
	recorder *MockFlagDeriverMockRecorder
}

// MockFlagDeriverMockRecorder is the mock recorder for MockFlagDeriver.
type MockFlagDeriverMockRecorder struct {
	mock *MockFlagDeriver
}

// NewMockFlagDeriver creates a new mock instance.
func NewMockFlagDeriver(ctrl *gomock.Controller) *MockFlagDeriver {
	mock := &MockFlagDeriver{ctrl: ctrl}
	mock.recorder = &MockFlagDeriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlagDeriver) EXPECT() *MockFlagDeriverMockRecorder {
	return m.recorder
}

// DeriveFlag mocks base method.
func (m *MockFlagDeriver) DeriveFlag(arg0 interface{}, arg1 string, arg2 ...Option) (IFlag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeriveFlag", varargs...)
	ret0, _ := ret[0].(IFlag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeriveFlag indicates an expected call of DeriveFlag.
func (mr *MockFlagDeriverMockRecorder) DeriveFlag(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeriveFlag", reflect.TypeOf((*MockFlagDeriver)(nil).DeriveFlag), varargs...)
}