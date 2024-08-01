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

"""Integration tests for the Integration resource
"""

import logging
import time
from typing import Dict, Tuple
from functools import partial

import boto3
import pytest
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_apigateway_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common.waiter import wait_until_deleted, safe_get
from .rest_api_test import simple_rest_api
from .resource_test import simple_resource

INTEGRATION_RESOURCE_PLURAL = "integrations"
MODIFY_WAIT_AFTER_SECONDS = 60


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture
def simple_integration(simple_resource, apigateway_client) -> Tuple[k8s.CustomResourceReference, Dict, Dict]:
    integration_name = random_suffix_name('simple-integration', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['INTEGRATION_NAME'] = integration_name
    (ref, resource_cr, rest_api_cr, resource_query) = simple_resource
    replacements['REST_API_REF_NAME'] = rest_api_cr['metadata']['name']
    replacements['RESOURCE_REF_NAME'] = resource_cr['metadata']['name']

    resource_data = load_apigateway_resource(
        'integration_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)
    resource_query = {**resource_query, **{'httpMethod': 'GET'}}
    apigateway_client.put_method(**resource_query, authorizationType='AWS_IAM')

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, INTEGRATION_RESOURCE_PLURAL,
        integration_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=15)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr, resource_query

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(apigateway_client.get_integration, **resource_query))
    apigateway_client.delete_method(**resource_query)


@service_marker
@pytest.mark.canary
class TestIntegration:
    def test_create_update_integration(self, simple_integration, apigateway_client):
        (ref, cr, resource_query) = simple_integration
        get_integration = partial(apigateway_client.get_integration, **resource_query)

        assert safe_get(get_integration) is not None

        updates = {
            'requestTemplates': {
                'application/json': '{}',
            },
            'passthroughBehavior': 'WHEN_NO_TEMPLATES',
            'timeoutInMillis': 99,
        }
        k8s.patch_custom_resource(ref, {'spec': updates})
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        aws_resource = get_integration()
        updated_fields = {field: aws_resource[field] for field in updates.keys()}
        assert updated_fields == updates
