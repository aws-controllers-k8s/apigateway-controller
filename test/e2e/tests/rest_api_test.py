# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the RestAPI resource
"""

import boto3
import logging
import time

import pytest
from functools import partial

from typing import Dict, Tuple
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest.k8s import condition
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_apigateway_resource
from e2e.common.waiter import wait_until_deleted, safe_get
from e2e.replacement_values import REPLACEMENT_VALUES


REST_API_RESOURCE_PLURAL = "restapis"
MODIFY_WAIT_AFTER_SECONDS = 30
MAX_WAIT_FOR_SYNCED_MINUTES = 1


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture
def simple_rest_api(apigateway_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    rest_api_name = random_suffix_name('simple-rest-api', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['REST_API_NAME'] = rest_api_name

    resource_data = load_apigateway_resource(
        'rest_api_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, REST_API_RESOURCE_PLURAL,
        rest_api_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=15)

    assert cr is not None
    assert cr['status']['id'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(apigateway_client.get_rest_api, restApiId=cr['status']['id']))


@service_marker
@pytest.mark.canary
class TestRestAPI:
    def test_create_update_rest_api(self, simple_rest_api, apigateway_client):
        (ref, cr) = simple_rest_api
        rest_api_id = cr['status']['id']
        assert safe_get(partial(apigateway_client.get_rest_api, restApiId=rest_api_id)) is not None

        updates = {
            'spec': {
                'apiKeySource': 'INVALID',
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_TERMINAL,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )

        updates = {
            'spec': {
                'apiKeySource': 'AUTHORIZER',
                'minimumCompressionSize': 8192,
                'binaryMediaTypes': ['application/octet-stream', 'application/vnd.apache.thrift.binary'],
                'description': 'Updated description',
                'tags': {
                    'k1': 'v10',
                    'k2': 'v20',
                    'k3': 'v3',
                    'k4': 'v4',
                }
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        assert k8s.get_resource_condition(ref, condition.CONDITION_TYPE_TERMINAL) is None

        aws_rest_api = apigateway_client.get_rest_api(restApiId=rest_api_id)
        expected_tags = updates['spec'].pop('tags')
        updated_fields = {field: aws_rest_api[field] for field in updates['spec'].keys()}
        assert updated_fields == updates['spec']
        tags.assert_equal_without_ack_tags(
            expected=expected_tags,
            actual=aws_rest_api['tags'],
        )
