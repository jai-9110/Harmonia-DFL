package edge

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"

	"harmonia.com/steward/operator/util"

	"harmonia.com/steward/config" // ----為了取得gitUserToken---------------------------
	
)

var localTrainRoundCount int

func pullBaseModelAndLocalTrain(appURI string, baseModelURL string, localModelURLs []string, epochCount int) func() {
	return func() {
		zap.L().Debug("pull aggregated model")
		util.PullData(baseModelURL)

		repoMetadata, err := util.ReadMetadata(baseModelURL)
		if err != nil {
			zap.L().Fatal("Cannot read repoMetadata", zap.Error(err))
			return
		}

		metadata := map[string]string{}
		for key, val := range repoMetadata["metadata"].(map[string]interface{}) {
			metadata[key] = val.(string)
		}
		metrics := map[string]float64{}
		for key, val := range repoMetadata["metrics"].(map[string]interface{}) {
			metrics[key] = val.(float64)
		}

		for index, localModelURL := range localModelURLs { // localModelURL string 改成 localModelURLs []string
			zap.L().Debug(fmt.Sprintf("localModelURL[%d]: %s", index, localModelURL)) // 檢查目前的localModelURL
			sendLocalTrainMessage(
				appURI,
				epochCount,
				baseModel{
					gitHttpURL: baseModelURL,
					metadata:   metadata,
					metrics:    metrics,
				},
				localModelURL,
			)
		}

	}
}

var StateTransit = util.StateTransit{
	reflect.TypeOf(idleState{}): {
		reflect.TypeOf(&trainPlanAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			trainPlanAction := action.(*trainPlanAction)
			if (util.TrainPlan{}) == trainPlanAction.trainPlan {
				zap.L().Warn("train plan is empty")
			}

			return trainInitState{
					init:            false,
					pretrainedModel: false,
					trainPlan:       trainPlanAction.trainPlan,
				}, []func(){
					func() {
						// Synchonous message
						sendInitMessage(operator.GetPayload().(Payload).GrpcServerURI)
						operator.Dispatch(&initMessageResponseAction{})
					},
				}
		},
	},
	reflect.TypeOf(trainInitState{}): {
		reflect.TypeOf(&trainStartAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			trainInitState := state.(trainInitState)
			return localTrainState{
					roundCount: 0,
					trainPlan:  trainInitState.trainPlan,
				}, []func(){
					func() {
						repoMetadata, err := util.ReadMetadata(operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL)
						if err != nil {
							zap.L().Fatal("Cannot read repoMetadata", zap.Error(err))
							return
						}

						metadata := map[string]string{}
						for key, val := range repoMetadata["metadata"].(map[string]interface{}) {
							metadata[key] = val.(string)
						}
						metrics := map[string]float64{}
						for key, val := range repoMetadata["metrics"].(map[string]interface{}) {
							metrics[key] = val.(float64)
						}
						
						edgeModelRepoGitHttpURLs := operator.GetPayload().(Payload).EdgeModelRepoGitHttpURLs // EdgeModelRepoGitHttpURL string 改成 EdgeModelRepoGitHttpURLs []string
						for index, edgeModelRepoGitHttpURL := range edgeModelRepoGitHttpURLs  { 

							zap.L().Debug(fmt.Sprintf("trainInitState_edgeModelRepoGitHttpURL[%d]: %s", index, edgeModelRepoGitHttpURL))

							sendLocalTrainMessage(
								operator.GetPayload().(Payload).GrpcServerURI,
								trainInitState.trainPlan.EpochCount,
								baseModel{
									gitHttpURL: operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL,
									metadata:   metadata,
									metrics:    metrics,
								},
								edgeModelRepoGitHttpURL,
							)
						}
					},
				}
		},
		reflect.TypeOf(&initMessageResponseAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			currentState := state.(trainInitState)
			if !currentState.pretrainedModel {
				return trainInitState{
					init:            true,
					pretrainedModel: currentState.pretrainedModel,
					trainPlan:       currentState.trainPlan,
				}, nil
			} else {
				return state, []func(){
					func() {
						operator.Dispatch(&trainStartAction{})
					},
				}
			}
		},
		reflect.TypeOf(&baseModelReceivedAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			trainInitState := state.(trainInitState)
			baseModelReceivedAction := action.(*baseModelReceivedAction)

			if baseModelReceivedAction.ref != fmt.Sprintf("refs/heads/%s", util.TrainBranch(trainInitState.trainPlan.Name, trainInitState.trainPlan.CommitID)) {
				return state, nil
			}

			return state, []func(){
				func() {
					util.CheckoutPretrainedModel(
						operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL,
						trainInitState.trainPlan.Name,
						trainInitState.trainPlan.CommitID,
					)
					operator.Dispatch(&pretrainedModelReadyAction{})
				},
			}
		},
		reflect.TypeOf(&pretrainedModelReadyAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			currentState := state.(trainInitState)
			if !currentState.init {
				return trainInitState{
					init:            currentState.init,
					pretrainedModel: true,
					trainPlan:       currentState.trainPlan,
				}, nil
			} else {
				return state, []func(){
					func() {
						operator.Dispatch(&trainStartAction{})
					},
				}
			}
		},
	},
	reflect.TypeOf(localTrainState{}): {
		reflect.TypeOf(&trainFinishAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			localTrainState := state.(localTrainState)
			trainFinishAction := action.(*trainFinishAction)

			localTrainRoundCount = localTrainState.roundCount + 1 //-----------將localTrainRoundCount提出單獨計算 避免在for迴圈中重複計算------------------------------------
			pushModel := func() {
				edgeModelRepoGitURLs := operator.GetPayload().(Payload).EdgeModelRepoGitHttpURLs // EdgeModelRepoGitHttpURL string 改成 EdgeModelRepoGitHttpURLs [] string

				for index, edgeModelRepoGitURL := range edgeModelRepoGitURLs {
					
					util.GitSetup(util.GitUser{ // 雖然有在steward.go中run過此func，當時執行完後只會連上最後一個token的gitea
						"Harmonia Operator",    // 導致在push local model回gitea時只會連到之前最後一個token的gitea
						"operator@harmonia",    // 因此在push local model回去前，再次根據gitUserToken連上gitea
						config.Config.GitUserTokens[index].Token,  // 確保能推到正確的gitea
					})

					zap.L().Debug(fmt.Sprintf("pushModel-edgeModelRepoGitURL[%d]: %s", index, edgeModelRepoGitURL))

					err := util.WriteMetadata(edgeModelRepoGitURL, map[string]interface{}{
						"trainPlanID":     localTrainState.trainPlan.CommitID,
						"datasetSize":     trainFinishAction.datasetSize,
						"metadata":        trainFinishAction.metadata,
						"metrics":         trainFinishAction.metrics,
						"roundNumber":     localTrainRoundCount,
						"pretrainedModel": localTrainState.trainPlan.PretrainedModelCommitID,
						// TODO
						// "baseModel"
					})

					if err != nil {
						zap.L().Error("push edge model error1", zap.Error(err))
						return
					}

					err = util.PushUpdates(edgeModelRepoGitURL, "")

					if err != nil {
						zap.L().Error("push edge model error2", zap.Error(err))
					}

				}
			}

			if (util.TrainPlan{}) == localTrainState.trainPlan {
				zap.L().Warn("train plan is empty")
			}
			return aggregateState{
				roundCount: localTrainRoundCount,
				trainPlan:  localTrainState.trainPlan,
			}, []func(){
				pushModel,
			}

		},
		reflect.TypeOf(&baseModelReceivedAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			localTrainState := state.(localTrainState)
			baseModelReceivedAction := action.(*baseModelReceivedAction)

			if baseModelReceivedAction.ref != fmt.Sprintf("refs/heads/%s", util.TrainBranch(localTrainState.trainPlan.Name, localTrainState.trainPlan.CommitID)) {
				return state, nil
			}

			return localTrainInterruptedState{
					roundCount: localTrainState.roundCount,
					trainPlan:  localTrainState.trainPlan,
				}, []func(){
					func() {
						sendTrainInterruptMessage(
							operator.GetPayload().(Payload).GrpcServerURI,
						)
						operator.Dispatch(&trainCleanupAction{
							roundCount: localTrainState.roundCount,
						})
					},
				}
		},
	},
	reflect.TypeOf(localTrainInterruptedState{}): {
		reflect.TypeOf(&trainCleanupAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			currentState := state.(localTrainInterruptedState)
			trainCleanupAction := action.(*trainCleanupAction)

			if currentState.roundCount != trainCleanupAction.roundCount {
				return currentState, nil
			}

			return localTrainState{
					roundCount: currentState.roundCount + 1,
					trainPlan:  currentState.trainPlan,
				}, []func(){
					pullBaseModelAndLocalTrain(
						operator.GetPayload().(Payload).GrpcServerURI,
						operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL,
						operator.GetPayload().(Payload).EdgeModelRepoGitHttpURLs,
						currentState.trainPlan.EpochCount,
					),
				}

		},
		reflect.TypeOf(&baseModelReceivedAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			currentState := state.(localTrainInterruptedState)
			baseModelReceivedAction := action.(*baseModelReceivedAction)

			if baseModelReceivedAction.ref != fmt.Sprintf("refs/heads/%s", util.TrainBranch(currentState.trainPlan.Name, currentState.trainPlan.CommitID)) {
				return state, nil
			}

			return localTrainInterruptedState{
					roundCount: currentState.roundCount + 1,
					trainPlan:  currentState.trainPlan,
				}, []func(){
					func() {
						sendTrainInterruptMessage(
							operator.GetPayload().(Payload).GrpcServerURI,
						)
						operator.Dispatch(&trainCleanupAction{
							roundCount: currentState.roundCount + 1,
						})
					},
				}
		},
	},
	reflect.TypeOf(aggregateState{}): {
		reflect.TypeOf(&baseModelReceivedAction{}): func(state util.State, action util.Action, operator util.AbstractOperator) (util.State, []func()) {
			aggregateState := state.(aggregateState)
			baseModelReceivedAction := action.(*baseModelReceivedAction)

			if (util.TrainPlan{}) == aggregateState.trainPlan {
				zap.L().Warn("train plan is empty")
			}
			if baseModelReceivedAction.ref != fmt.Sprintf("refs/heads/%s", util.TrainBranch(aggregateState.trainPlan.Name, aggregateState.trainPlan.CommitID)) {
				return state, nil
			}
			if aggregateState.roundCount == aggregateState.trainPlan.RoundCount {
				return idleState{}, []func(){
					func() {
						util.PullData(operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL)
						sendTrainFinishMessage(
							operator.GetPayload().(Payload).GrpcServerURI,
						)
						operator.TrainFinish()
					},
				}
			} else {
				return localTrainState{
						roundCount: aggregateState.roundCount,
						trainPlan:  aggregateState.trainPlan,
					}, []func(){
						pullBaseModelAndLocalTrain(
							operator.GetPayload().(Payload).GrpcServerURI,
							operator.GetPayload().(Payload).AggregatedModelRepoGitHttpURL,
							operator.GetPayload().(Payload).EdgeModelRepoGitHttpURLs,
							aggregateState.trainPlan.EpochCount,
						),
					}
			}
		},
	},
}

var InitState = idleState{}
