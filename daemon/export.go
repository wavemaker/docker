package daemon

import (
	"io"

	"github.com/docker/docker/engine"
	"github.com/docker/docker/archive"
	"github.com/docker/docker/pkg/log"
)

func (daemon *Daemon) ContainerExport(job *engine.Job) engine.Status {
	if len(job.Args) != 1 {
		return job.Errorf("Usage: %s container_id", job.Name)
	}
	name := job.Args[0]
	readwrite := job.GetenvBool("readwrite")
	log.Debugf("the boolean readwrite is %t ", readwrite)
	if container := daemon.Get(name); container != nil {
        var (
            data archive.Archive
            err error
        )

	    if readwrite {
	        log.Debugf("Exporting only the rw directory ")
	        data, err = container.ExportRw()
	    } else {
	        log.Debugf("Exporting Entire rootfs/basefs of container ")
		    data, err = container.Export()
		}
		if err != nil {
			return job.Errorf("%s: %s", name, err)
		}
		defer data.Close()

		// Stream the entire contents of the container (basically a volatile snapshot)
		if _, err := io.Copy(job.Stdout, data); err != nil {
			return job.Errorf("%s: %s", name, err)
		}
		// FIXME: factor job-specific LogEvent to engine.Job.Run()
		container.LogEvent("export")
		return engine.StatusOK
	}
	return job.Errorf("No such container: %s", name)
}
