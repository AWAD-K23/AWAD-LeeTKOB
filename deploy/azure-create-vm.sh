#!/usr/bin/env bash
# Create the LeeTKOB deployment VM on Azure (run it yourself; it spends credit).
#   az login
#   bash deploy/azure-create-vm.sh            # defaults below
#   LOCATION=eastasia DNS_LABEL=leetkob-k23 bash deploy/azure-create-vm.sh
# It creates, in a new or an existing resource group (RG=<name>): an Ubuntu 24.04 VM, a static
# public IP "<vm>-ip" with a DNS name, and a firewall "<vm>-nsg" that allows 80/443 from anywhere
# and SSH only from this machine's IP. To stop paying for compute: az vm deallocate -g RG -n VM.
# Remove it after PA#5/PA#6 are marked: delete the VM, then its disk, NIC, <vm>-ip and <vm>-nsg.
# Do not delete the whole group if it holds other resources.
set -euo pipefail

RG="${RG:-leetkob-pa3}"
# The region must be one where the Gemini API is offered: eastasia is Hong Kong, and Gemini
# answers every call from there with "User location is not supported" (seen 2026-11-24).
# Azure for Students allows only malaysiawest, eastasia, japaneast, uaenorth and
# indonesiacentral (policy sys.regionrestriction). On 2026-11-24 the only small size with quota
# among the Gemini regions was the ARM64 B2pls_v2 in japaneast; every image we build is multi-arch.
# The VM may live in another region than its group.
LOCATION="${LOCATION:-japaneast}"
VM="${VM:-leetkob-vm}"
SIZE="${SIZE:-Standard_B2pls_v2}"     # 2 vCPU (ARM64), 4 GiB, burstable
IMAGE="${IMAGE:-Canonical:ubuntu-24_04-lts:server-arm64:latest}"   # x86 sizes: ...:server:latest
DNS_LABEL="${DNS_LABEL:-leetkob-awad-k23}"
ADMIN="${ADMIN:-azureuser}"
SSH_KEY="${SSH_KEY:-$HOME/.ssh/id_ed25519.pub}"

[ -f "$SSH_KEY" ] || { echo "No public key at $SSH_KEY" >&2; exit 1; }
az account show --query "{subscription:name, user:user.name}" -o table

MY_IP="$(curl -fsS https://api.ipify.org)"
echo "SSH will be allowed only from ${MY_IP}/32"

if [ "$(az group exists --name "$RG")" = "true" ]; then
  echo "== Using existing resource group $RG ($(az group show --name "$RG" --query location -o tsv)); VM in $LOCATION"
else
  echo "== Resource group $RG in $LOCATION"
  az group create --name "$RG" --location "$LOCATION" -o none
fi

echo "== VM $VM ($SIZE, Ubuntu 24.04), static IP, DNS label $DNS_LABEL"
az vm create \
  --resource-group "$RG" --name "$VM" --location "$LOCATION" \
  --image "$IMAGE" --size "$SIZE" \
  --admin-username "$ADMIN" --ssh-key-values "$SSH_KEY" --authentication-type ssh \
  --public-ip-address "${VM}-ip" --nsg "${VM}-nsg" \
  --public-ip-sku Standard --public-ip-address-allocation static \
  --public-ip-address-dns-name "$DNS_LABEL" \
  --os-disk-size-gb 30 --nsg-rule NONE -o none

NSG="${VM}-nsg"
echo "== Firewall $NSG: SSH from ${MY_IP}, HTTP/HTTPS from anywhere"
az network nsg rule create -g "$RG" --nsg-name "$NSG" -n ssh-from-admin --priority 1000 \
  --access Allow --protocol Tcp --direction Inbound --destination-port-ranges 22 \
  --source-address-prefixes "${MY_IP}/32" -o none
az network nsg rule create -g "$RG" --nsg-name "$NSG" -n web --priority 1010 \
  --access Allow --protocol Tcp --direction Inbound --destination-port-ranges 80 443 \
  --source-address-prefixes Internet -o none

FQDN="$(az network public-ip show -g "$RG" -n "${VM}-ip" --query dnsSettings.fqdn -o tsv)"
IP="$(az network public-ip show -g "$RG" -n "${VM}-ip" --query ipAddress -o tsv)"
echo
echo "VM ready: ${FQDN} (${IP})"
echo "Next (runbook step 2):"
echo "  scp deploy/vm-bootstrap.sh deploy/vm-accounts.sh ${ADMIN}@${FQDN}:"
echo "  ssh ${ADMIN}@${FQDN} 'sudo DOMAIN=${FQDN} bash vm-bootstrap.sh'"
