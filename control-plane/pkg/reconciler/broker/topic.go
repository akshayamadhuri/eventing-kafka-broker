/*
 * Copyright 2020 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package broker

import (
	"fmt"

	"github.com/Shopify/sarama"
	eventing "knative.dev/eventing/pkg/apis/eventing/v1beta1"
)

func (r *Reconciler) CreateTopic(broker *eventing.Broker) (string, error) {
	topic := Topic(broker)
	topicDetail := r.topicDetailFromBrokerConfig(broker)

	r.KafkaClusterAdminLock.Lock()
	defer r.KafkaClusterAdminLock.Unlock()

	createTopicError := r.KafkaClusterAdmin.CreateTopic(topic, topicDetail, true)
	if err, ok := createTopicError.(*sarama.TopicError); ok && err.Err == sarama.ErrTopicAlreadyExists {
		return topic, nil
	}

	return topic, createTopicError
}

func (r *Reconciler) deleteTopic(broker *eventing.Broker) (string, error) {

	r.KafkaClusterAdminLock.RLock()
	defer r.KafkaClusterAdminLock.RUnlock()

	topic := Topic(broker)
	return topic, r.KafkaClusterAdmin.DeleteTopic(topic)
}

func Topic(broker *eventing.Broker) string {
	return fmt.Sprintf("%s%s-%s", TopicPrefix, broker.Namespace, broker.Name)
}
