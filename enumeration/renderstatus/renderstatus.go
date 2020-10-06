package renderstatus

// RenderStatus represents the status of a scene render
type RenderStatus string

// Created - Scene is created and registered, but no render has even been started
var Created RenderStatus = "CREATED"

// Pending - Render has been requested, waiting to begin
var Pending RenderStatus = "PENDING"

// Starting - Render has been started, initialization is in progress (loading from database, assembling primitives, etc)
var Starting RenderStatus = "STARTING"

// Running - Render is actively running
var Running RenderStatus = "RUNNING"

// Stopping - Render has been requested to stop and is attemping to stop
var Stopping RenderStatus = "STOPPING"

// Stopped - Render has been manually stopped
var Stopped RenderStatus = "STOPPED"

// Completed - Render has completed
var Completed RenderStatus = "COMPLETED"

// Error - Render has been cancelled due to an unexpected error
var Error RenderStatus = "ERROR"
