package payment_workflow

import (
	"context"
	"time"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/apperror"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	authPaymentWorkflowName    = "auth_payment_workflow"
	capturePaymentWorkflowName = "capture_payment_workflow"
	requestActivity            = `request_activity`
	checkTransactionActivity   = `check_transaction_activity`
	publishActivity            = `publish_activity`
	captureActivity            = `capture_activity`
	authActivity               = `auth_activity`

)

var nonRetryableErrorTypes = []string{}

type (
	paymentWorkflow struct {
		workflow temporal_client.WorkflowExecution
		temporal temporal_client.Temporal
	}

	PaymentWorkflow interface {
		AuthPaymentWorkflow(ctx context.Context, request *requests.AuthPaymentRequest) (*responses.AuthPaymentResponse, *apperror.CustomError)
		CapturePaymentWorkflow(ctx context.Context, request *requests.CapturePaymentRequest) (*responses.CapturePaymentResponse, *apperror.CustomError)
	}
)

func (p *paymentWorkflow) registerWorkflowAndActivities() {
	activityOption := &workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        5,
			InitialInterval:        2 * time.Second,
			BackoffCoefficient:     2,
			MaximumInterval:        30 * time.Second,
			NonRetryableErrorTypes: nonRetryableErrorTypes,
		},
	}
	p.workflow.AddTransitionActivityWithOptions(requestActivity, "", p.request, activityOption, checkTransactionActivity)
	p.workflow.AddTransitionActivityWithOptions(checkTransactionActivity, "", p.checkTransaction, activityOption, publishActivity)
	p.workflow.AddTransitionActivityWithOptions(publishActivity, "", p.publish, activityOption, captureActivity)

	p.workflow.RegisterWorkflow(authPaymentWorkflowName, p.auth)
	p.workflow.RegisterWorkflow(capturePaymentWorkflowName, p.capture)
}

func (p *paymentWorkflow) AuthPaymentWorkflow(ctx context.Context, request *requests.AuthPaymentRequest) (*responses.AuthPaymentResponse, *apperror.CustomError) {
	return nil, nil
}

func (p *paymentWorkflow) CapturePaymentWorkflow(ctx context.Context, request *requests.CapturePaymentRequest) (*responses.CapturePaymentResponse, *apperror.CustomError) {
	return nil, nil
}

func NewPaymentWorkflow(
	workflow temporal_client.WorkflowExecution,
	temporal temporal_client.Temporal,
) *paymentWorkflow {
	p := &paymentWorkflow{
		workflow: workflow,
		temporal: temporal,
	}
	p.registerWorkflowAndActivities()
	return p
}
