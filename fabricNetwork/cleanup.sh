docker container stop $(docker ps -a -q)
docker container rm $(docker ps -a -q)
docker volume rm $(docker volume list -q)
sudo rm -r test-network/docker-compose
sudo rm -r test-network/configtx
sudo rm test-network/crypto.yaml
sudo rm -r test-network/organizations
sudo rm -r test-network/channel-artifacts
sudo rm -r test-network/system-genesis-block
# sudo rm -r test-network/chaincode