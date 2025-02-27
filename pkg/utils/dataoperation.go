/*
  Copyright 2023 The Fluid Authors.

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

package utils

import (
	"context"
	"fmt"

	datav1alpha1 "github.com/fluid-cloudnative/fluid/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListDataOperationJobByCronjob gets the DataOperation(i.e. DataMigrate, DataLoad) job by cronjob given its name and namespace
func ListDataOperationJobByCronjob(c client.Client, cronjobNamespacedName types.NamespacedName) ([]batchv1.Job, error) {
	jobLabelSelector, err := labels.Parse(fmt.Sprintf("cronjob=%s", cronjobNamespacedName.Name))
	if err != nil {
		return nil, err
	}
	var jobList batchv1.JobList
	if err := c.List(context.TODO(), &jobList, &client.ListOptions{
		LabelSelector: jobLabelSelector,
		Namespace:     cronjobNamespacedName.Namespace,
	}); err != nil {
		return nil, err
	}
	return jobList.Items, nil
}

func GetOperationStatus(obj client.Object) (*datav1alpha1.OperationStatus, error) {
	if obj == nil {
		return nil, nil
	}

	if dataLoad, ok := obj.(*datav1alpha1.DataLoad); ok {
		return dataLoad.Status.DeepCopy(), nil
	} else if dataMigrate, ok := obj.(*datav1alpha1.DataMigrate); ok {
		return dataMigrate.Status.DeepCopy(), nil
	} else if dataBackup, ok := obj.(*datav1alpha1.DataBackup); ok {
		return dataBackup.Status.DeepCopy(), nil
	} else if dataProcess, ok := obj.(*datav1alpha1.DataProcess); ok {
		return dataProcess.Status.DeepCopy(), nil
	}

	return nil, fmt.Errorf("obj is not of any data operation type")
}

func GetPrecedingOperationStatus(client client.Client, opRef *datav1alpha1.OperationRef) (*datav1alpha1.OperationStatus, error) {
	if opRef == nil {
		return nil, nil
	}

	switch opRef.OperationKind {
	case datav1alpha1.DataBackupType:
		object, err := GetDataBackup(client, opRef.Name, opRef.Namespace)
		if err != nil {
			return nil, err
		}
		return &object.Status, nil
	case datav1alpha1.DataLoadType:
		object, err := GetDataLoad(client, opRef.Name, opRef.Namespace)
		if err != nil {
			return nil, err
		}
		return &object.Status, nil
	case datav1alpha1.DataMigrateType:
		object, err := GetDataMigrate(client, opRef.Name, opRef.Namespace)
		if err != nil {
			return nil, err
		}
		return &object.Status, nil
	case datav1alpha1.DataProcessType:
		object, err := GetDataProcess(client, opRef.Name, opRef.Namespace)
		if err != nil {
			return nil, err
		}
		return &object.Status, nil
	default:
		return nil, fmt.Errorf("unknown data operation kind")
	}
}

func HasPrecedingOperation(obj client.Object) (has bool, err error) {
	if obj == nil {
		return false, nil
	}

	if dataLoad, ok := obj.(*datav1alpha1.DataLoad); ok {
		return dataLoad.Spec.RunAfter != nil, nil
	} else if dataMigrate, ok := obj.(*datav1alpha1.DataMigrate); ok {
		return dataMigrate.Spec.RunAfter != nil, nil
	} else if dataBackup, ok := obj.(*datav1alpha1.DataBackup); ok {
		return dataBackup.Spec.RunAfter != nil, nil
	} else if dataProcess, ok := obj.(*datav1alpha1.DataProcess); ok {
		return dataProcess.Spec.RunAfter != nil, nil
	}

	return false, fmt.Errorf("obj is not of any data operation type")
}
