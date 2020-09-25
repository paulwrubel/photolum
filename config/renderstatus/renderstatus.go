package renderstatus

// RenderStatus represents the status of a scene render
type RenderStatus string

// Created - Scene is created and registered, but no render has even been started
var Created RenderStatus = "Created"

// Pending - Render has been requested, waiting to begin
var Pending RenderStatus = "Pending"

// Running - Render is actively running
var Running RenderStatus = "Running"

// Stopped - Render has been manually stopped
var Stopped RenderStatus = "Stopped"

// Completed - Render has completed
var Completed RenderStatus = "Completed"

// Error - Render has been cancelled due to an unexpected error
var Error RenderStatus = "Error"
