# resource_list.py
import boto3
import requests
import json


# credential key parsing
def key_parsing():
    with open('./credentials.txt', 'r') as f:
        keys = f.readlines()
        keys = [x.replace('\n', '') for x in keys]
        companys = [company.replace('[', '').replace(']', '') for company in keys if '[' in company]
        accesskeys = [accesskey.replace('aws_access_key_id = ', '') for accesskey in keys if 'aws_access_key_id' in accesskey]
        secretkeys = [secretkey.replace('aws_secret_access_key = ', '') for secretkey in keys if'aws_secret_access_key' in secretkey]
    return list(zip(companys, accesskeys, secretkeys))


# ec2_list check
def ec2_list(client, region):
    ec2 = client.client('ec2', region_name=region)
    instance_information = ec2.describe_instances(
        Filters=[
            {
                'Name': 'instance-state-name',
                'Values': ['running', ]
            }
        ],
    )
    instance_desc = [reservation['Instances'] for reservation in instance_information['Reservations']]
    instance_id = [instance['InstanceId'] for instance_list in instance_desc for instance in instance_list]
    return instance_id


# sns topic create
def create_sns(client):
    sns = client.client('sns')
    topic = sns.create_topic(
        Name='test2'
    )
    subscribe = sns.subscribe(
        TopicArn=topic['TopicArn'],
        Protocol='email-json',
        Endpoint='kky@mz.co.kr',
        ReturnSubscriptionArn=True
    )
    test = client.resource('sns').Topic(topic['TopicArn']).publish(Message=json.dumps({"this": "that"}), MessageAttributes='',)
    print(test)
    #return msg


# main function
def main():
    keys = key_parsing()
    client = boto3.session.Session(keys[0][1], keys[0][2])
    regions = client.client('ec2')
    response = regions.describe_regions()
    region_list = [region['RegionName'] for region in response['Regions']]
    print(create_sns(client))
    #for i in region_list:
        #ec2 = ec2_list(client, i)
        #print(ec2)


main()

