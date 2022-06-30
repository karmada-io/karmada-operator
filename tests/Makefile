INVENTORY=$(PWD)/../inventory/sample/${CI_JOB_NAME}-${BUILD_NUMBER}.ini

$(HOME)/.ssh/id_rsa:
	mkdir -p $(HOME)/.ssh
	echo $(PRIVATE_KEY) | base64 -d > $(HOME)/.ssh/id_rsa
	chmod 400 $(HOME)/.ssh/id_rsa

init-do: $(HOME)/.ssh/id_rsa
	echo $(DO_PRIVATE_KEY) | base64 -d > $(HOME)/.ssh/id_rsa

init-packet:
	echo $(PACKET_VM_SSH_PRIVATE_KEY) | base64 -d > $(HOME)/.ssh/id_rsa
	chmod 400 $(HOME)/.ssh/id_rsa

create-packet: init-packet
	ansible-playbook cloud_playbooks/create-packet.yml -c local \
	$(ANSIBLE_LOG_LEVEL) \
	-e @"files/${CI_JOB_NAME}.yml" \
	-e test_id=$(TEST_ID) \
	-e inventory_path=$(INVENTORY)

delete-packet:
	ansible-playbook cloud_playbooks/delete-packet.yml -c local \
	$(ANSIBLE_LOG_LEVEL) \
	-e @"files/${CI_JOB_NAME}.yml" \
	-e test_id=$(TEST_ID) \
	-e inventory_path=$(INVENTORY)