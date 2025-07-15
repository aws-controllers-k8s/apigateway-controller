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

"""Integration tests for the API Gateway DomainName resource
"""

import logging
import time

import boto3
import pytest

from acktest.k8s import resource as k8s
from acktest.k8s import condition
from acktest.resources import random_suffix_name

from e2e import service_marker, load_apigateway_resource, CRD_GROUP, CRD_VERSION
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.replacement_values import REPLACEMENT_VALUES

CREATE_WAIT_AFTER_SECONDS = 60

@service_marker
@pytest.mark.canary
class TestDomainName:

    def test_crud_domain_name(self):
        test_data = REPLACEMENT_VALUES.copy()
        domain_name = random_suffix_name("test-domain", 24)
        test_data["DOMAIN_RES_NAME"] = domain_name
        test_data["DOMAIN_NAME"] = f"{domain_name}.example.com"
        test_data["CERT_ARN"] = "arn:aws:acm:us-west-2:123456789012:certificate/invalid-cert-id"
        
        # Create a custom resource reference
        domain_name_ref = k8s.CustomResourceReference(
            CRD_GROUP, CRD_VERSION, "domainnames",
            domain_name, namespace="default",
        )
        
        # Load the resource data
        domain_name_data = load_apigateway_resource(
            "domain_name_simple",
            additional_replacements=test_data,
        )
        
        logging.debug(f"domain name ref is {domain_name_ref}, data: {domain_name_data}")
        
        # Attempting Create expecting BadRequestException
        k8s.create_custom_resource(domain_name_ref, domain_name_data)
        k8s.wait_resource_consumed_by_controller(domain_name_ref)
        condition.assert_type_status(domain_name_ref, condition.CONDITION_TYPE_TERMINAL)

        # Verify terminal condition with expected error message
        expected_msg = "BadRequestException"
        terminal_condition = k8s.get_resource_condition(domain_name_ref, condition.CONDITION_TYPE_TERMINAL)
        assert expected_msg in terminal_condition['message']
        
        # Clean up
        k8s.delete_custom_resource(domain_name_ref)