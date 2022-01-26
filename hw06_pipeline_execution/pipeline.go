package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	out := in
	for _, stage := range stages {
		out = doneStage(stage(out), done)
		// out = stage(out)

	}
	return out
}

func doneStage(in In, done In) Out {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			case <-done:
				return
			}
		}
	}()
	return out
}
