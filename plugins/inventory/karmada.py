from __future__ import (absolute_import, division, print_function)
__metaclass__ = type

DOCUMENTATION = '''
    name: karmada-operator
    plugin_type: inventory
    author:
      - Rong Zhang <@riverzhang>
    short_description: Kubernetes (K8s) inventory source for operator sdk
    description:
      - Fetch running KarmadaDeployment for one or more namespace
      - Dynamically generate inventory by parsing the CR information of karmadaDeployment.

    options:
      plugin:
         description: token that ensures this is a source file for the 'karmada-operator' plugin.
         required: True
         choices: ['karmada']
         type: str
      cr_cwd:
         description: custom resource directory where ansible run is stored
         env:
            - name: INVENTORY_CR_CWD
      connection:
          description:
          - Optional dict of cluster connection settings. If no connection are provided, the default
            I(~/.kube/config) and active context will be used, and objects will be returned for all namespaces
            the active user is authorized to access.
          suboptions:
              name:
                  description:
                  - Optional name to assign to the cluster. If not provided, a name is constructed from the server
                    and port.
              kubeconfig:
                  description:
                  - Path to an existing Kubernetes config file. If not provided, and no other connection
                    options are provided, the OpenShift client will attempt to load the default
                    configuration file from I(~/.kube/config.json). Can also be specified via K8S_AUTH_KUBECONFIG
                    environment variable.
              context:
                  description:
                  - The name of a context found in the config file. Can also be specified via K8S_AUTH_CONTEXT environment
                    variable.
              host:
                  description:
                  - Provide a URL for accessing the API. Can also be specified via K8S_AUTH_HOST environment variable.
              api_key:
                  description:
                  - Token used to authenticate with the API. Can also be specified via K8S_AUTH_API_KEY environment
                    variable.
              username:
                  description:
                  - Provide a username for authenticating with the API. Can also be specified via K8S_AUTH_USERNAME
                    environment variable.
              password:
                  description:
                  - Provide a password for authenticating with the API. Can also be specified via K8S_AUTH_PASSWORD
                    environment variable.
              client_cert:
                  description:
                  - Path to a certificate used to authenticate with the API. Can also be specified via K8S_AUTH_CERT_FILE
                    environment variable.
                  aliases: [ cert_file ]
              client_key:
                  description:
                  - Path to a key file used to authenticate with the API. Can also be specified via K8S_AUTH_KEY_FILE
                    environment variable.
                  aliases: [ key_file ]
              ca_cert:
                  description:
                  - Path to a CA certificate used to authenticate with the API. Can also be specified via
                    K8S_AUTH_SSL_CA_CERT environment variable.
                  aliases: [ ssl_ca_cert ]
              validate_certs:
                  description:
                  - "Whether or not to verify the API server's SSL certificates. Can also be specified via
                    K8S_AUTH_VERIFY_SSL environment variable."
                  type: bool
                  aliases: [ verify_ssl ]
    requirements:
    - "python >= 2.7"
    - "openshift >= 0.6"
    - "PyYAML >= 3.11"
'''

EXAMPLES = '''
# File could be named inventory.yaml or inventory.yml
# Authenticate with token, and return all pods and services for all namespaces
plugin: karmada
connection:
  host: https://193.168.6.4:8443
  api_key: xxxxxxxxxxxxxxxx
  validate_certs: false
# Use a custom config file, and a specific context.
plugin: karmada
connection:
  kubeconfig: /path/to/config
  context: 'cluster.local'
'''

import os
import json

from ansible.plugins.filter.core import to_uuid
from ansible_collections.kubernetes.core.plugins.module_utils.common import K8sAnsibleMixin, HAS_K8S_MODULE_HELPER, k8s_import_exception, get_api_client
from ansible.plugins.inventory import BaseInventoryPlugin, Constructable, Cacheable

try:
    from openshift.dynamic.exceptions import DynamicApiError, ResourceNotFoundError
except ImportError:
    pass


def format_dynamic_api_exc(exc):
    if exc.body:
        if exc.headers and exc.headers.get('Content-Type') == 'application/json':
            message = json.loads(exc.body).get('message')
            if message:
                return message
        return exc.body
    else:
        return '%s Reason: %s' % (exc.status, exc.reason)


class InventoryException(Exception):
    pass


class InventoryModule(BaseInventoryPlugin, Constructable, Cacheable, K8sAnsibleMixin):
    NAME = 'karmada'

    transport = 'kubectl'

    cr_cwd = None

    def parse(self, inventory, loader, path, cache=True):
        super(InventoryModule, self).parse(inventory, loader, path)
        cache_key = self._get_cache_prefix(path)
        config_data = self._read_config_data(path)
        self.setup(config_data, cache, cache_key)

    def setup(self, config_data, cache, cache_key):
        if not HAS_K8S_MODULE_HELPER:
            raise InventoryException(
                "This module requires the OpenShift Python client. Try `pip install openshift`. Detail: {0}".format(k8s_import_exception)
            )

        source_data = None
        if cache and cache_key in self._cache:
            try:
                source_data = self._cache[cache_key]
            except KeyError:
                pass

        if not source_data:
            connection = config_data.get('connection')

            if connection:
                if not isinstance(connection, dict):
                    raise InventoryException("Expecting connection to be a dictionary.")
                self.client = get_api_client(**connection)
            else:
                self.client = get_api_client()

            # get cr info
            self.get_cr_info()

            # build karmadadeployment inventory of cr
            self.build_karmada_inventory()

    # get current directory (playbook) to infer CR meta from it, to later filter inventory
    # /tmp/ansible-operator/runner/<group>/<version>/<Kind>/<namespace>/<cr_name>/project
    # https://github.com/operator-framework/operator-sdk/blob/f298f7c92ee154a0e8123fb13398a7f21720cf0e/internal/ansible/runner/runner.go#L203
    def get_cr_info(self):
        try:
            #print(self.cr_crd)
            # Preference order: cr_cwd variable > inventory.yml > environment variable > cwd
            if self.cr_cwd is not None:
                cwd = self.cr_cwd
            elif self.get_option("cr_cwd") is not None:
                cwd = self.get_option("cr_cwd")
            elif "INVENTORY_CR_CWD" in os.environ:
                cwd = os.getenv('INVENTORY_CR_CWD')
            else:
                cwd = os.getcwd()

            dir_parts = cwd.split(os.path.sep)
            self.cr_name = dir_parts[-2]
            self.cr_namespace = dir_parts[-3]
            self.cr_kind = dir_parts[-4]
            self.cr_version = dir_parts[-5]
            self.cr_group = dir_parts[-6]
            self.api_version = self.cr_group + '/' + self.cr_version
        except Exception as cwd_exc:
            raise InventoryException('Error getting CR from cwd: %s' % cwd_exc)

    def build_karmada_inventory(self):
        try:
            cr_api = self.client.resources.get(api_version=self.api_version, kind=self.cr_kind)
        except ResourceNotFoundError as notFound:
            raise InventoryException('CR %s not found: %s' % (self.cr_name, notFound))

        try:
            obj = cr_api.get(name=self.cr_name, namespace=self.cr_namespace)
        except DynamicApiError as exc:
            self.display.debug(exc)
            raise InventoryException('Error fetching CR %s: %s' % (self.cr_name, format_dynamic_api_exc(exc)))

        add_member_group = self.inventory.add_group('add-member')
        for member in obj.spec.members:
            self.inventory.add_host(member.name, add_member_group)
            self.inventory.set_variable(member.name, "ansible_connection", "local")
            self.inventory.set_variable(member.name, "member_cluster_name", member.name)
            self.inventory.set_variable(member.name, "sync_mode", member.syncMode)
            self.inventory.set_variable(member.name, "secret_name", member.kubeConfigSecretName)

        # todo(riverzhang)
        # del-member group of delete member cluster from karmada
