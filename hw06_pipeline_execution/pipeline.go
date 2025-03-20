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
		for {
			select {
			case <-done:
				close(chNext)
				for data := range in {
					_ = data
				}
				return
			case data, ok := <-in:
				if !ok {
					close(chNext)
					return
				}
				select {
				case <-done:
					close(chNext)
					for data := range in {
						_ = data
					}
					return
				case chNext <- data:
				}
			}
		}
	}()
	return chNext
}
