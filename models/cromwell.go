package models

import "fmt"

type IdAndStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type WorkflowQueryParameter struct {
	Submission string `json:"submission,omitempty"`
	Start      string `json:"start,omitempty"`
	End        string `json:"end,omitempty"`
	Status     string `json:"status,omitempty"`
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
}

type WorkflowQueryResult struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Submission string `json:"submission"`
	Start      string `json:"start"`
	End        string `json:"end"`
}

func (r *WorkflowQueryResult) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		r.ID, r.Submission, r.Start, r.End, r.Status, r.Name)
}

type WorkflowQueryResponse struct {
	Results           []*WorkflowQueryResult `json:"results"`
	TotalResultsCount int                    `json:"totalResultsCount"`
}

type DescriptorTypeAndVersion struct {
	DescriptorType        string `json:"descriptorType"`
	DescriptorTypeVersion string `json:"descriptorTypeVersion"`
}

type ValueType struct {
	TypeName         string    `json:"typeName"`
	OptionalType     string    `json:"optionalType"`
	ArrayType        string    `json:"arrayType"`
	MapType          string    `json:"mapType"`
	TupleTypes       []*string `json:"tupleTypes"`
	ObjectFieldTypes []struct {
		FieldName string `json:"fieldName"`
		FieldType string `json:"fieldType"`
	} `json:"objectFieldTypes"`
}

type ToolInputParameter struct {
	Name            string    `json:"name"`
	ValueType       ValueType `json:"valueType"`
	Optional        bool      `json:"optional"`
	Default         string    `json:"default"`
	TypeDisplayName string    `json:"typeDisplayName"`
}

type ToolOutputParameter struct {
	Name            string    `json:"name"`
	ValueType       ValueType `json:"valueType"`
	TypeDisplayName string    `json:"typeDisplayName"`
}

type WorkflowDescription struct {
	Valid                   bool                      `json:"valid"`
	Errors                  []string                  `json:"errors"`
	ValidWorkflow           bool                      `json:"validWorkflow"`
	Name                    string                    `json:"name"`
	Inputs                  []*ToolInputParameter     `json:"inputs"`
	Outputs                 []*ToolOutputParameter    `json:"outputs"`
	SubmittedDescriptorType *DescriptorTypeAndVersion `json:"submittedDescriptorType"`
	IsRunnableWorkflow      bool                      `json:"isRunnableWorkflow"`
}
