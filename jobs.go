package gocd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
)

// ScheduledJob instance
type ScheduledJob struct {
	Name         string    `xml:"name,attr"`
	JobID        string    `xml:"id,attr"`
	BuildLocator string    `xml:"buildLocator"`
	Link         LinkInXML `xml:"link"`
}

// JobURL - Full URL location of the scheduled job
func (job *ScheduledJob) JobURL() string {
	return job.Link.Href
}

// LinkInXML - <link rel="..." href="..."> tag
type LinkInXML struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// GetScheduledJobs - Lists all the current job instances which are scheduled but not yet assigned to any agent.
func (c *Client) GetScheduledJobs() ([]*ScheduledJob, error) {
	var errors *multierror.Error

	type ScheduledJobsResponse struct {
		XMLName xml.Name        `xml:"scheduledJobs"`
		Jobs    []*ScheduledJob `xml:"job"`
	}

	var jobs ScheduledJobsResponse
	_, body, errs := c.Request.
		Get(c.resolve("/go/api/jobs/scheduled.xml")).
		End()
	multierror.Append(errors, errs...)
	xmlErr := xml.Unmarshal([]byte(body), &jobs)
	multierror.Append(errors, xmlErr)

	return jobs.Jobs, errors.ErrorOrNil()
}

// JobHistory - Represents the
type JobHistory struct {
	AgentUUID           string   `json:"agent_uuid"`
	Name                string   `json:"name"`
	JobStateTransitions []string `json:"job_state_transitions"`
	ScheduledDate       int      `json:"scheduled_date"`
	OriginalJobID       string   `json:"original_job_id"`
	PipelineCounter     int      `json:"pipeline_counter"`
	PipelineName        string   `json:"pipeline_name"`
	Result              string   `json:"result"`
	State               string   `json:"state"`
	ID                  int      `json:"id"`
	StageCounter        string   `json:"stage_counter"`
	StageName           string   `json:"stage_name"`
	ReRun               bool     `json:"rerun"`
}

// GetJobHistory - The job history allows users to list job instances of specified job. Supports pagination using offset which tells the API how many instances to skip.
func (c *Client) GetJobHistory(pipeline, stage, job string, offset int) ([]*JobHistory, error) {
	var errors *multierror.Error
	_, body, errs := c.Request.
		Get(c.resolve(fmt.Sprintf("/go/api/jobs/%s/%s/%s/history/%d", pipeline, stage, job, offset))).
		Set("Accept", "application/vnd.go.cd.v2+json").
		End()
	multierror.Append(errors, errs...)

	type JobHistoryResponse struct {
		Jobs []*JobHistory `json:"jobs"`
	}
	var jobs *JobHistoryResponse
	jsonErr := json.Unmarshal([]byte(body), &jobs)
	multierror.Append(errors, jsonErr)
	return jobs.Jobs, errors.ErrorOrNil()
}