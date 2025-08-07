package source

type TopicPuller struct {
	Topics []string
	Source Source
}

func NewTopicPuller(topic []string, source Source) *TopicPuller {
	return &TopicPuller{Topics: topic, Source: source}
}
