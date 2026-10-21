package job

import (
	"testing"

	upgradeapi "github.com/rancher/system-upgrade-controller/pkg/apis/upgrade.cattle.io"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFailedMessage(t *testing.T) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apply-plan1-on-node1-with-hash",
			Namespace: "system-upgrade",
			Labels: map[string]string{
				upgradeapi.LabelNode: "node1",
			},
		},
	}

	if message, failed := FailedMessage(job); failed || message != "" {
		t.Errorf("expected Job without conditions to not be failed, got %q, %v", message, failed)
	}

	job.Status.Conditions = []batchv1.JobCondition{{
		Type:    batchv1.JobFailed,
		Status:  corev1.ConditionTrue,
		Reason:  "BackoffLimitExceeded",
		Message: "Job has reached the specified backoff limit",
	}}
	expected := "Job system-upgrade/apply-plan1-on-node1-with-hash failed on Node node1: BackoffLimitExceeded: Job has reached the specified backoff limit"
	if message, failed := FailedMessage(job); !failed || message != expected {
		t.Errorf("expected failed Job message %q, got %q, %v", expected, message, failed)
	}
}
