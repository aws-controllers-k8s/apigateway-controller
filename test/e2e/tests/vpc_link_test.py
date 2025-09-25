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

"""Integration tests for the VPCLink resource
"""

import logging
import time
from typing import Dict, Tuple
from functools import partial

import boto3
import pytest
from acktest import tags
from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_apigateway_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common.waiter import wait_until_deleted, safe_get
from e2e.bootstrap_resources import get_bootstrap_resources

VPC_LINK_RESOURCE_PLURAL = "vpclinks"
MODIFY_WAIT_AFTER_SECONDS = 40
MAX_WAIT_FOR_SYNCED_MINUTES = 20


@pytest.fixture(scope='module')
def apigateway_client():
    return boto3.client(SERVICE_NAME)


@pytest.fixture
def simple_vpc_link(apigateway_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    vpc_link_name = random_suffix_name('simple-vpc-link', 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements['VPC_LINK_NAME'] = vpc_link_name
    replacements['NLB_TARGET_ARN'] = get_bootstrap_resources().NetworkLoadBalancer.arn

    resource_data = load_apigateway_resource(
        'vpc_link_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, VPC_LINK_RESOURCE_PLURAL,
        vpc_link_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=15)

    assert cr is not None
    assert cr['status']['id'] is not None
    assert cr['status']['status'] == 'PENDING'
    assert k8s.get_resource_exists(ref)
    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(apigateway_client.get_vpc_link, vpcLinkId=cr['status']['id']))


@service_marker
@pytest.mark.canary
class TestVPCLink:
    def test_create_update_vpc_link(self, simple_vpc_link, apigateway_client):
        (ref, cr) = simple_vpc_link
        vpc_link_id = cr['status']['id']
        assert safe_get(partial(apigateway_client.get_vpc_link, vpcLinkId=vpc_link_id)) is not None

        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_READY,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )

        vpc_link = k8s.get_resource(ref)
        assert vpc_link['status']['status'] == 'AVAILABLE'

        updates = {
            'spec': {
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

        aws_res = apigateway_client.get_vpc_link(vpcLinkId=vpc_link_id)
        expected_tags = updates['spec'].pop('tags')
        updated_fields = {field: aws_res[field] for field in updates['spec'].keys()}
        assert updated_fields == updates['spec']
        tags.assert_equal_without_ack_tags(
            expected=expected_tags,
            actual=aws_res['tags'],
        )
