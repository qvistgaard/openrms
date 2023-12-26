package fuel

type average struct {
	last    float32
	usage   []float32
	average float32
}

func (a average) reportUsage(usage float32) average {
	av := average{}

	if usage < a.last {
		av.usage = []float32{a.average}
	} else {
		var u float32
		if a.last == 0 {
			u = usage
		} else {
			u = usage - a.last
		}
		av.usage = append(a.usage, u)
	}

	if len(av.usage) > 10 {
		av.usage = av.usage[1:]
	}

	av.average = av.averageUsage()
	av.last = usage
	return av
}

func (a average) averageUsage() float32 {
	sum := float32(0)
	for _, value := range a.usage {
		sum += value
	}
	return sum / float32(len(a.usage))
}
