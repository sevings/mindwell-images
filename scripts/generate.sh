#!/bin/sh

swagger generate server -f ../mindwell-server/web/swagger.yaml -P models.UserID -A mindwell-images -O PutUsersMeAvatar -M Avatar -M UserID -M Error
