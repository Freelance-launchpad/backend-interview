FROM 406151221390.dkr.ecr.eu-west-1.amazonaws.com/dockerhub/sqitch/sqitch:v1.4.1

COPY ./services/activity-reports/activity-reports/db/sqitch /etc/activity-reports/sqitch
WORKDIR /etc/activity-reports/sqitch

ENTRYPOINT ["sqitch", "deploy", "--verify"]
