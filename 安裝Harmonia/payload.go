package edge

type Payload struct {
	GrpcServerURI string
	TrainPlanRepoGitHttpURL string
	AggregatedModelRepoGitHttpURL string
	EdgeModelRepoGitHttpURLs []string 
}
