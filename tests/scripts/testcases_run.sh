#!/bin/bash
set -euxo pipefail

export ANSIBLE_REMOTE_USER=$SSH_USER
export ANSIBLE_BECOME=true
export ANSIBLE_BECOME_USER=root

cd tests && make create-${CI_PLATFORM} -s ; cd -
ansible-playbook tests/cloud_playbooks/wait-for-ssh.yml