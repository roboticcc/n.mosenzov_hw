package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(work(in, done))
	}

	return in
}

func work(in In, done In) In {
	chNext := make(Bi)

	go func() {
		defer close(chNext)
		for {
			select {
			case <-done:
				<-in
				return
			case data, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					<-in
					return
				case chNext <- data:
				}
			}
		}
	}()

	return chNext
}
