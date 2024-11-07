#!/usr/bin/env bash
certbot -n -d ${AWS_BEANSTALK_DOMAIN} certonly --nginx --agree-tos --email ${CERTBOT_EMAIL}