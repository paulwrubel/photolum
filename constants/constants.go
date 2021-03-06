package constants

import "math"

var PostgresHostnameEnvironmentKey = "PHOTOLUM_PG_HOSTNAME"
var PostgresUsernameEnvironmentKey = "PHOTOLUM_PG_USER"
var PostgresPasswordEnvironmentKey = "PHOTOLUM_PG_PASSWORD"

var ParametersMinimumDimension uint32 = 5
var ParametersMaximumDimension uint32 = 50000
var ParametersMinimumTotalPixels uint32 = 100
var ParametersMaximumTotalPixels uint32 = 25000000
var ParametersMaximumMaxBounces uint32 = 100
var ParametersMaximumTMax float64 = math.MaxFloat64

var CameraMinimumVerticalFOV float64 = 10.0
var CameraMaximumVerticalFOV float64 = 120.0
var CameraMinimumAperture float64 = 0.0
var CameraMaximumAperture float64 = math.MaxFloat64
var CameraMinimumFocusDistance float64 = 0.0
var CameraMaximumFocusDistance float64 = math.MaxFloat64
