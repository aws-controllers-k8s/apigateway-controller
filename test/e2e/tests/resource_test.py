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

RESOURCE_RESOURCE_PLURAL = "resources"
MODIFY_WAIT_AFTER_SECONDS = 60


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture
def simple_resource(simple_rest_api, apigateway_client) -> Tuple[k8s.CustomResourceReference, Dict, Dict, Dict]:
    resource_name = random_suffix_name('simple-resource', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['RESOURCE_NAME'] = resource_name
    (ref, rest_api_cr) = simple_rest_api
    replacements['REST_API_REF_NAME'] = rest_api_cr['spec']['name']
    replacements['PARENT_ID'] = rest_api_cr['status']['rootResourceID']

    resource_data = load_apigateway_resource(
        'resource_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_RESOURCE_PLURAL,
        resource_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=15)

    assert cr is not None
    assert cr['status']['id'] is not None
    assert k8s.get_resource_exists(ref)
    resource_query = {
        'restApiId': rest_api_cr['status']['id'],
        'resourceId': cr['status']['id'],
    }

    yield ref, cr, rest_api_cr, resource_query

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(apigateway_client.get_resource, **resource_query))


@service_marker
@pytest.mark.canary
class TestResource:
    def test_create_update_resource(self, simple_resource, apigateway_client):
        (ref, cr, _, resource_query) = simple_resource
        assert safe_get(partial(apigateway_client.get_resource, **resource_query)) is not None

        updates = {
            'spec': {
                'pathPart': 'updated-path',
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        aws_resource = apigateway_client.get_resource(**resource_query)
        assert aws_resource['pathPart'] == updates['spec']['pathPart']
