/*
Copyright 2024 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName="nb"
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.status.type`,description="The type of backup, such as full, incr"
// +kubebuilder:printcolumn:name="BACKUP",type=string,JSONPath=`.status.backupName`,description="The name of the backup generated by nebula"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`,description="The current status of the backup"
// +kubebuilder:printcolumn:name="Started",type=date,JSONPath=`.status.timeStarted`,description="The time at which the backup was started"
// +kubebuilder:printcolumn:name="Completed",type=date,JSONPath=`.status.timeCompleted`,description="The time at which the backup was completed"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

type NebulaBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NebulaBackupList contains a list of NebulaBackup.
type NebulaBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NebulaBackup `json:"items"`
}

// BackupType represents the backup type.
// +k8s:openapi-gen=true
type BackupType string

const (
	// BackupTypeFull represents the full backup of nebula cluster.
	BackupTypeFull BackupType = "full"
	// BackupTypeIncr represents the incremental backup of nebula cluster.
	BackupTypeIncr BackupType = "incr"
)

func (t BackupType) String() string {
	return string(t)
}

// BackupConditionType represents a valid condition of a Backup.
type BackupConditionType string

const (
	// BackupRunning means the backup is currently being executed.
	BackupRunning BackupConditionType = "Running"
	// BackupComplete means the backup has successfully executed and the
	// backup data has been loaded into nebula cluster.
	BackupComplete BackupConditionType = "Complete"
	// BackupClean means the clean job has been clean the backup data.
	BackupClean BackupConditionType = "Clean"
	// BackupFailed means the backup has failed.
	BackupFailed BackupConditionType = "Failed"
	// BackupInvalid means invalid backup CR.
	BackupInvalid BackupConditionType = "Invalid"
)

// BackupCondition describes the observed state of a Backup at a certain point.
type BackupCondition struct {
	// Type of the condition.
	Type BackupConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human-readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// BackupSpec contains the specification for a backup of a nebula cluster backup.
type BackupSpec struct {
	// Container image.
	// +optional
	Image string `json:"image,omitempty"`

	// Version tag for container image.
	// +optional
	Version string `json:"version,omitempty"`

	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// +kubebuilder:default=Always
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// +optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// +optional
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`

	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// CleanBackupData denotes whether to clean backup data when the object is deleted from the cluster,
	// if not set, the backup data will be retained
	// +optional
	CleanBackupData *bool `json:"cleanBackupData,omitempty"`

	// The job that status is failed and completed will be removed automatically.
	// +optional
	AutoRemoveFinished *bool `json:"autoRemoveFinished,omitempty"`

	Config *BackupConfig `json:"config,omitempty"`
}

type BackupConfig struct {
	NamespacedObjectReference `json:",inline"`
	// The name of the base backup and only used for incremental backup.
	BaseBackupName *string `json:"baseBackupName,omitempty"`
	// StorageProvider configures where and how backups should be stored.
	StorageProvider `json:",inline"`
}

// BackupStatus represents the current status of a nebula cluster backup.
type BackupStatus struct {
	// Type is the backup type for nebula cluster.
	Type BackupType `json:"type,omitempty"`
	// BackupName is the name of the backup generated by nebula.
	BackupName string `json:"backupName,omitempty"`
	// TimeStarted is the time at which the backup was started.
	// +nullable
	TimeStarted *metav1.Time `json:"timeStarted,omitempty"`
	// TimeCompleted is the time at which the backup was completed.
	// +nullable
	TimeCompleted *metav1.Time `json:"timeCompleted,omitempty"`
	// Phase is a user readable state inferred from the underlying Backup conditions.
	Phase BackupConditionType `json:"phase,omitempty"`
	// +nullable
	Conditions []BackupCondition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&NebulaBackup{}, &NebulaBackupList{})
}
