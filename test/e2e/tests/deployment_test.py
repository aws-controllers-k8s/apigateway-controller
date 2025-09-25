# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the Resource resource
"""

import logging
import time
from functools import partial
from typing import Dict, Tuple

import boto3
import pytest
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_apigateway_resource
from e2e.common.waiter import wait_until_deleted, safe_get
from e2e.replacement_values import REPLACEMENT_VALUES
from .rest_api_test import simple_rest_api
from .resource_test import simple_resource
from .integration_test import simple_integration
from acktest.k8s import condition

RESOURCE_RESOURCE_PLURAL = "deployments"
MODIFY_WAIT_AFTER_SECONDS = 60
MAX_RETRIES = 3
WAIT_TIME = 30


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture(scope='module')
def simple_deployment(apigateway_client, simple_integration):
    deployment_name = random_suffix_name('simple-deployment', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['DEPLOYMENT_NAME'] = deployment_name
    (_, _, _, rest_api_cr) = simple_integration
    replacements['REST_API_REF_NAME'] = rest_api_cr['spec']['name']

    deployment_data = load_apigateway_resource(
        'deployment_simple',
        additional_replacements=replacements,
    )
    logging.debug(deployment_name)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_RESOURCE_PLURAL,
        deployment_name, namespace='default',
    )

    k8s.create_custom_resource(ref, deployment_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=30)

    assert cr is not None
    assert cr['status']['id'] is not None
    assert k8s.get_resource_exists(ref)
    resource_query = {
        'restApiId': rest_api_cr['status']['id']
    }
    k8s.wait_on_condition(
        ref,
        condition.CONDITION_TYPE_READY,
        "True",
        wait_periods=60,
    )
    cr = k8s.get_resource(ref)
    resource_query['deploymentId'] = cr['status']['id']
    yield ref, cr, rest_api_cr, resource_query

    _, deleted = k8s.delete_custom_resource(ref, 10, 60)
    assert deleted


@service_marker
@pytest.mark.canary
class TestDeployment:
    def test_create_update_deployment(self, simple_deployment, apigateway_client):
        (ref, _, _, resource_query) = simple_deployment
        assert safe_get(partial(apigateway_client.get_deployment,
                        **resource_query)) is not None

        updates = {
            'spec': {
                'description': 'updated description',
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        aws_resource = apigateway_client.get_deployment(**resource_query)
        assert aws_resource['description'] == updates['spec']['description']
