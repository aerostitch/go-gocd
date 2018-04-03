package gocd

import (
	"fmt"
)

// PipelineInstance represents a pipeline instance (for a given run)
type PipelineInstance struct {
	ID                  int                `json:"id"`
	Name                string             `json:"name"`
	Label               string             `json:"label"`
	NaturalOrder        float32            `json:"natural_order"`
	CanRun              bool               `json:"can_run"`
	Comment             string             `json:"comment"`
	Counter             int                `json:"counter"`
	PreparingToSchedule bool               `json:"preparing_to_schedule"`
	Stages              []StageRun         `json:"stages"`
	BuildCause          PipelineBuildCause `json:"build_cause"`
}

// PipelineBuildCause represent what triggered the build of the pipeline
type PipelineBuildCause struct {
	Approver          string             `json:"approver"`
	MaterialRevisions []MaterialRevision `json:"material_revisions"`
	TriggerForced     bool               `json:"trigger_forced"`
	TriggerMessage    string             `json:"trigger_message"`
}

// PipelineHistoryPage represents a page of the history of run of a pipeline
type PipelineHistoryPage struct {
	Pipelines  []PipelineInstance `json:"pipelines"`
	Pagination Pagination         `json:"pagination"`
}

// GetPipelineInstance returns the pipeline instance corresponding to the given
// pipeline name and counter
func (c *DefaultClient) GetPipelineInstance(name string, counter int) (*PipelineInstance, error) {
	res := new(PipelineInstance)
	err := c.getJSON(fmt.Sprintf("/go/api/pipelines/%s/instance/%d", name, counter), nil, res)
	return res, err
}

// GetPipelineHistoryPage allows users to list pipeline instances. Supports
// pagination using offset which tells the API how many instances to skip.
// Note that te history is listed in reverse chronological order meaning the
// setting an offset to 1 will skip the last run of the pipeline and will give
// you a page of pipeline runs history which is 10 by default.
func (c *DefaultClient) GetPipelineHistoryPage(name string, offset int) (*PipelineHistoryPage, error) {
	res := new(PipelineHistoryPage)
	err := c.getJSON(fmt.Sprintf("/go/api/pipelines/%s/history/%d", name, offset), nil, res)
	return res, err
}
