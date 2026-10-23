package main

func (p *Plugin) runJob() {
	// Include job logic here
	p.client.Log.Debug("Background job is running")
}
