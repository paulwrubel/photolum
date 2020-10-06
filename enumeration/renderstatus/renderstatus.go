package renderstatus

// RenderStatus represents the status of a scene render
type RenderStatus string

// Created - Scene is created and registered, but no render has even been started
var Created RenderStatus = "Created"

// Pending - Render has been requested, waiting to begin
var Pending RenderStatus = "Pending"

// Starting - Render has been started, initialization is in progress (loading from database, assembling primitives, etc)
var Starting RenderStatus = "Starting"

// Running - Render is actively running
var Running RenderStatus = "Running"

// Stopping - Render has been requested to stop and is attemping to stop
var Stopping RenderStatus = "Stopping"

// Stopped - Render has been manually stopped
var Stopped RenderStatus = "Stopped"

// Completed - Render has completed
var Completed RenderStatus = "Completed"

// Error - Render has been cancelled due to an unexpected error
var Error RenderStatus = "Error"
