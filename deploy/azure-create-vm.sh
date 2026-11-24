#!/usr/bin/env bash
# Create the LeeTKOB deployment VM on Azure (run it yourself; it spends credit).
#   az login
#   bash deploy/azure-create-vm.sh            # defaults below
#   LOCATION=eastasia DNS_LABEL=leetkob-k23 bash deploy/azure-create-vm.sh
# It creates one resource group holding: an Ubuntu 24.04 VM, a static public IP with a DNS
# name, and a firewall that allows 80/443 from anywhere and SSH only from this machine's IP.
# Delete everything later with:  az group delete --name "$RG"   (only after PA#5/PA#6 are marked)
set -euo pipefail

RG="${RG:-leetkob-pa3}"
LOCATION="${LOCATION:-southeastasia}"
VM="${VM:-leetkob-vm}"
SIZE="${SIZE:-Standard_B2s}"          # 2 vCPU, 4 GiB
DNS_LABEL="${DNS_LABEL:-leetkob-awad-k23}"
ADMIN="${ADMIN:-azureuser}"
SSH_KEY="${SSH_KEY:-$HOME/.ssh/id_ed25519.pub}"

[ -f "$SSH_KEY" ] || { echo "No public key at $SSH_KEY" >&2; exit 1; }
az account show --query "{subscription:name, user:user.name}" -o table

MY_IP="$(curl -fsS https://api.ipify.org)"
echo "SSH will be allowed only from ${MY_IP}/32"

echo "== Resource group $RG in $LOCATION"
az group create --name "$RG" --location "$LOCATION" -o none

echo "== VM $VM ($SIZE, Ubuntu 24.04), static IP, DNS label $DNS_LABEL"
az vm create \
  --resource-group "$RG" --name "$VM" --location "$LOCATION" \
  --image Canonical:ubuntu-24_04-lts:server:latest --size "$SIZE" \
  --admin-username "$ADMIN" --ssh-key-values "$SSH_KEY" --authentication-type ssh \
  --public-ip-sku Standard --public-ip-address-allocation static \
  --public-ip-address-dns-name "$DNS_LABEL" \
  --os-disk-size-gb 30 --nsg-rule NONE -o none

NSG="$(az network nsg list -g "$RG" --query "[0].name" -o tsv)"
echo "== Firewall $NSG: SSH from ${MY_IP}, HTTP/HTTPS from anywhere"
az network nsg rule create -g "$RG" --nsg-name "$NSG" -n ssh-from-admin --priority 1000 \
  --access Allow --protocol Tcp --direction Inbound --destination-port-ranges 22 \
  --source-address-prefixes "${MY_IP}/32" -o none
az network nsg rule create -g "$RG" --nsg-name "$NSG" -n web --priority 1010 \
  --access Allow --protocol Tcp --direction Inbound --destination-port-ranges 80 443 \
  --source-address-prefixes Internet -o none

FQDN="$(az network public-ip list -g "$RG" --query "[0].dnsSettings.fqdn" -o tsv)"
IP="$(az network public-ip list -g "$RG" --query "[0].ipAddress" -o tsv)"
echo
echo "VM ready: ${FQDN} (${IP})"
echo "Next (runbook step 2):"
echo "  scp deploy/vm-bootstrap.sh deploy/vm-accounts.sh ${ADMIN}@${FQDN}:"
echo "  ssh ${ADMIN}@${FQDN} 'sudo DOMAIN=${FQDN} bash vm-bootstrap.sh'"
