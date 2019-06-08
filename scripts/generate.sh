#!/bin/sh

swagger generate server -f ../mindwell-server/web/swagger.yaml -P models.UserID -A mindwell-images \
 -O PutMeAvatar -O PutMeCover -O PostImages \
 -M Avatar -M Cover -M ImageSize -M Image -M UserID -M Error -M User
