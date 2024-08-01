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

"""Integration tests for the Stage resource
"""

import logging
import time
from typing import Dict, Tuple
from functools import partial

import boto3
import pytest
from acktest import tags
from acktest.k8s import resource as k8s
from acktest.k8s import condition
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_apigateway_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common.waiter import wait_until_deleted, safe_get
from .rest_api_test import simple_rest_api
from .resource_test import simple_resource
from .integration_test import simple_integration


STAGE_RESOURCE_PLURAL = 'stages'
MODIFY_WAIT_AFTER_SECONDS = 60
MAX_WAIT_FOR_SYNCED_MINUTES = 1


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture
def simple_stage(simple_integration, apigateway_client) -> Tuple[k8s.CustomResourceReference, Dict, str]:
    stage_name = random_suffix_name('simple-stage', 32)
    (ref, cr, resource_query) = simple_integration
    rest_api_id = resource_query['restApiId']
    deployment_res = apigateway_client.create_deployment(restApiId=rest_api_id, description=stage_name)

    resource_data = load_apigateway_resource(
        'stage_simple',
        additional_replacements={
            **REPLACEMENT_VALUES,
            **{
                'STAGE_NAME': stage_name,
                'REST_API_REF_NAME': cr['spec']['restAPIRef']['from']['name'],
                'DEPLOYMENT_ID': deployment_res['id'],
            },
        },
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, STAGE_RESOURCE_PLURAL,
        stage_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=15)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr, rest_api_id

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(apigateway_client.get_stage, restApiId=rest_api_id, stageName=stage_name))
    apigateway_client.delete_deployment(restApiId=rest_api_id, deploymentId=deployment_res['id'])


@service_marker
@pytest.mark.canary
class TestStage:
    def test_create_update_stage(self, simple_stage, apigateway_client):
        (ref, cr, rest_api_id) = simple_stage

        get_stage = partial(apigateway_client.get_stage, restApiId=rest_api_id, stageName=cr['spec']['stageName'])
        assert safe_get(get_stage) is not None

        updates = {
            'canarySettings': {
                'stageVariableOverrides': {
                    'v1': 'k1',
                    'v2': 'k2',
                },
                'percentTraffic': 10.1,
            },
            'description': 'Updated description',
            'variables': {
                'v1': 'v1',
                'v2': 'v2',
                'var1': 'val1',
                'var2': 'val2',
            },
            'tags': {
                'k1': 'v1',
                'k2': 'v10',
                'k3': 'v3',
                'k4': 'v4',
            }
        }
        k8s.patch_custom_resource(ref, {'spec': updates})
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        assert k8s.get_resource_condition(ref, condition.CONDITION_TYPE_TERMINAL) is None

        aws_resource = get_stage()
        expected_tags = updates.pop('tags')
        updated_canary_fields = {field: aws_resource['canarySettings'][field] for field in updates['canarySettings'].keys()}
        updated_fields = {
            **{field: aws_resource[field] for field in updates.keys() if field != 'canarySettings'},
            **{'canarySettings': updated_canary_fields},
        }
        assert updated_fields == updates
        tags.assert_equal_without_ack_tags(
            expected=expected_tags,
            actual=aws_resource['tags'],
        )
