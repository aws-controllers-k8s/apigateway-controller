"""Utilities for working with API Gateway resources"""

import datetime
import time
import typing

import boto3
import pytest

from e2e import SERVICE_NAME

DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS = 60*30
DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS = 15

GetResourceFunc = typing.NewType(
    'GetResourceFunc',
    typing.Callable[[], dict],
)


def wait_until_deleted(
        get_resource: GetResourceFunc,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS,
) -> None:
    """Waits until a resource is deleted from the API Gateway API

    Usage:
        from e2e.common.waiter import wait_until_deleted

        wait_until_deleted(partial(apigateway_client.get_resource, **resource_query))

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while safe_get(get_resource) is not None:
        if datetime.datetime.now() >= timeout:
            pytest.fail('Timed out waiting for resource to be deleted in API Gateway')
        time.sleep(interval_seconds)


def safe_get(get_resource: GetResourceFunc):
    try:
        return get_resource()
    except boto3.client(SERVICE_NAME).exceptions.NotFoundException:
        return None
