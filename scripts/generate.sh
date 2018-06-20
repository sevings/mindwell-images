#!/bin/sh

swagger generate server -f ../mindwell-server/web/swagger.yaml -P models.UserID -A mindwell-images -O PutUsersMeAvatar -O PutUsersMeCover -M Avatar -M Cover -M UserID -M Error
