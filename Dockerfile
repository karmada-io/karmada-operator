FROM rongzhang/ansible-operator:v1.22.3

COPY requirements.yml ${HOME}/requirements.yml
RUN ansible-galaxy collection install -r ${HOME}/requirements.yml \
 && mkdir -p ${HOME}/.ansible/plugins/inventory/ && chmod -R ug+rwx ${HOME}/.ansible

COPY watches.yaml ${HOME}/watches.yaml
COPY roles/ ${HOME}/roles/
COPY playbooks/ ${HOME}/playbooks/
COPY plugins/inventory/karmada.py ${HOME}/.ansible/plugins/inventory/karmada.py
COPY inventory.yml ${HOME}/inventory.yml

ENV ANSIBLE_INVENTORY=${HOME}/inventory.yml \
    ANSIBLE_INVENTORY_ENABLED=auto,karmada
