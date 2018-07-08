#!/bin/sh

swagger generate server -f ../mindwell-server/web/swagger.yaml -P models.UserID -A mindwell-images -O PutMeAvatar -O PutMeCover -M Avatar -M Cover -M UserID -M Error
